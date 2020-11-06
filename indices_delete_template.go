// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/v7/uritemplates"
)

// IndicesDeleteTemplateService deletes templates.
//
// Index templates have changed during in 7.8 update of Elasticsearch.
// This service implements the legacy version (7.7 or lower). If you want
// the new version, please use the IndicesDeleteIndexTemplateService.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.9/indices-delete-template-v1.html
// for more details.
type IndicesDeleteTemplateService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name          string
	timeout       string
	masterTimeout string
}

// NewIndicesDeleteTemplateService creates a new IndicesDeleteTemplateService.
func NewIndicesDeleteTemplateService(client *Client) *IndicesDeleteTemplateService {
	return &IndicesDeleteTemplateService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *IndicesDeleteTemplateService) Pretty(pretty bool) *IndicesDeleteTemplateService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *IndicesDeleteTemplateService) Human(human bool) *IndicesDeleteTemplateService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *IndicesDeleteTemplateService) ErrorTrace(errorTrace bool) *IndicesDeleteTemplateService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *IndicesDeleteTemplateService) FilterPath(filterPath ...string) *IndicesDeleteTemplateService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *IndicesDeleteTemplateService) Header(name string, value string) *IndicesDeleteTemplateService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *IndicesDeleteTemplateService) Headers(headers http.Header) *IndicesDeleteTemplateService {
	s.headers = headers
	return s
}

// Name is the name of the template.
func (s *IndicesDeleteTemplateService) Name(name string) *IndicesDeleteTemplateService {
	s.name = name
	return s
}

// Timeout is an explicit operation timeout.
func (s *IndicesDeleteTemplateService) Timeout(timeout string) *IndicesDeleteTemplateService {
	s.timeout = timeout
	return s
}

// MasterTimeout specifies the timeout for connection to master.
func (s *IndicesDeleteTemplateService) MasterTimeout(masterTimeout string) *IndicesDeleteTemplateService {
	s.masterTimeout = masterTimeout
	return s
}

// buildURL builds the URL for the operation.
func (s *IndicesDeleteTemplateService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_template/{name}", map[string]string{
		"name": s.name,
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
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *IndicesDeleteTemplateService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *IndicesDeleteTemplateService) Do(ctx context.Context) (*IndicesDeleteTemplateResponse, error) {
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
	ret := new(IndicesDeleteTemplateResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// IndicesDeleteTemplateResponse is the response of IndicesDeleteTemplateService.Do.
type IndicesDeleteTemplateResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index,omitempty"`
}
