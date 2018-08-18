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

// XpackWatcherActivateWatchService is documented at https://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-activate-watch.html.
type XpackWatcherActivateWatchService struct {
	client        *Client
	pretty        bool
	watchId       string
	masterTimeout string
	bodyJson      interface{}
	bodyString    string
}

// NewXpackWatcherActivateWatchService creates a new XpackWatcherActivateWatchService.
func NewXpackWatcherActivateWatchService(client *Client) *XpackWatcherActivateWatchService {
	return &XpackWatcherActivateWatchService{
		client: client,
	}
}

// WatchId is documented as: Watch ID.
func (s *XpackWatcherActivateWatchService) WatchId(watchId string) *XpackWatcherActivateWatchService {
	s.watchId = watchId
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *XpackWatcherActivateWatchService) MasterTimeout(masterTimeout string) *XpackWatcherActivateWatchService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherActivateWatchService) Pretty(pretty bool) *XpackWatcherActivateWatchService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: Execution control.
func (s *XpackWatcherActivateWatchService) BodyJson(body interface{}) *XpackWatcherActivateWatchService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: Execution control.
func (s *XpackWatcherActivateWatchService) BodyString(body string) *XpackWatcherActivateWatchService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherActivateWatchService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/watch/{watch_id}/_activate", map[string]string{
		"watch_id": s.watchId,
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
func (s *XpackWatcherActivateWatchService) Validate() error {
	var invalid []string
	if s.watchId == "" {
		invalid = append(invalid, "WatchId")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XpackWatcherActivateWatchService) Do(ctx context.Context) (*XpackWatcherActivateWatchResponse, error) {
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
	ret := new(XpackWatcherActivateWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherActivateWatchResponse is the response of XpackWatcherActivateWatchService.Do.
type XpackWatcherActivateWatchResponse struct {
	Status WatchStatus `json:"status"`
}
