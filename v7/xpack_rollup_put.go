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

// XPackRollupPutService create or update a rollup job by its job id.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/rollup-put-job.html.
type XPackRollupPutService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	jobId string
	body  interface{}
}

// NewXPackRollupPutService creates a new XPackRollupPutService.
func NewXPackRollupPutService(client *Client) *XPackRollupPutService {
	return &XPackRollupPutService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *XPackRollupPutService) Pretty(pretty bool) *XPackRollupPutService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *XPackRollupPutService) Human(human bool) *XPackRollupPutService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *XPackRollupPutService) ErrorTrace(errorTrace bool) *XPackRollupPutService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *XPackRollupPutService) FilterPath(filterPath ...string) *XPackRollupPutService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *XPackRollupPutService) Header(name string, value string) *XPackRollupPutService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *XPackRollupPutService) Headers(headers http.Header) *XPackRollupPutService {
	s.headers = headers
	return s
}

// JobId is id of the rollup to create.
func (s *XPackRollupPutService) JobId(jobId string) *XPackRollupPutService {
	s.jobId = jobId
	return s
}

// Body specifies the role. Use a string or a type that will get serialized as JSON.
func (s *XPackRollupPutService) Body(body interface{}) *XPackRollupPutService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackRollupPutService) buildURL() (string, url.Values, error) {
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
func (s *XPackRollupPutService) Validate() error {
	var invalid []string
	if s.jobId == "" {
		invalid = append(invalid, "Job ID")
	}
	if s.body == nil {
		invalid = append(invalid, "Body")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XPackRollupPutService) Do(ctx context.Context) (*XPackRollupPutResponse, error) {
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
		Method:  "PUT",
		Path:    path,
		Params:  params,
		Body:    s.body,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XPackRollupPutResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XPackRollupPutResponse is the response of XPackRollupPutService.Do.
type XPackRollupPutResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
