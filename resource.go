// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// Resource represents the apis Flare track and the info to detect changes on documents.
type Resource struct {
	ID        string
	Addresses []string
	Path      string
	Change    ResourceChange
	CreatedAt time.Time
}

// The types of value Flare support to detect document change.
const (
	ResourceChangeInteger = "integer"
	ResourceChangeDate    = "date"
)

// ResourceChange holds the information to detect document change.
type ResourceChange struct {
	Field      string
	Kind       string
	DateFormat string
}

// Valid indicates if the current resourceChange is valid.
func (rc *ResourceChange) Valid() error {
	if rc.Field == "" {
		return errors.New("missing field")
	}

	if rc.Kind == "" {
		return errors.New("missing kind")
	}

	switch rc.Kind {
	case ResourceChangeDate:
		if rc.DateFormat == "" {
			return errors.New("missing dateFormat")
		}
	case ResourceChangeInteger:
		if rc.DateFormat != "" {
			return errors.New("dateFormat should not be present if the kind is integer")
		}
	}
	return nil
}

// ResourceRepositorier is used to interact with Resource repository.
type ResourceRepositorier interface {
	FindAll(context.Context, *Pagination) ([]Resource, *Pagination, error)
	FindOne(context.Context, string) (*Resource, error)
	FindByURI(context.Context, string) (*Resource, error)
	Create(context.Context, *Resource) error
	Delete(context.Context, string) error
}

// ResourceRepositoryError represents all the errors the repository can return.
type ResourceRepositoryError interface {
	AlreadyExists() bool
	PathConflict() bool
	NotFound() bool
}
