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

// XPackRollupStartService starts the rollup job if it is not already running.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/rollup-start-job.html.
type XPackRollupStartService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	jobId string
}

// NewXPackRollupStartService creates a new XPackRollupStartService.
func NewXPackRollupStartService(client *Client) *XPackRollupStartService {
	return &XPackRollupStartService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *XPackRollupStartService) Pretty(pretty bool) *XPackRollupStartService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *XPackRollupStartService) Human(human bool) *XPackRollupStartService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *XPackRollupStartService) ErrorTrace(errorTrace bool) *XPackRollupStartService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *XPackRollupStartService) FilterPath(filterPath ...string) *XPackRollupStartService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *XPackRollupStartService) Header(name string, value string) *XPackRollupStartService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *XPackRollupStartService) Headers(headers http.Header) *XPackRollupStartService {
	s.headers = headers
	return s
}

// JobId is id of the rollup to retrieve.
func (s *XPackRollupStartService) JobId(jobId string) *XPackRollupStartService {
	s.jobId = jobId
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackRollupStartService) buildURL() (string, url.Values, error) {
	// Build URL path
	path, err := uritemplates.Expand("/_rollup/job/{job_id}/_start", map[string]string{
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
func (s *XPackRollupStartService) Validate() error {
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
func (s *XPackRollupStartService) Do(ctx context.Context) (*XPackRollupStartResponse, error) {
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
		Method:  "POST",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XPackRollupStartResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XPackRollupStartResponse is the response of XPackRollupStartService.Do.
type XPackRollupStartResponse struct {
	Started bool `json:"started"`
}
