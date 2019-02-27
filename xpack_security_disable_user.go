// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"github.com/thales-e-security/elastic/uritemplates"
	"net/url"
)

// XPackSecurityDisableUserService disables a native user by its name.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.6/security-api-disable-user.html.
type XPackSecurityDisableUserService struct {
	client *Client
	pretty bool
	name   string
}

// NewXPackSecurityDisableUserService creates a new XPackSecurityDisableUserService.
func NewXPackSecurityDisableUserService(client *Client) *XPackSecurityDisableUserService {
	return &XPackSecurityDisableUserService{
		client: client,
	}
}

// Name is name of the user to create.
func (s *XPackSecurityDisableUserService) Name(name string) *XPackSecurityDisableUserService {
	s.name = name
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *XPackSecurityDisableUserService) Pretty(pretty bool) *XPackSecurityDisableUserService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityDisableUserService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_xpack/security/user/{name}/_disable", map[string]string{
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
func (s *XPackSecurityDisableUserService) Validate() error {
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
func (s *XPackSecurityDisableUserService) Do(ctx context.Context) error {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return err
	}

	// Get HTTP response
	_, err = s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "PUT",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return err
	}

	return nil
}
