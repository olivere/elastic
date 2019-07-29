// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/v7/uritemplates"
)

// XPackSecurityChangePasswordService changes a native user's password.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.1/security-api-change-password.html.
type XPackSecurityChangePasswordService struct {
	client   *Client
	pretty   bool
	username string
	password string
	refresh  string
	body     interface{}
}

// NewXPackSecurityChangePasswordService creates a new XPackSecurityChangePasswordService.
func NewXPackSecurityChangePasswordService(client *Client) *XPackSecurityChangePasswordService {
	return &XPackSecurityChangePasswordService{
		client: client,
	}
}

// Username is name of the user to change.
func (s *XPackSecurityChangePasswordService) Username(username string) *XPackSecurityChangePasswordService {
	s.username = username
	return s
}

// Password is the new value of the password.
func (s *XPackSecurityChangePasswordService) Password(password string) *XPackSecurityChangePasswordService {
	s.password = password
	return s
}

// Refresh, if "true" (the default), refreshes the affected shards to make this operation
// visible to search, if "wait_for" then wait for a refresh to make this operation visible
// to search, if "false" then do nothing with refreshes.
func (s *XPackSecurityChangePasswordService) Refresh(refresh string) *XPackSecurityChangePasswordService {
	s.refresh = refresh
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XPackSecurityChangePasswordService) Pretty(pretty bool) *XPackSecurityChangePasswordService {
	s.pretty = pretty
	return s
}

// Body specifies the password. Use a string or a type that will get serialized as JSON.
func (s *XPackSecurityChangePasswordService) Body(body interface{}) *XPackSecurityChangePasswordService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityChangePasswordService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/security/user/{username}/_password", map[string]string{
		"username": s.username,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if v := s.refresh; v != "" {
		params.Set("refresh", v)
	}
	if s.pretty {
		params.Set("pretty", "true")
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XPackSecurityChangePasswordService) Validate() error {
	var invalid []string
	if s.username == "" {
		invalid = append(invalid, "Userame")
	}
	if s.password == "" && s.body == nil {
		invalid = append(invalid, "Body")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *XPackSecurityChangePasswordService) Do(ctx context.Context) (*XPackSecurityChangeUserPasswordResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	var body interface{}
	if s.body != nil {
		body = s.body
	} else {
		body = map[string]interface{}{
			"password": s.password,
		}
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "POST",
		Path:   path,
		Params: params,
		Body:   body,
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

// XPackSecurityChangeUserPasswordResponse is the response of
// XPackSecurityChangePasswordService.Do.
//
// A successful call returns an empty JSON structure: {}.
type XPackSecurityChangeUserPasswordResponse struct {
}
