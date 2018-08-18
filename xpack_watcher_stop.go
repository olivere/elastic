// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"net/url"
)

// XpackWatcherStopService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-stop.html.
type XpackWatcherStopService struct {
	client *Client
	pretty bool
}

// NewXpackWatcherStopService creates a new XpackWatcherStopService.
func NewXpackWatcherStopService(client *Client) *XpackWatcherStopService {
	return &XpackWatcherStopService{
		client: client,
	}
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherStopService) Pretty(pretty bool) *XpackWatcherStopService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherStopService) buildURL() (string, url.Values, error) {
	// Build URL path
	path := "/_xpack/watcher/_stop"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherStopService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *XpackWatcherStopService) Do(ctx context.Context) (*XpackWatcherStopResponse, error) {
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
	ret := new(XpackWatcherStopResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherStopResponse is the response of XpackWatcherStopService.Do.
type XpackWatcherStopResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
