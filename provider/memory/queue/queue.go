// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

import (
	"context"
	"sync"
	"time"
)

type message struct {
	content     []byte
	locked      bool
	lockedUntil time.Time
	context     context.Context
}

type Client struct {
	mutex          sync.Mutex
	messages       [][]byte
	timeoutProcess time.Duration
}

func (q *Client) Push(_ context.Context, content []byte) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, content)
	return nil
}

// TODO: falta colocar configuracoes de timeout, erro, etc... aqui.
func (q *Client) Pull(ctx context.Context, fn func(context.Context, []byte) error) error {
	q.mutex.Lock()

	if len(q.messages) == 0 {
		q.mutex.Unlock()
		<-time.After(1 * time.Second)
		return nil
	}
	defer q.mutex.Unlock()

	ctx, ctxCancel := context.WithTimeout(ctx, q.timeoutProcess)
	defer ctxCancel()

	if err := fn(ctx, q.messages[0]); err != nil {
		return err
	}
	q.messages = q.messages[1:]
	return nil
}

func NewClient(options ...func(*Client)) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	if c.timeoutProcess == 0 {
		c.timeoutProcess = time.Duration(1 * time.Hour)
	}

	return c
}

func ClientProcessTimeout(timeout time.Duration) func(*Client) {
	return func(c *Client) {
		c.timeoutProcess = timeout
	}
}
