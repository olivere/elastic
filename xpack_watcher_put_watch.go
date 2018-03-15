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

// XpackWatcherPutWatchService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-put-watch.html.
type XpackWatcherPutWatchService struct {
	client        *Client
	pretty        bool
	id            string
	active        *bool
	masterTimeout string
	bodyJson      interface{}
	bodyString    string
}

// NewXpackWatcherPutWatchService creates a new XpackWatcherPutWatchService.
func NewXpackWatcherPutWatchService(client *Client) *XpackWatcherPutWatchService {
	return &XpackWatcherPutWatchService{
		client: client,
	}
}

// Id is documented as: Watch ID.
func (s *XpackWatcherPutWatchService) Id(id string) *XpackWatcherPutWatchService {
	s.id = id
	return s
}

// Active is documented as: Specify whether the watch is in/active by default.
func (s *XpackWatcherPutWatchService) Active(active bool) *XpackWatcherPutWatchService {
	s.active = &active
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *XpackWatcherPutWatchService) MasterTimeout(masterTimeout string) *XpackWatcherPutWatchService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherPutWatchService) Pretty(pretty bool) *XpackWatcherPutWatchService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: The watch.
func (s *XpackWatcherPutWatchService) BodyJson(body interface{}) *XpackWatcherPutWatchService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: The watch.
func (s *XpackWatcherPutWatchService) BodyString(body string) *XpackWatcherPutWatchService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherPutWatchService) buildURL() (string, url.Values, error) {
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
	if s.active != nil {
		params.Set("active", fmt.Sprintf("%v", *s.active))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherPutWatchService) Validate() error {
	var invalid []string
	if s.id == "" {
		invalid = append(invalid, "Id")
	}
	if s.bodyString == "" && s.bodyJson == nil {
		invalid = append(invalid, "BodyJson")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XpackWatcherPutWatchService) Do(ctx context.Context) (*XpackWatcherPutWatchResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request body
	var body interface{}
	if s.bodyJson != nil {
		body = s.bodyJson
	} else {
		body = s.bodyString
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "PUT",
		Path:   path,
		Params: params,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XpackWatcherPutWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherPutWatchResponse is the response of XpackWatcherPutWatchService.Do.
type XpackWatcherPutWatchResponse struct {
}
