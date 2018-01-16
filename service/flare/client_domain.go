// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	document "github.com/diegobernardes/flare/domain/document/http"
	resource "github.com/diegobernardes/flare/domain/resource/http"
	subscription "github.com/diegobernardes/flare/domain/subscription/http"
	infraHTTP "github.com/diegobernardes/flare/infra/http"
)

func (c *Client) initDomainResource() (*resource.Handler, error) {
	writer, err := infraHTTP.NewWriter(c.logger)
	if err != nil {
		return nil, errors.Wrap(err, "error during http.Writer initialization")
	}

	defaultLimit, err := c.getDefaultLimitDomainPagination()
	if err != nil {
		return nil, err
	}

	handler, err := resource.NewHandler(
		resource.HandlerGetResourceID(func(r *http.Request) string { return chi.URLParam(r, "id") }),
		resource.HandlerGetResourceURI(func(id string) string {
			return fmt.Sprintf("/resources/%s", id)
		}),
		resource.HandlerParsePagination(infraHTTP.ParsePagination(defaultLimit)),
		resource.HandlerWriter(writer),
		resource.HandlerRepository(c.repository.Resource()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during resource.Handler initialization")
	}

	return handler, nil
}

func (c *Client) initDomainSubscription() (*subscription.Handler, error) {
	writer, err := infraHTTP.NewWriter(c.logger)
	if err != nil {
		return nil, errors.Wrap(err, "error during http.Writer initialization")
	}

	defaultLimit, err := c.getDefaultLimitDomainPagination()
	if err != nil {
		return nil, err
	}

	subscriptionService, err := subscription.NewHandler(
		subscription.HandlerParsePagination(infraHTTP.ParsePagination(defaultLimit)),
		subscription.HandlerWriter(writer),
		subscription.HandlerGetResourceID(func(r *http.Request) string {
			return chi.URLParam(r, "resourceID")
		}),
		subscription.HandlerGetSubscriptionID(func(r *http.Request) string {
			return chi.URLParam(r, "id")
		}),
		subscription.HandlerGetSubscriptionURI(func(resourceId, id string) string {
			return fmt.Sprintf("/resources/%s/subscriptions/%s", resourceId, id)
		}),
		subscription.HandlerResourceRepository(c.repository.Resource()),
		subscription.HandlerSubscriptionRepository(c.repository.Subscription()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during subscription.Handler initialization")
	}

	return subscriptionService, nil
}

func (c *Client) initDomainDocument() (*document.Handler, error) {
	writer, err := infraHTTP.NewWriter(c.logger)
	if err != nil {
		return nil, errors.Wrap(err, "error during http.Writer initialization")
	}

	documentHandler, err := document.NewHandler(
		document.HandlerDocumentRepository(c.repository.Document()),
		document.HandlerGetDocumentID(func(r *http.Request) string { return chi.URLParam(r, "*") }),
		document.HandlerResourceRepository(c.repository.Resource()),
		document.HandlerSubscriptionTrigger(c.worker.subscriptionPartition),
		document.HandlerWriter(writer),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error during document.Handler initialization")
	}

	return documentHandler, nil
}

func (c *Client) getDefaultLimitDomainPagination() (int, error) {
	key := "domain.pagination.default-limit"
	if !c.config.IsSet(key) {
		return 30, nil
	}

	defaultLimit := c.config.GetInt(key)
	if defaultLimit <= 0 {
		return 0, fmt.Errorf("invalid pagination.default-limit '%d'", defaultLimit)
	}
	return defaultLimit, nil
}
