// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// XpackWatcherDeleteWatchService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-delete-watch.html.
type XpackWatcherDeleteWatchService struct {
	client        *Client
	pretty        bool
	id            string
	masterTimeout string
}

// NewXpackWatcherDeleteWatchService creates a new XpackWatcherDeleteWatchService.
func NewXpackWatcherDeleteWatchService(client *Client) *XpackWatcherDeleteWatchService {
	return &XpackWatcherDeleteWatchService{
		client: client,
	}
}

// Id is documented as: Watch ID.
func (s *XpackWatcherDeleteWatchService) Id(id string) *XpackWatcherDeleteWatchService {
	s.id = id
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *XpackWatcherDeleteWatchService) MasterTimeout(masterTimeout string) *XpackWatcherDeleteWatchService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherDeleteWatchService) Pretty(pretty bool) *XpackWatcherDeleteWatchService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherDeleteWatchService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/watch/{id}", map[string]string{
		"id": s.id,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherDeleteWatchService) Validate() error {
	var invalid []string
	if s.id == "" {
		invalid = append(invalid, "Id")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XpackWatcherDeleteWatchService) Do(ctx context.Context) (*XpackWatcherDeleteWatchResponse, error) {
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
		Method: "DELETE",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XpackWatcherDeleteWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherDeleteWatchResponse is the response of XpackWatcherDeleteWatchService.Do.
type XpackWatcherDeleteWatchResponse struct {
	Found   bool   `json:"found"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
}
