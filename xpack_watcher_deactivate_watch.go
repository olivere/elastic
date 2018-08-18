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

// XpackWatcherDeactivateWatchService is documented at https://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-deactivate-watch.html.
type XpackWatcherDeactivateWatchService struct {
	client        *Client
	pretty        bool
	watchId       string
	masterTimeout string
	bodyJson      interface{}
	bodyString    string
}

// NewXpackWatcherDeactivateWatchService creates a new XpackWatcherDeactivateWatchService.
func NewXpackWatcherDeactivateWatchService(client *Client) *XpackWatcherDeactivateWatchService {
	return &XpackWatcherDeactivateWatchService{
		client: client,
	}
}

// WatchId is documented as: Watch ID.
func (s *XpackWatcherDeactivateWatchService) WatchId(watchId string) *XpackWatcherDeactivateWatchService {
	s.watchId = watchId
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *XpackWatcherDeactivateWatchService) MasterTimeout(masterTimeout string) *XpackWatcherDeactivateWatchService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherDeactivateWatchService) Pretty(pretty bool) *XpackWatcherDeactivateWatchService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: Execution control.
func (s *XpackWatcherDeactivateWatchService) BodyJson(body interface{}) *XpackWatcherDeactivateWatchService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: Execution control.
func (s *XpackWatcherDeactivateWatchService) BodyString(body string) *XpackWatcherDeactivateWatchService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherDeactivateWatchService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/watch/{watch_id}/_deactivate", map[string]string{
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
func (s *XpackWatcherDeactivateWatchService) Validate() error {
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
func (s *XpackWatcherDeactivateWatchService) Do(ctx context.Context) (*XpackWatcherDeactivateWatchResponse, error) {
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
	ret := new(XpackWatcherDeactivateWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherDeactivateWatchResponse is the response of XpackWatcherDeactivateWatchService.Do.
type XpackWatcherDeactivateWatchResponse struct {
	Status WatchStatus `json:"status"`
}
