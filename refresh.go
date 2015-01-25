// Copyright 2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

type RefreshService struct {
	client  *Client
	indices []string
	force   *bool
	pretty  bool
	debug   bool
}

func NewRefreshService(client *Client) *RefreshService {
	builder := &RefreshService{
		client:  client,
		indices: make([]string, 0),
	}
	return builder
}

func (s *RefreshService) Index(index string) *RefreshService {
	s.indices = append(s.indices, index)
	return s
}

func (s *RefreshService) Indices(indices ...string) *RefreshService {
	s.indices = append(s.indices, indices...)
	return s
}

func (s *RefreshService) Force(force bool) *RefreshService {
	s.force = &force
	return s
}

func (s *RefreshService) Pretty(pretty bool) *RefreshService {
	s.pretty = pretty
	return s
}

func (s *RefreshService) Debug(debug bool) *RefreshService {
	s.debug = debug
	return s
}

func (s *RefreshService) Do() (*RefreshResult, error) {
	// Build url
	urls := "/"

	// Indices part
	indexPart := make([]string, 0)
	for _, index := range s.indices {
		index, err := uritemplates.Expand("{index}", map[string]string{
			"index": index,
		})
		if err != nil {
			return nil, err
		}
		indexPart = append(indexPart, index)
	}
	if len(indexPart) > 0 {
		urls += strings.Join(indexPart, ",")
	}

	urls += "/_refresh"

	// Parameters
	params := make(url.Values)
	if s.force != nil {
		params.Set("force", fmt.Sprintf("%v", *s.force))
	}
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	if s.debug {
		s.client.dumpRequest((*http.Request)(req))
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

	if s.debug {
		s.client.dumpResponse(res)
	}

	ret := new(RefreshResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a refresh request.

type RefreshResult struct {
	Shards shardsInfo `json:"_shards,omitempty"`
}
