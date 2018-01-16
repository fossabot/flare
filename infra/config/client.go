// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Client struct {
	Content string
	viper   *viper.Viper
}

func (c *Client) IsSet(key string) bool { return c.viper.IsSet(key) }

func (c *Client) GetString(key string) string { return c.viper.GetString(key) }

func (c *Client) GetStringSlice(key string) []string { return c.viper.GetStringSlice(key) }

func (c *Client) GetInt(key string) int { return c.viper.GetInt(key) }

func (c *Client) GetBool(key string) bool { return c.viper.GetBool(key) }

func (c *Client) GetDuration(key string) (time.Duration, error) {
	value := c.GetString(key)
	if value == "" {
		return 0, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error during parse '%s' to time.Duration", key))
	}
	return duration, nil
}

func (c *Client) Init() error {
	c.viper = viper.New()
	c.viper.SetConfigType("toml")

	if err := c.viper.ReadConfig(bytes.NewBufferString(c.Content)); err != nil {
		return errors.Wrap(err, "error during config parse")
	}
	return nil
}
