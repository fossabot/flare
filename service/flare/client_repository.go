// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/diegobernardes/flare"
	memory "github.com/diegobernardes/flare/provider/memory/repository"
	"github.com/diegobernardes/flare/provider/mongodb"
	mongoDBRepository "github.com/diegobernardes/flare/provider/mongodb/repository"
)

type repositoryGrouper interface {
	Resource() flare.ResourceRepositorier
	Subscription() flare.SubscriptionRepositorier
	Document() flare.DocumentRepositorier
}

func (c *Client) initRepository() error {
	var (
		err      error
		provider = c.config.GetString("provider.repository")
	)

	switch provider {
	case "", "memory":
		c.repository, err = c.initRepositoryMemory()
	case "mongodb":
		c.repository, err = c.initRepositoryMongoDB()
	default:
		c.repository, err = nil, fmt.Errorf("invalid repository.provider '%s'", provider)
	}

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) initRepositoryMemory() (repositoryGrouper, error) {
	return memory.NewClient(), nil
}

func (c *Client) initRepositoryMongoDB() (repositoryGrouper, error) {
	options, err := c.initRepositoryMongoDBOptions()
	if err != nil {
		return nil, err
	}

	client, err := mongodb.NewClient(options...)
	if err != nil {
		return nil, errors.Wrap(err, "error during MongoDB repository initialization")
	}

	repository, err := mongoDBRepository.NewClient(client)
	if err != nil {
		return nil, errors.Wrap(err, "error during repository initialization")
	}

	return repository, nil
}

func (c *Client) initRepositoryMongoDBOptions() ([]func(*mongodb.Client), error) {
	var options []func(*mongodb.Client)

	addrs := c.config.GetStringSlice("provider.mongodb.addrs")
	if len(addrs) == 0 {
		addrs = append(addrs, "localhost:27017")
	}
	options = append(options, mongodb.ClientAddrs(addrs))

	database := c.config.GetString("provider.mongodb.database")
	if database == "" {
		database = "flare"
	}
	options = append(options, mongodb.ClientDatabase(database))

	username := c.config.GetString("provider.mongodb.username")
	if username != "" {
		options = append(options, mongodb.ClientUsername(username))
	}

	password := c.config.GetString("provider.mongodb.password")
	if password != "" {
		options = append(options, mongodb.ClientPassword(password))
	}

	replicaSet := c.config.GetString("provider.mongodb.replica-set")
	if replicaSet != "" {
		options = append(options, mongodb.ClientReplicaSet(replicaSet))
	}

	poolLimit := c.config.GetInt("provider.mongodb.pool-limit")
	if poolLimit >= 0 {
		options = append(options, mongodb.ClientPoolLimit(poolLimit))
	} else if poolLimit < 0 {
		return nil, fmt.Errorf("invalid mongodb.pool-limit, it cannot be less then zero '%d'", poolLimit)
	}

	timeout, err := c.config.GetDuration("provider.mongodb.timeout")
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf(
				"error during transform '%s' to time.Duration",
				c.config.GetString("provider.mongodb.timeout"),
			),
		)
	}
	if timeout < 0 {
		return nil, fmt.Errorf("mongodb.timeout can't be less then zero '%d'", timeout)
	}
	if timeout > 0 {
		options = append(options, mongodb.ClientTimeout(timeout))
	}

	return options, nil
}

func (c *Client) stopRepository(rawGroup repositoryGrouper) {
	group, ok := rawGroup.(*mongoDBRepository.Client)
	if !ok {
		return
	}
	group.Stop()
}
