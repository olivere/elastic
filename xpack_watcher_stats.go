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

// XpackWatcherStatsService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-stats.html.
type XpackWatcherStatsService struct {
	client          *Client
	pretty          bool
	metric          string
	emitStacktraces *bool
}

// NewXpackWatcherStatsService creates a new XpackWatcherStatsService.
func NewXpackWatcherStatsService(client *Client) *XpackWatcherStatsService {
	return &XpackWatcherStatsService{
		client: client,
	}
}

// Metric is documented as: Controls what additional stat metrics should be include in the response.
func (s *XpackWatcherStatsService) Metric(metric string) *XpackWatcherStatsService {
	s.metric = metric
	return s
}

// EmitStacktraces is documented as: Emits stack traces of currently running watches.
func (s *XpackWatcherStatsService) EmitStacktraces(emitStacktraces bool) *XpackWatcherStatsService {
	s.emitStacktraces = &emitStacktraces
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XpackWatcherStatsService) Pretty(pretty bool) *XpackWatcherStatsService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XpackWatcherStatsService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/watcher/stats", map[string]string{
		"metric": s.metric,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.emitStacktraces != nil {
		params.Set("emit_stacktraces", fmt.Sprintf("%v", *s.emitStacktraces))
	}
	if s.metric != "" {
		params.Set("metric", s.metric)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XpackWatcherStatsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *XpackWatcherStatsService) Do(ctx context.Context) (*XpackWatcherStatsResponse, error) {
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
		Method: "GET",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XpackWatcherStatsResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XpackWatcherStatsResponse is the response of XpackWatcherStatsService.Do.
type XpackWatcherStatsResponse struct {
	Stats []WatcherStats `json:"stats"`
}

type WatcherStats struct {
	WatcherState        string                 `json:"watcher_state"`
	WatchCount          int                    `json:"watch_count"`
	ExecutionThreadPool map[string]interface{} `json:"execution_thread_pool"`
}
