// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"github.com/pkg/errors"

	document "github.com/diegobernardes/flare/domain/document/http"
	resource "github.com/diegobernardes/flare/domain/resource/http"
	subscription "github.com/diegobernardes/flare/domain/subscription/http"
)

func (c *Client) initServer(
	resourceHandler *resource.Handler,
	subscriptionHandler *subscription.Handler,
	documentHandler *document.Handler,
) error {
	key := "http.server.enable"
	if c.config.IsSet(key) && !c.config.GetBool(key) {
		return nil
	}

	c.server = &server{logger: c.logger}
	c.server.handler.resource = resourceHandler
	c.server.handler.subscription = subscriptionHandler
	c.server.handler.document = documentHandler

	addr := c.config.GetString("http.server.addr")
	if addr != "" {
		c.server.addr = addr
	}

	if err := c.server.init(); err != nil {
		return errors.Wrap(err, "error during HTTP server initialization")
	}

	if err := c.server.start(); err != nil {
		return errors.Wrap(err, "error during HTTP server start")
	}

	return nil
}
