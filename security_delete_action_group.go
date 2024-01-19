// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/disaster37/opensearch/v2/uritemplates"
)

// SecurityDeleteActionGroupService delete a action group by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#delete-action-group
type SecurityDeleteActionGroupService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityDeleteActionGroupService creates a new SecurityDeleteActionGroupService.
func NewSecurityDeleteActionGroupService(client *Client) *SecurityDeleteActionGroupService {
	return &SecurityDeleteActionGroupService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityDeleteActionGroupService) Pretty(pretty bool) *SecurityDeleteActionGroupService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityDeleteActionGroupService) Human(human bool) *SecurityDeleteActionGroupService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityDeleteActionGroupService) ErrorTrace(errorTrace bool) *SecurityDeleteActionGroupService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityDeleteActionGroupService) FilterPath(filterPath ...string) *SecurityDeleteActionGroupService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityDeleteActionGroupService) Header(name string, value string) *SecurityDeleteActionGroupService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityDeleteActionGroupService) Headers(headers http.Header) *SecurityDeleteActionGroupService {
	s.headers = headers
	return s
}

// Name is name of the action group to delete.
func (s *SecurityDeleteActionGroupService) Name(name string) *SecurityDeleteActionGroupService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityDeleteActionGroupService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/actiongroups/{name}", map[string]string{
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
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *SecurityDeleteActionGroupService) Validate() error {
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
func (s *SecurityDeleteActionGroupService) Do(ctx context.Context) (*SecurityDeleteActionGroupResponse, error) {
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
	ret := new(SecurityDeleteActionGroupResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityDeleteActionGroupResponse the response of SecurityDeleteActionGroupService.Do.
type SecurityDeleteActionGroupResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
