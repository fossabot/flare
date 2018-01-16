// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

import "context"

type Client struct {
	Content []byte
	err     error
}

func (c *Client) Push(ctx context.Context, content []byte) error {
	if c.err != nil {
		return c.err
	}
	c.Content = content
	return nil
}

func NewClient(options ...func(*Client)) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	return c
}

func ClientError(err error) func(*Client) {
	return func(c *Client) {
		c.err = err
	}
}
