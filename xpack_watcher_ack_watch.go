// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// XpackWatcherAckWatchService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-ack-watch.html.
type XpackWatcherAckWatchService struct {
	client        *Client
	pretty        bool
	watchId       string
	actionId      []string
	masterTimeout string
	bodyJson      interface{}
	bodyString    string
}

// NewXpackWatcherAckWatchService creates a new XpackWatcherAckWatchService.
func NewXpackWatcherAckWatchService(client *Client) *XpackWatcherAckWatchService {
	return &XpackWatcherAckWatchService{
		client:   client,
		actionId: make([]string, 0),
	}
}

// WatchId is documented as: Watch ID.
func (s *XpackWatcherAckWatchService) WatchId(watchId string) *XpackWatcherAckWatchService {
	s.watchId = watchId
	return s
}

// ActionId is documented as: A comma-separated list of the action ids to be acked.
func (s *XpackWatcherAckWatchService) ActionId(actionId []string) *XpackWatcherAckWatchService {
	s.actionId = actionId
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *XpackWatcherAckWatchService) MasterTimeout(masterTimeout string) *XpackWatcherAckWatchService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherAckWatchService) Pretty(pretty bool) *XpackWatcherAckWatchService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: Execution control.
func (s *XpackWatcherAckWatchService) BodyJson(body interface{}) *XpackWatcherAckWatchService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: Execution control.
func (s *XpackWatcherAckWatchService) BodyString(body string) *XpackWatcherAckWatchService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherAckWatchService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/watch/{watch_id}/_ack", map[string]string{
		"action_id": strings.Join(s.actionId, ","),
		"watch_id":  s.watchId,
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
func (s *XpackWatcherAckWatchService) Validate() error {
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
func (s *XpackWatcherAckWatchService) Do(ctx context.Context) (*XpackWatcherAckWatchResponse, error) {
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
	ret := new(XpackWatcherAckWatchResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherAckWatchResponse is the response of XpackWatcherAckWatchService.Do.
type XpackWatcherAckWatchResponse struct {
	Status AckWatchStatus `json:"status"`
}

type AckWatchStatus struct {
	State            map[string]interface{}            `json:"state"`
	LastChecked      string                            `json:"last_checked"`
	LastMetCondition string                            `json:"last_met_condition"`
	Actions          map[string]map[string]interface{} `json:"actions"`
	ExecutionState   string                            `json:"execution_state"`
	Version          int                               `json:"version"`
}
