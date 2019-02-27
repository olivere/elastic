// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thales-e-security/elastic/uritemplates"
	"net/url"
)

// XPackSecurityChangeUserPasswordService changes a native user's password.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.6/security-api-change-password.html.
type XPackSecurityChangeUserPasswordService struct {
	client *Client
	pretty bool
	name   string
	body   interface{}
}

// NewXPackSecurityChangeUserPasswordService creates a new XPackSecurityChangeUserPasswordService.
func NewXPackSecurityChangeUserPasswordService(client *Client) *XPackSecurityChangeUserPasswordService {
	return &XPackSecurityChangeUserPasswordService{
		client: client,
	}
}

// Name is name of the user to change.
func (s *XPackSecurityChangeUserPasswordService) Name(name string) *XPackSecurityChangeUserPasswordService {
	s.name = name
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XPackSecurityChangeUserPasswordService) Pretty(pretty bool) *XPackSecurityChangeUserPasswordService {
	s.pretty = pretty
	return s
}

// Body specifies the password. Use a string or a type that will get serialized as JSON.
func (s *XPackSecurityChangeUserPasswordService) Body(body interface{}) *XPackSecurityChangeUserPasswordService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityChangeUserPasswordService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/security/user/{name}/_password", map[string]string{
		"name": s.name,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XPackSecurityChangeUserPasswordService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
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
func (s *XPackSecurityChangeUserPasswordService) Do(ctx context.Context) (*XPackSecurityChangeUserPasswordResponse, error) {
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
		Method: "POST",
		Path:   path,
		Params: params,
		Body:   s.body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XPackSecurityChangeUserPasswordResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XPackSecurityChangeUserPasswordResponse is the response of XPackSecurityChangeUserPasswordService.Do.
// A successful call returns an empty JSON structure '{}'
type XPackSecurityChangeUserPasswordResponse struct {
}
