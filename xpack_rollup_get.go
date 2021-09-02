// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/v7/uritemplates"
)

// XPackRollupGetService retrieves a role by its name.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/rollup-apis.html.
type XPackRollupGetService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	jobId string
}

// NewXPackRollupGetService creates a new XPackRollupGetService.
func NewXPackRollupGetService(client *Client) *XPackRollupGetService {
	return &XPackRollupGetService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *XPackRollupGetService) Pretty(pretty bool) *XPackRollupGetService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *XPackRollupGetService) Human(human bool) *XPackRollupGetService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *XPackRollupGetService) ErrorTrace(errorTrace bool) *XPackRollupGetService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *XPackRollupGetService) FilterPath(filterPath ...string) *XPackRollupGetService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *XPackRollupGetService) Header(name string, value string) *XPackRollupGetService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *XPackRollupGetService) Headers(headers http.Header) *XPackRollupGetService {
	s.headers = headers
	return s
}

// JobId is id of the rollup to retrieve.
func (s *XPackRollupGetService) JobId(jobId string) *XPackRollupGetService {
	s.jobId = jobId
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackRollupGetService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_rollup/job/{job_id}", map[string]string{
		"job_id": s.jobId,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XPackRollupGetService) Validate() error {
	var invalid []string
	if s.jobId == "" {
		invalid = append(invalid, "Job ID")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XPackRollupGetService) Do(ctx context.Context) (*XPackRollupGetResponse, error) {
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
		Method:  "GET",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := XPackRollupGetResponse{}
	if err := json.Unmarshal(res.Body, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// XPackRollupGetResponse is the response of XPackRollupGetService.Do.
type XPackRollupGetResponse struct {
	Jobs []XPackRollup `json:"jobs"`
}

// XPackRollup is the role object.
type XPackRollup struct {
	Config XPackRollupConfig `json:"config"`
	Status XPackRollupStatus `json:"status"`
	Stats  XPackRollupStats  `json:"stats"`
}

type XPackRollupConfig struct {
	Id           string                 `json:"id"`
	Cron         string                 `json:"cron"`
	IndexPattern string                 `json:"index_pattern"`
	RollupIndex  string                 `json:"rollup_index"`
	Groups       map[string]interface{} `json:"groups"`
	Metrics      []XPackRollupMetrics   `json:"metrics"`
	Timeout      string                 `json:"timeout"`
	PageSize     int                    `json:"page_size"`
}

type XPackRollupMetrics struct {
	Field   string   `json:"field"`
	Metrics []string `json:"metrics"`
}

type XPackRollupStatus struct {
	JobState      string `json:"job_state"`
	UpgradedDocId bool   `json:"upgraded_doc_id"`
}

type XPackRollupStats struct {
	PageProcessed      int `json:"pages_processed"`
	DocumentsProcessed int `json:"documents_processed"`
	RollupsIndexed     int `json:"rollups_indexed"`
	TriggerCount       int `json:"trigger_count"`
	IndexFailures      int `json:"index_failures"`
	IndexTimeInMs      int `json:"index_time_in_ms"`
	IndexTotal         int `json:"index_total"`
	SearchFailures     int `json:"search_failures"`
	SearchTimeInMs     int `json:"search_time_in_ms"`
	SearchTotal        int `json:"search_total"`
	ProcessingTimeInMs int `json:"processing_time_in_ms"`
	ProcessingTotal    int `json:"processing_total"`
}
