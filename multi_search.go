// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
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

// MultiSearch executes one or more searches in one roundtrip.
// See http://www.elasticsearch.org/guide/reference/api/multi-search/
type MultiSearchService struct {
	client     *Client
	requests   []*SearchRequest
	indices    []string
	pretty     bool
	debug      bool
	routing    string
	preference string
}

func NewMultiSearchService(client *Client) *MultiSearchService {
	builder := &MultiSearchService{
		client:   client,
		requests: make([]*SearchRequest, 0),
		indices:  make([]string, 0),
		debug:    false,
		pretty:   false,
	}
	return builder
}

func (s *MultiSearchService) Add(requests ...*SearchRequest) *MultiSearchService {
	s.requests = append(s.requests, requests...)
	return s
}

func (s *MultiSearchService) Index(index string) *MultiSearchService {
	s.indices = append(s.indices, index)
	return s
}

func (s *MultiSearchService) Indices(indices ...string) *MultiSearchService {
	s.indices = append(s.indices, indices...)
	return s
}

func (s *MultiSearchService) Pretty(pretty bool) *MultiSearchService {
	s.pretty = pretty
	return s
}

func (s *MultiSearchService) Debug(debug bool) *MultiSearchService {
	s.debug = debug
	return s
}

func (s *MultiSearchService) Do() (*MultiSearchResult, error) {
	// Build url
	urls := "/_msearch"

	// Parameters
	params := make(url.Values)
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("GET", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	lines := make([]string, 0)
	for _, sr := range s.requests {
		// Set default indices if not specified in the request
		if !sr.HasIndices() && len(s.indices) > 0 {
			sr = sr.Indices(s.indices...)
		}

		header, err := json.Marshal(sr.header())
		if err != nil {
			return nil, err
		}
		body, err := json.Marshal(sr.body())
		if err != nil {
			return nil, err
		}
		lines = append(lines, string(header))
		lines = append(lines, string(body))
	}
	req.SetBodyString(strings.Join(lines, "\n") + "\n") // Don't forget trailing \n

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

	ret := new(MultiSearchResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type MultiSearchResult struct {
	Responses []*SearchResult `json:"responses,omitempty"`
}
