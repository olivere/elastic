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

// SmDeletePolicyService delete a policy by its name.
// See https://opensearch.org/docs/latest/im-plugin/ism/api/#delete-policy
type SmDeletePolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSmDeletePolicyService creates a new SmDeletePolicyService.
func NewSmDeletePolicyService(client *Client) *SmDeletePolicyService {
	return &SmDeletePolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SmDeletePolicyService) Pretty(pretty bool) *SmDeletePolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SmDeletePolicyService) Human(human bool) *SmDeletePolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SmDeletePolicyService) ErrorTrace(errorTrace bool) *SmDeletePolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SmDeletePolicyService) FilterPath(filterPath ...string) *SmDeletePolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SmDeletePolicyService) Header(name string, value string) *SmDeletePolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SmDeletePolicyService) Headers(headers http.Header) *SmDeletePolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to delete.
func (s *SmDeletePolicyService) Name(name string) *SmDeletePolicyService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SmDeletePolicyService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_sm/policies/{name}", map[string]string{
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
func (s *SmDeletePolicyService) Validate() error {
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
func (s *SmDeletePolicyService) Do(ctx context.Context) (*SmDeletePolicyResponse, error) {
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
	ret := new(SmDeletePolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SmDeletePolicyResponse is the response of SmDeletePolicyService.Do.
type SmDeletePolicyResponse struct {
	Index          string         `json:"_index"`
	ID             string         `json:"_id"`
	Version        int64          `json:"_version"`
	Result         string         `json:"result"`
	ForcedRefresh  bool           `json:"forced_refresh"`
	Shard          map[string]any `json:"_shards"`
	SequenceNumber int64          `json:"_seq_no"`
	PrimaryTerm    int64          `json:"_primary_term"`
}
