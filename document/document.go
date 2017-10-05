// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package document

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/diegobernardes/flare"
)

type document flare.Document

func (d *document) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id               string      `json:"id"`
		ChangeFieldValue interface{} `json:"changeFieldValue"`
		UpdatedAt        string      `json:"updatedAt"`
	}{
		Id:               d.Id,
		ChangeFieldValue: d.ChangeFieldValue,
		UpdatedAt:        d.UpdatedAt.Format(time.RFC3339),
	})
}

func transformDocument(d *flare.Document) *document {
	return (*document)(d)
}

type response struct {
	Error    *responseError `json:"error,omitempty"`
	Document *document      `json:"document,omitempty"`
}

type responseError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (s *Service) writeError(w http.ResponseWriter, err error, title string, status int) {
	resp := &response{Error: &responseError{Status: status}}

	if err != nil {
		resp.Error.Detail = err.Error()
	}

	if title != "" {
		resp.Error.Title = title
	}

	s.writeResponse(w, resp, status, nil)
}