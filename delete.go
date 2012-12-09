// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DeleteService struct {
	client *Client
	index  string
	_type  string
	id     string
	routing string
	refresh *bool
	version *int
}

func NewDeleteService(client *Client) *DeleteService {
	builder := &DeleteService{
		client: client,
	}
	return builder
}

func (s *DeleteService) Index(index string) *DeleteService {
	s.index = index
	return s
}

func (s *DeleteService) Type(_type string) *DeleteService {
	s._type = _type
	return s
}

func (s *DeleteService) Id(id string) *DeleteService {
	s.id = id
	return s
}

func (s *DeleteService) Parent(parent string) *DeleteService {
	if s.routing == "" {
		s.routing = parent
	}
	return s
}

func (s *DeleteService) Refresh(refresh bool) *DeleteService {
	s.refresh = &refresh
	return s
}

func (s *DeleteService) Version(version int) *DeleteService {
	s.version = &version
	return s
}

func (s *DeleteService) Do() (*DeleteResult, error) {
	// Build url
	urls := "/{index}/{type}/{id}"
	urls = strings.Replace(urls, "{index}", cleanPathString(s.index), 1)
	urls = strings.Replace(urls, "{type}", cleanPathString(s._type), 1)
	urls = strings.Replace(urls, "{id}", cleanPathString(s.id), 1)

	// Parameters
	params := make(url.Values)
	if s.refresh != nil {
		params.Set("refresh", fmt.Sprintf("%v", *s.refresh))
	}
	if s.version != nil {
		params.Set("version", fmt.Sprintf("%d", *s.version))
	}
	if s.routing != "" {
		params.Set("routing", fmt.Sprintf("%s", s.routing))
	}
	urls += "?" + params.Encode()

	// Set up a new request
	req, err := s.client.NewRequest("DELETE", urls)
	if err != nil {
		return nil, err
	}

	// Get response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	ret := new(DeleteResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a delete request.

type DeleteResult struct {
	Ok  bool `json:"ok"`
	Found bool `json:"found"`
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id"`
	Version int64 `json:"_version"`
}
