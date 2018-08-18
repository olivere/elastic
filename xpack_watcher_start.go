// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"net/url"
)

// XpackWatcherStartService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-start.html.
type XpackWatcherStartService struct {
	client *Client
	pretty bool
}

// NewXpackWatcherStartService creates a new XpackWatcherStartService.
func NewXpackWatcherStartService(client *Client) *XpackWatcherStartService {
	return &XpackWatcherStartService{
		client: client,
	}
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherStartService) Pretty(pretty bool) *XpackWatcherStartService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherStartService) buildURL() (string, url.Values, error) {
	// Build URL path
	path := "/_xpack/watcher/_start"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherStartService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *XpackWatcherStartService) Do(ctx context.Context) (*XpackWatcherStartResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "POST",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XpackWatcherStartResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherStartResponse is the response of XpackWatcherStartService.Do.
type XpackWatcherStartResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
