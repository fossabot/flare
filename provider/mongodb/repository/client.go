// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/pkg/errors"

	"github.com/diegobernardes/flare"
	mongodb "github.com/diegobernardes/flare/provider/mongodb"
)

type Client struct {
	base         *mongodb.Client
	resource     Resource
	subscription Subscription
	document     Document
}

func (c *Client) Resource() flare.ResourceRepositorier {
	return &c.resource
}

func (c *Client) Subscription() flare.SubscriptionRepositorier {
	return &c.subscription
}

func (c *Client) Document() flare.DocumentRepositorier {
	return &c.document
}

func (c *Client) Stop() error {
	return nil
}

func NewClient(base *mongodb.Client) (*Client, error) {
	c := &Client{base: base}
	c.resource.subscriptionRepository = &c.subscription
	c.subscription.resourceRepository = &c.resource
	c.resource.client = base
	c.subscription.client = base
	c.document.client = base

	if err := c.resource.init(); err != nil {
		return nil, errors.Wrap(err, "error during resource repository initialization")
	}

	if err := c.subscription.init(); err != nil {
		return nil, errors.Wrap(err, "error during subscription repository initialization")
	}

	if err := c.document.init(); err != nil {
		return nil, errors.Wrap(err, "error during document repository initialization")
	}

	return c, nil
}
