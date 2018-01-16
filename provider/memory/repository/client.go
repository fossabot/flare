// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/diegobernardes/flare"
)

type Client struct {
	resource        Resource
	resourceOptions []func(*Resource)
	subscription    Subscription
	document        Document
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

func NewClient(options ...func(*Client)) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	c.resource.repository = &c.subscription
	c.subscription.resourceRepository = &c.resource
	c.subscription.documentRepository = &c.document

	c.resource.init(c.resourceOptions...)
	c.subscription.init()
	c.document.init()
	return c
}

func ClientResourceOptions(options ...func(*Resource)) func(*Client) {
	return func(c *Client) {
		c.resourceOptions = options
	}
}
