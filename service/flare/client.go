// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/diegobernardes/flare/infra/config"
)

// Variables set with ldflags during compilation.
var (
	Version   = ""
	BuildTime = ""
	Commit    = ""
	GoVersion = runtime.Version()
)

type Client struct {
	config     config.Client
	logger     log.Logger
	repository repositoryGrouper
	server     *server
	worker     *worker
}

func (c *Client) Start() error {
	if err := c.config.Init(); err != nil {
		return errors.Wrap(err, "error during config initialization")
	}

	if err := c.initLogger(); err != nil {
		return errors.Wrap(err, "error during log initialization")
	}

	if err := c.initRepository(); err != nil {
		return errors.Wrap(err, "error during repository initialization")
	}

	if err := c.worker.init(); err != nil {
		return errors.Wrap(err, "error during worker initialization")
	}

	resourceHandler, err := c.initDomainResource()
	if err != nil {
		return errors.Wrap(err, "error during domain resource initialization")
	}

	subscriptionHandler, err := c.initDomainSubscription()
	if err != nil {
		return errors.Wrap(err, "error during domain subscription initializer")
	}

	documentHandler, err := c.initDomainDocument()
	if err != nil {
		return errors.Wrap(err, " error during domain document initialization")
	}

	if err := c.initServer(resourceHandler, subscriptionHandler, documentHandler); err != nil {
		return errors.Wrap(err, "error during HTTP server initialization")
	}

	return nil
}

func (c *Client) Stop() error {
	level.Info(c.logger).Log("message", "signal to close the process received")
	fmt.Println("closing server")
	if err := c.server.stop(); err != nil {
		panic(err)
	}

	fmt.Println("closing worker")
	if err := c.worker.stop(); err != nil {
		panic(err)
	}

	fmt.Println("everything is closed")

	// depois parar as filas
	// depois parar os bancos

	return nil
}

func (c *Client) Setup(ctx context.Context) error {
	return nil
}

func NewClient(options ...func(*Client)) (*Client, error) {
	c := &Client{worker: &worker{}}
	c.worker.client = c

	for _, option := range options {
		option(c)
	}

	return c, nil
}

func ClientConfig(config string) func(*Client) {
	return func(c *Client) {
		c.config.Content = config
	}
}
