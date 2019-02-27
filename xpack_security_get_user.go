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

// XPackSecurityGetUserService retrieves a native user by its name.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.6/security-api-get-user.html.
type XPackSecurityGetUserService struct {
	client *Client
	pretty bool
	name   string
}

// NewXPackSecurityGetUserService creates a new XPackSecurityGetUserService.
func NewXPackSecurityGetUserService(client *Client) *XPackSecurityGetUserService {
	return &XPackSecurityGetUserService{
		client: client,
	}
}

// Name is name of the user to retrieve.
func (s *XPackSecurityGetUserService) Name(name string) *XPackSecurityGetUserService {
	s.name = name
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XPackSecurityGetUserService) Pretty(pretty bool) *XPackSecurityGetUserService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityGetUserService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/security/user/{name}", map[string]string{
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
func (s *XPackSecurityGetUserService) Validate() error {
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
func (s *XPackSecurityGetUserService) Do(ctx context.Context) (*XPackSecurityGetUserResponse, error) {
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
	ret := XPackSecurityGetUserResponse{}
	if err := json.Unmarshal(res.Body, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// XPackSecurityGetUserResponse is the response of XPackSecurityGetUserService.Do.
type XPackSecurityGetUserResponse map[string]XPackSecurityUser

// XPackSecurityUser is the native user object.
//
// The Java source for this struct is defined here:
// https://github.com/elastic/elasticsearch/blob/master/x-pack/plugin/core/src/main/java/org/elasticsearch/xpack/core/security/user/User.java
type XPackSecurityUser struct {
	Username string                 `json:"username"`
	Roles    []string               `json:"roles"`
	FullName string                 `json:"full_name"`
	Email    string                 `json:"email"`
	Metadata map[string]interface{} `json:"metadata"`
	Enabled  bool                   `json:"enabled"`
}
