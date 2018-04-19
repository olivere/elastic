// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"net/url"
)

// CatIndicesService executes cat commands against the cluster.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.2/cat.html
// for details.
type CatService struct {
	client *Client
	pretty bool
	model  string
}

// NewCatService creates a new CatService.
func NewCatService(client *Client) *CatService {
	return &CatService{
		client: client,
		pretty: true,
	}
}

func (s *CatService) Indices() *CatService {
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatService) Pretty(pretty bool) *CatService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *CatService) buildURL() (string, url.Values, error) {
	path := "/_cat/indices"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *CatService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatService) Do(ctx context.Context) (string, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return "", err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return "", err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "GET",
		Path:   path,
		Params: params,
	})

	if err != nil {
		return "", err
	}

	return string(res.Body), nil
}

// CatResponse is the response of CatService.Do.
type CatResponse struct {
	Response string
}
