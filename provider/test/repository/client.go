// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/diegobernardes/flare"
	memoryRepository "github.com/diegobernardes/flare/provider/memory/repository"
)

type Client struct {
	base *memoryRepository.Client

	document        *Document
	documentOptions []func(*Document)

	resource        *Resource
	resourceOptions []func(*Resource)

	subscription        *Subscription
	subscriptionOptions []func(*Subscription)
}

func (c *Client) Resource() flare.ResourceRepositorier {
	return c.resource
}

func (c *Client) Subscription() flare.SubscriptionRepositorier {
	return c.subscription
}

func (c *Client) Document() flare.DocumentRepositorier {
	return c.document
}

func NewClient(options ...func(*Client)) *Client {
	return NewClientWithBase(memoryRepository.NewClient(), options...)
}

func NewClientWithBase(base *memoryRepository.Client, options ...func(*Client)) *Client {
	c := &Client{base: base}

	c.documentOptions = append(c.documentOptions, DocumentRepository(c.base.Document()))
	c.resourceOptions = append(c.resourceOptions, ResourceRepository(c.base.Resource()))
	c.subscriptionOptions = append(
		c.subscriptionOptions,
		SubscriptionRepository(c.base.Subscription()),
	)

	for _, option := range options {
		option(c)
	}

	c.document = newDocument(c.documentOptions...)
	c.resource = newResource(c.resourceOptions...)
	c.subscription = newSubscription(c.subscriptionOptions...)
	return c
}

func ClientDocumentOptions(options ...func(*Document)) func(*Client) {
	return func(c *Client) {
		c.documentOptions = append(c.documentOptions, options...)
	}
}

func ClientResourceOptions(options ...func(*Resource)) func(*Client) {
	return func(c *Client) {
		c.resourceOptions = append(c.resourceOptions, options...)
	}
}

func ClientSubscriptionOptions(options ...func(*Subscription)) func(*Client) {
	return func(c *Client) {
		c.subscriptionOptions = append(c.subscriptionOptions, options...)
	}
}
