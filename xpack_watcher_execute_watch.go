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

// XpackWatcherExecuteWatchService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-execute-watch.html.
type XpackWatcherExecuteWatchService struct {
	client     *Client
	pretty     bool
	id         string
	debug      *bool
	bodyJson   interface{}
	bodyString string
}

// NewXpackWatcherExecuteWatchService creates a new XpackWatcherExecuteWatchService.
func NewXpackWatcherExecuteWatchService(client *Client) *XpackWatcherExecuteWatchService {
	return &XpackWatcherExecuteWatchService{
		client: client,
	}
}

// Id is documented as: Watch ID.
func (s *XpackWatcherExecuteWatchService) Id(id string) *XpackWatcherExecuteWatchService {
	s.id = id
	return s
}

// Debug is documented as: indicates whether the watch should execute in debug mode.
func (s *XpackWatcherExecuteWatchService) Debug(debug bool) *XpackWatcherExecuteWatchService {
	s.debug = &debug
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherExecuteWatchService) Pretty(pretty bool) *XpackWatcherExecuteWatchService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: Execution control.
func (s *XpackWatcherExecuteWatchService) BodyJson(body interface{}) *XpackWatcherExecuteWatchService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: Execution control.
func (s *XpackWatcherExecuteWatchService) BodyString(body string) *XpackWatcherExecuteWatchService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherExecuteWatchService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/watch/{id}/_execute", map[string]string{
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
	if s.debug != nil {
		params.Set("debug", fmt.Sprintf("%v", *s.debug))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherExecuteWatchService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *XpackWatcherExecuteWatchService) Do(ctx context.Context) (*XpackWatcherExecuteWatchResponse, error) {
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
		Method: "POST",
		Path:   path,
		Params: params,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XpackWatcherExecuteWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherExecuteWatchResponse is the response of XpackWatcherExecuteWatchService.Do.
type XpackWatcherExecuteWatchResponse struct {
	Id          string      `json:"_id"`
	WatchRecord WatchRecord `json:"watch_record"`
}

type WatchRecord struct {
	WatchId   string                            `json:"watch_id"`
	Node      string                            `json:"node"`
	Messages  []string                          `json:"messages"`
	State     string                            `json:"state"`
	Status    WatchRecordStatus                 `json:"status"`
	Input     map[string]map[string]interface{} `json:"input"`
	Condition map[string]map[string]interface{} `json:"condition"`
	Result    map[string]interface{}            `json:"Result"`
}

type WatchRecordStatus struct {
	Version          int                               `json:"version"`
	State            map[string]interface{}            `json:"state"`
	LastChecked      string                            `json:"last_checked"`
	LastMetCondition string                            `json:"last_met_condition"`
	Actions          map[string]map[string]interface{} `json:"actions"`
	ExecutionState   string                            `json:"execution_state"`
}
