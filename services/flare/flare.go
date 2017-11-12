// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pkg/errors"

	"github.com/diegobernardes/flare"
	"github.com/diegobernardes/flare/document"
	infraHTTP "github.com/diegobernardes/flare/infra/http"
	"github.com/diegobernardes/flare/infra/task"
	"github.com/diegobernardes/flare/resource"
	"github.com/diegobernardes/flare/subscription"
)

const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

// Client is used to initialize Flare.
type Client struct {
	server    *server
	rawConfig string
	config    *config
	logger    log.Logger
}

// Start is used to start the service.
func (c *Client) Start() error {
	config, err := newConfig(configContent(c.rawConfig))
	if err != nil {
		return err
	}
	c.config = config

	if err = c.initLogger(); err != nil {
		return errors.Wrap(err, "error during log initialization")
	}
	c.logger.Log("message", "Starting Flare")

	documentRepository, err := c.config.documentRepository()
	if err != nil {
		return err
	}

	subscriptionRepository, err := c.config.subscriptionRepository()
	if err != nil {
		return err
	}

	resourceService, resourceRepository, err := c.initResourceService(subscriptionRepository)
	if err != nil {
		level.Debug(c.logger).Log(
			"error", err.Error(), "message", "error during resource service initialization",
		)
		return err
	}

	documentService, err := c.initDocumentService(
		documentRepository,
		resourceRepository,
		subscriptionRepository,
	)
	if err != nil {
		return errors.Wrap(err, "error during document service initialization")
	}

	subscriptionService, err := c.initSubscriptionService(resourceRepository, subscriptionRepository)
	if err != nil {
		return errors.Wrap(err, "error during subscription service initialization")
	}

	srv, err := newServer(
		serverAddr(config.getString("http.addr")),
		serverHandlerResource(resourceService),
		serverHandlerSubscription(subscriptionService),
		serverHandlerDocument(documentService),
		serverLogger(c.logger),
	)
	if err != nil {
		return errors.Wrap(err, "error during server initialization")
	}
	c.server = srv
	c.server.start()
	return nil
}

// Stop is used to graceful stop the service.
func (c *Client) Stop() error {
	if err := c.server.stop(); err != nil {
		return errors.Wrap(err, "error during server stop")
	}
	return nil
}

func (c *Client) initLogger() error {
	logger, err := c.initLoggerOutput()
	if err != nil {
		return err
	}

	logger, err = c.initLoggerLevel(logger)
	if err != nil {
		return err
	}

	c.logger = log.With(logger, "time", log.DefaultTimestampUTC)
	return nil
}

func (c *Client) initLoggerOutput() (log.Logger, error) {
	output := c.config.getString("log.output")
	if output == "" {
		output = "stdout"
	}

	format := c.config.getString("log.format")
	if format == "" {
		format = "human"
	}

	switch output {
	case "discard":
		return log.NewNopLogger(), nil
	case "stdout":
		switch format {
		case "human":
			return term.NewLogger(
				log.NewSyncWriter(os.Stdout),
				log.NewLogfmtLogger,
				c.loggerColor,
			), nil
		case "json":
			return log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), nil
		default:
			return nil, fmt.Errorf("invalid log.format '%s'", format)
		}
	default:
		return nil, fmt.Errorf("invalid log.output '%s'", output)
	}
}

func (c *Client) initLoggerLevel(logger log.Logger) (log.Logger, error) {
	logLevel := c.config.getString("log.level")
	if logLevel == "" {
		logLevel = logLevelDebug
	}

	switch logLevel {
	case logLevelDebug:
		return level.Debug(logger), nil
	case logLevelInfo:
		return level.Info(logger), nil
	case logLevelWarn:
		return level.Warn(logger), nil
	case logLevelError:
		return level.Error(logger), nil
	default:
		return nil, fmt.Errorf("invalid log.level '%s'", logLevel)
	}
}

func (c *Client) initResourceService(
	subscriptionRepository flare.SubscriptionRepositorier,
) (*resource.Service, flare.ResourceRepositorier, error) {
	repository, err := c.config.resourceRepository()
	if err != nil {
		panic(err)
	}

	resourceService, err := resource.NewService(
		resource.ServiceGetResourceID(func(r *http.Request) string { return chi.URLParam(r, "id") }),
		resource.ServiceGetResourceURI(func(id string) string {
			return fmt.Sprintf("/resources/%s", id)
		}),
		resource.ServiceParsePagination(infraHTTP.ParsePagination(c.config.getInt("http.default-limit"))),
		resource.ServiceWriteResponse(infraHTTP.WriteResponse(c.logger)),
		resource.ServiceRepository(repository),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error during resource.Service initialization")
	}

	return resourceService, repository, nil
}

func (c *Client) initSubscriptionService(
	resourceRepository flare.ResourceRepositorier,
	subscriptionRepository flare.SubscriptionRepositorier,
) (*subscription.Service, error) {
	subscriptionService, err := subscription.NewService(
		subscription.ServiceDefaultLimit(c.config.getInt("http.default-limit")),
		subscription.ServiceGetResourceId(func(r *http.Request) string {
			return chi.URLParam(r, "resourceId")
		}),
		subscription.ServiceGetSubscriptionId(func(r *http.Request) string {
			return chi.URLParam(r, "id")
		}),
		subscription.ServiceGetSubscriptionURI(func(resourceId, id string) string {
			return fmt.Sprintf("/resources/%s/subscriptions/%s", resourceId, id)
		}),
		subscription.ServiceLogger(c.logger),
		subscription.ServiceResourceRepository(resourceRepository),
		subscription.ServiceSubscriptionRepository(subscriptionRepository),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during subscription.Service initialization")
	}

	return subscriptionService, nil
}

func (c *Client) initDocumentService(
	dr flare.DocumentRepositorier,
	rr flare.ResourceRepositorier,
	sr flare.SubscriptionRepositorier,
) (*document.Service, error) {
	documentPusher, documentPuller, err := c.config.queue("document")
	if err != nil {
		return nil, errors.Wrap(err, "error during queue initialization")
	}

	subscriptionPusher, subscriptionPuller, err := c.config.queue("subscription")
	if err != nil {
		return nil, errors.Wrap(err, "error during queue initialization")
	}

	trigger := &subscription.Trigger{}
	triggerWorker, err := task.NewWorker(
		task.WorkerGoroutines(1),
		task.WorkerProcessor(trigger),
		task.WorkerPuller(subscriptionPuller),
		task.WorkerPusher(subscriptionPusher),
		task.WorkerTimeoutProcess(10*time.Second),
		task.WorkerTimeoutPush(10*time.Second),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during worker initialization")
	}

	err = trigger.Init(
		subscription.TriggerRepository(sr),
		subscription.TriggerHTTPClient(http.DefaultClient),
		subscription.TriggerDocumentRepository(dr),
		subscription.TriggerPusher(triggerWorker),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during subscription.Trigger initialization")
	}
	triggerWorker.Start()

	documentWorker := &document.Worker{}
	jobWorker, err := task.NewWorker(
		task.WorkerGoroutines(1),
		task.WorkerProcessor(documentWorker),
		task.WorkerPuller(documentPuller),
		task.WorkerPusher(documentPusher),
		task.WorkerTimeoutProcess(10*time.Second),
		task.WorkerTimeoutPush(10*time.Second),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during worker initialization")
	}

	err = documentWorker.Init(
		document.WorkerDocumentRepository(dr),
		document.WorkerResourceRepository(rr),
		document.WorkerSubscriptionRepository(sr),
		document.WorkerSubscriptionTrigger(trigger),
		document.WorkerPusher(jobWorker),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during worker initialization")
	}
	jobWorker.Start()

	documentService, err := document.NewService(
		document.ServiceDocumentRepository(dr),
		document.ServiceResourceRepository(rr),
		document.ServiceGetDocumentId(func(r *http.Request) string { return chi.URLParam(r, "*") }),
		document.ServiceLogger(c.logger),
		document.ServiceWorker(documentWorker),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during document.Service initialization")
	}

	return documentService, nil
}

func (c *Client) loggerColor(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}

		switch keyvals[i+1].(level.Value).String() {
		case logLevelDebug:
			return term.FgBgColor{Fg: term.DarkGray}
		case logLevelInfo:
			return term.FgBgColor{Fg: term.Gray}
		case logLevelWarn:
			return term.FgBgColor{Fg: term.Yellow}
		case logLevelError:
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}
	return term.FgBgColor{}
}

// NewClient returns a client to Flare service.
func NewClient(options ...func(*Client)) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	return c
}

// ClientConfig set the config to initialize the Flare client.
func ClientConfig(config string) func(*Client) {
	return func(c *Client) { c.rawConfig = config }
}
