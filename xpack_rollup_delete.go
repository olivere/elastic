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

// XPackRollupDeleteService delete a rollup job by its job id.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/rollup-delete-job.html.
type XPackRollupDeleteService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	jobId string
}

// XPackRollupDeleteService creates a new XPackRollupDeleteService.
func NewXPackRollupDeleteService(client *Client) *XPackRollupDeleteService {
	return &XPackRollupDeleteService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *XPackRollupDeleteService) Pretty(pretty bool) *XPackRollupDeleteService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *XPackRollupDeleteService) Human(human bool) *XPackRollupDeleteService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *XPackRollupDeleteService) ErrorTrace(errorTrace bool) *XPackRollupDeleteService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *XPackRollupDeleteService) FilterPath(filterPath ...string) *XPackRollupDeleteService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *XPackRollupDeleteService) Header(name string, value string) *XPackRollupDeleteService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *XPackRollupDeleteService) Headers(headers http.Header) *XPackRollupDeleteService {
	s.headers = headers
	return s
}

// JobId is id of the rollup to delete.
func (s *XPackRollupDeleteService) JobId(jobId string) *XPackRollupDeleteService {
	s.jobId = jobId
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackRollupDeleteService) buildURL() (string, url.Values, error) {
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
func (s *XPackRollupDeleteService) Validate() error {
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
func (s *XPackRollupDeleteService) Do(ctx context.Context) (*XPackRollupDeleteResponse, error) {
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
		Method:  "DELETE",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XPackRollupDeleteResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XPackRollupDeleteResponse is the response of XPackRollupDeleteService.Do.
type XPackRollupDeleteResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
