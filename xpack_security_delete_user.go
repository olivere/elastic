// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// XPackSecurityDeleteUserService delete a native user by its name.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.6/security-api-delete-user.html.
type XPackSecurityDeleteUserService struct {
	client *Client
	pretty bool
	name   string
}

// NewXPackSecurityDeleteUserService creates a new XPackSecurityDeleteUserService.
func NewXPackSecurityDeleteUserService(client *Client) *XPackSecurityDeleteUserService {
	return &XPackSecurityDeleteUserService{
		client: client,
	}
}

// Name is name of the native user to delete.
func (s *XPackSecurityDeleteUserService) Name(name string) *XPackSecurityDeleteUserService {
	s.name = name
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XPackSecurityDeleteUserService) Pretty(pretty bool) *XPackSecurityDeleteUserService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityDeleteUserService) buildURL() (string, url.Values, error) {
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
func (s *XPackSecurityDeleteUserService) Validate() error {
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
func (s *XPackSecurityDeleteUserService) Do(ctx context.Context) (*XPackSecurityDeleteUserResponse, error) {
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
		Method: "DELETE",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(XPackSecurityDeleteUserResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// XPackSecurityDeleteUserResponse is the response of XPackSecurityDeleteUserService.Do.
type XPackSecurityDeleteUserResponse struct {
	Found bool `json:"found"`
}
