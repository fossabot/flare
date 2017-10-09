// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/diegobernardes/flare"
)

// Subscription implements the data layer for the subscription service.
type Subscription struct {
	client            *Client
	database          string
	collection        string
	collectionTrigger string
}

// FindAll returns a list of subscriptions.
func (s *Subscription) FindAll(
	_ context.Context, pagination *flare.Pagination, id string,
) ([]flare.Subscription, *flare.Pagination, error) {
	var (
		group         errgroup.Group
		subscriptions []flare.Subscription
		total         int
	)

	group.Go(func() error {
		session := s.client.session()
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()

		totalResult, err := session.DB(s.database).C(s.collection).Find(bson.M{}).Count()
		if err != nil {
			return err
		}
		total = totalResult
		return nil
	})

	group.Go(func() error {
		session := s.client.session()
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)

		q := session.
			DB(s.database).
			C(s.collection).
			Find(bson.M{}).
			Sort("createdAt").
			Limit(pagination.Limit)
		if pagination.Offset != 0 {
			q = q.Skip(pagination.Offset)
		}

		return q.All(&subscriptions)
	})

	if err := group.Wait(); err != nil {
		return nil, nil, errors.Wrap(err, "error during MongoDB access")
	}

	return subscriptions, &flare.Pagination{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Total:  total,
	}, nil
}

// FindOne return the Subscription that match the id.
func (s *Subscription) FindOne(
	_ context.Context, resourceId, id string,
) (*flare.Subscription, error) {
	session := s.client.session()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	result := &flare.Subscription{}
	err := session.DB(s.database).C(s.collection).Find(bson.M{"id": id}).One(result)
	if err == mgo.ErrNotFound {
		return nil, &errMemory{message: fmt.Sprintf(
			"subscription '%s' at resource '%s' not found", id, resourceId,
		), notFound: true}
	}
	return result, errors.Wrap(err, "error during subscription search")
}

// Create a subscription.
func (s *Subscription) Create(_ context.Context, subscription *flare.Subscription) error {
	session := s.client.session()
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	resourceEntity := &resourceEntity{}
	err := session.DB(s.database).C(s.collection).Find(bson.M{
		"resource.id":  subscription.Resource.Id,
		"endpoint.url": subscription.Endpoint.URL.String(),
	}).One(resourceEntity)
	if err == nil {
		return fmt.Errorf("already has a subscription '%s' with this endpoint", resourceEntity.Id)
	}
	if err != nil && err != mgo.ErrNotFound {
		return errors.Wrap(err, "error during subscription search")
	}

	subscription.CreatedAt = time.Now()
	return errors.Wrap(
		session.DB(s.database).C(s.collection).Insert(subscription),
		"error during subscription create",
	)
}

// HasSubscription check if a resource has subscriptions.
func (s *Subscription) HasSubscription(ctx context.Context, resourceId string) (bool, error) {
	session := s.client.session()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	count, err := session.
		DB(s.database).
		C(s.collection).
		Find(bson.M{"resource.id": resourceId}).
		Count()
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, &errMemory{message: fmt.Sprintf(
				"subscriptions not found for resource '%s'", resourceId), notFound: true,
			}
		}
		return false, err
	}
	return count > 0, nil
}

// Delete a given subscription.
func (s *Subscription) Delete(_ context.Context, resourceId, id string) error {
	session := s.client.session()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(s.database).C(s.collection)

	if err := c.Remove(bson.M{"id": id, "resource.id": resourceId}); err != nil {
		if err == mgo.ErrNotFound {
			return &errMemory{message: fmt.Sprintf(
				"subscription '%s' at resource '%s' not found", id, resourceId,
			), notFound: true}
		}
	}
	return nil
}

// Trigger process the update on a document.
func (s *Subscription) Trigger(
	ctx context.Context,
	kind string,
	doc *flare.Document,
	fn func(context.Context, flare.Subscription, string) error,
) error {
	session := s.client.session()
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var subscriptions []flare.Subscription
	err := session.
		DB(s.database).
		C(s.collection).
		Find(bson.M{"resource.id": doc.Resource.Id}).
		All(&subscriptions)
	if err != nil {
		return errors.Wrap(err, "error while subscription search")
	}

	group, groupCtx := errgroup.WithContext(ctx)
	for i := range subscriptions {
		group.Go(s.triggerProcess(groupCtx, subscriptions[i], doc, kind, fn))
	}

	return errors.Wrap(group.Wait(), "error during processing")
}

func (s *Subscription) loadReferenceDocument(
	session *mgo.Session,
	subs flare.Subscription,
	doc *flare.Document,
) (*flare.Document, error) {
	content := make(map[string]interface{})
	err := session.
		DB(s.database).
		C(s.collectionTrigger).
		Find(bson.M{"subscriptionId": subs.Id, "document.id": doc.Id}).
		One(&content)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error during search")
	}

	return &flare.Document{
		Id:               doc.Id,
		ChangeFieldValue: content["document"].(map[string]interface{})["changeFieldValue"],
	}, nil
}

func (s *Subscription) newEntry(
	groupCtx context.Context,
	kind string,
	session *mgo.Session,
	subs flare.Subscription,
	doc *flare.Document,
	fn func(context.Context, flare.Subscription, string) error,
) error {
	if kind == flare.SubscriptionTriggerDelete {
		return nil
	}

	err := s.upsertSubscriptionTrigger(session, subs, doc)
	if err != nil {
		return errors.Wrap(err, "error during document upsert")
	}

	if err = fn(groupCtx, subs, flare.SubscriptionTriggerCreate); err != nil {
		return errors.Wrap(err, "error during document subscription processing")
	}
	return nil
}

func (s *Subscription) triggerProcessDelete(
	groupCtx context.Context,
	kind string,
	session *mgo.Session,
	subs flare.Subscription,
	doc *flare.Document,
	fn func(context.Context, flare.Subscription, string) error,
) error {
	err := session.
		DB(s.database).
		C(s.collectionTrigger).
		Remove(bson.M{"subscriptionId": subs.Id, "document.id": doc.Id})
	if err != nil {
		return errors.Wrap(err, "error during subscriptionTriggers delete")
	}

	if err = fn(groupCtx, subs, flare.SubscriptionTriggerDelete); err != nil {
		return errors.Wrap(err, "error during document subscription processing")
	}
	return nil
}

func (s *Subscription) upsertSubscriptionTrigger(
	session *mgo.Session,
	subs flare.Subscription,
	doc *flare.Document,
) error {
	_, err := session.
		DB(s.database).
		C(s.collectionTrigger).
		Upsert(
			bson.M{"subscriptionId": subs.Id, "document.id": doc.Id},
			bson.M{"subscriptionId": subs.Id, "document": bson.M{
				"id":               doc.Id,
				"changeFieldValue": doc.ChangeFieldValue,
				"updatedAt":        time.Now(),
			}},
		)
	if err != nil {
		return errors.Wrap(err, "error during update subscriptionTriggers")
	}
	return nil
}

func (s *Subscription) triggerProcess(
	groupCtx context.Context,
	subs flare.Subscription,
	doc *flare.Document,
	kind string,
	fn func(context.Context, flare.Subscription, string) error,
) func() error {
	return func() error {
		session := s.client.session()
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()

		referenceDocument, err := s.loadReferenceDocument(session, subs, doc)
		if err != nil {
			return errors.Wrap(err, "error during reference document search")
		}

		if referenceDocument == nil {
			return s.newEntry(groupCtx, kind, session, subs, doc, fn)
		}

		if kind == flare.SubscriptionTriggerDelete {
			return s.triggerProcessDelete(groupCtx, kind, session, subs, doc, fn)
		}

		newer, err := doc.Newer(referenceDocument)
		if err != nil {
			return errors.Wrap(err, "error during check if document is newer")
		}
		if !newer {
			return nil
		}

		err = s.upsertSubscriptionTrigger(session, subs, doc)
		if err != nil {
			return errors.Wrap(err, "error during update subscriptionTriggers")
		}

		if err = fn(groupCtx, subs, flare.SubscriptionTriggerUpdate); err != nil {
			return errors.Wrap(err, "error during document subscription processing")
		}
		return nil
	}
}

// NewSubscription returns a configured subscription repository.
func NewSubscription(options ...func(*Subscription)) (*Subscription, error) {
	s := &Subscription{}
	for _, option := range options {
		option(s)
	}

	if s.client == nil {
		return nil, errors.New("invalid client")
	}

	s.collection = "subscriptions"
	s.collectionTrigger = "subscriptionTriggers"
	s.database = s.client.database
	return s, nil
}

// SubscriptionClient set the client to access MongoDB.
func SubscriptionClient(client *Client) func(*Subscription) {
	return func(s *Subscription) {
		s.client = client
	}
}
