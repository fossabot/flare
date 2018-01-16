// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	subscriptionWorker "github.com/diegobernardes/flare/domain/subscription/worker"
	infraWorker "github.com/diegobernardes/flare/infra/worker"
	memoryQueue "github.com/diegobernardes/flare/provider/memory/queue"
)

type worker struct {
	client                *Client
	base                  []infraWorker.Client
	workers               map[string]infraWorker.Client
	subscriptionPartition *subscriptionWorker.Partition
	subscriptionSpread    *subscriptionWorker.Spread
	subscriptionDelivery  *subscriptionWorker.Delivery
}

func (w *worker) init() error {
	key := "worker.enable"
	if w.client.config.IsSet(key) && !w.client.config.GetBool(key) {
		return nil
	}

	if err := w.initSubscriptionDelivery(); err != nil {
		return errors.Wrap(err, " error during subscription.delivery worker initialization")
	}

	if err := w.initSubscriptionSpread(); err != nil {
		return errors.Wrap(err, " error during subscription.spread worker initialization")
	}

	if err := w.initSubscriptionPartition(); err != nil {
		return errors.Wrap(err, "error during subscription.partition worker initialization")
	}

	for _, b := range w.base {
		b.Start()
	}

	return nil
}

func (w *worker) initSubscriptionPartition() error {
	concurrency := w.client.config.GetInt("worker.subscription.partition.concurrency")
	if concurrency == 0 {
		concurrency = 10
	} else if concurrency < 0 {
		return fmt.Errorf("invalid concurrency '%d' for worker subscription.partition", concurrency)
	}

	pusher, puller, err := w.initQueueMemorySubscriptionPartition()
	if err != nil {
		return errors.Wrap(err, "error during pusher and puller initialization")
	}

	unitOfWork := &subscriptionWorker.Partition{}

	client, err := infraWorker.NewClient(
		infraWorker.WorkerGoroutines(concurrency),
		infraWorker.WorkerLogger(w.client.logger),
		infraWorker.WorkerProcessor(unitOfWork),
		infraWorker.WorkerPuller(puller),
		infraWorker.WorkerPusher(pusher),
	)
	if err != nil {
		return errors.Wrap(err, "error during subscription.partition initialization")
	}

	err = unitOfWork.Init(
		subscriptionWorker.PartitionResourceRepository(w.client.repository.Resource()),
		subscriptionWorker.PartitionPusher(client),
		subscriptionWorker.PartitionOutput(w.subscriptionSpread),
	)
	if err != nil {
		return errors.Wrap(err, "error during worker processor initialization")
	}

	client.Start()
	w.subscriptionPartition = unitOfWork

	return nil
}

func (w *worker) initSubscriptionDelivery() error {
	concurrency := w.client.config.GetInt("worker.subscription.delivery.concurrency")
	if concurrency == 0 {
		concurrency = 100
	} else if concurrency < 0 {
		return fmt.Errorf("invalid concurrency '%d' for worker subscription.delivery", concurrency)
	}

	pusher, puller, err := w.initQueueMemorySubscriptionDelivery()
	if err != nil {
		return errors.Wrap(err, "error during pusher and puller initialization")
	}

	unitOfWork := &subscriptionWorker.Delivery{}

	client, err := infraWorker.NewClient(
		infraWorker.WorkerGoroutines(concurrency),
		infraWorker.WorkerLogger(w.client.logger),
		infraWorker.WorkerProcessor(unitOfWork),
		infraWorker.WorkerPuller(puller),
		infraWorker.WorkerPusher(pusher),
	)
	if err != nil {
		return errors.Wrap(err, "error during subscription.partition initialization")
	}

	err = unitOfWork.Init(
		subscriptionWorker.DeliveryPusher(client),
		subscriptionWorker.DeliverySubscriptionRepository(w.client.repository.Subscription()),
		subscriptionWorker.DeliveryHTTPClient(http.DefaultClient),
	)
	if err != nil {
		return errors.Wrap(err, "error during worker processor initialization")
	}

	client.Start()
	w.subscriptionDelivery = unitOfWork

	return nil
}

func (w *worker) initSubscriptionSpread() error {
	concurrency := w.client.config.GetInt("worker.subscription.spread.concurrency")
	if concurrency == 0 {
		concurrency = 10
	} else if concurrency < 0 {
		return fmt.Errorf("invalid concurrency '%d' for worker subscription.partition", concurrency)
	}

	pusher, puller, err := w.initQueueMemorySubscriptionSpread()
	if err != nil {
		return errors.Wrap(err, "error during pusher and puller initialization")
	}

	unitOfWork := &subscriptionWorker.Spread{}

	client, err := infraWorker.NewClient(
		infraWorker.WorkerGoroutines(concurrency),
		infraWorker.WorkerLogger(w.client.logger),
		infraWorker.WorkerProcessor(unitOfWork),
		infraWorker.WorkerPuller(puller),
		infraWorker.WorkerPusher(pusher),
	)
	if err != nil {
		return errors.Wrap(err, "error during subscription.partition initialization")
	}

	err = unitOfWork.Init(
		// TODO: faltando a concorrencia...
		subscriptionWorker.SpreadSubscriptionRepository(w.client.repository.Subscription()),
		subscriptionWorker.SpreadPusher(client),
		subscriptionWorker.SpreadOutput(w.subscriptionDelivery),
	)
	if err != nil {
		return errors.Wrap(err, "error during worker processor initialization")
	}

	client.Start()
	w.subscriptionSpread = unitOfWork

	return nil
}

func (w *worker) initQueue() error {
	provider := w.client.config.GetString("queue.provider")
	switch provider {
	case "memory":
	case "aws.sqs":
	default:
		return fmt.Errorf("unknown queue.provider '%s'", provider)
	}
	// buscar as queues do config e retornar aqui
	// falta fazer um metodo para pegar os timeouts

	return nil
}

func (w *worker) initWorkerQueueMemory() {
	// [provider]
	// [provider.memory.queue.subscription.partition]
	//   [provider.memory.queue.subscription.partition.ingress]
	//   [provider.memory.queue.subscription.partition.egress]
	//   [provider.memory.queue.subscription.partition.process]
}

func (w *worker) initQueueMemorySubscriptionPartition() (
	infraWorker.Pusher, infraWorker.Puller, error,
) {
	timeout, err := w.client.config.GetDuration(
		"provider.memory.queue.subscription.partition.process-timeout",
	)
	if err != nil {
		return nil, nil, err
	}

	queue := memoryQueue.NewClient(memoryQueue.ClientProcessTimeout(timeout))
	return queue, queue, nil
}

func (w *worker) initQueueMemorySubscriptionSpread() (
	infraWorker.Pusher, infraWorker.Puller, error,
) {
	timeout, err := w.client.config.GetDuration(
		"provider.memory.queue.subscription.spread.process-timeout",
	)
	if err != nil {
		return nil, nil, err
	}

	queue := memoryQueue.NewClient(memoryQueue.ClientProcessTimeout(timeout))
	return queue, queue, nil
}

func (w *worker) initQueueMemorySubscriptionDelivery() (
	infraWorker.Pusher, infraWorker.Puller, error,
) {
	timeout, err := w.client.config.GetDuration(
		"provider.memory.queue.subscription.delivery.process-timeout",
	)
	if err != nil {
		return nil, nil, err
	}

	queue := memoryQueue.NewClient(memoryQueue.ClientProcessTimeout(timeout))
	return queue, queue, nil
}

func (w *worker) stop() error {
	for _, b := range w.base {
		b.Stop()
	}

	return nil
}
