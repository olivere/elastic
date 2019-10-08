// Copyright 2012-2019 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/v7/uritemplates"
)

// XPackSecurityEnableUserService retrieves a user by its name.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/security-api-get-user.html.
type XPackSecurityEnableUserService struct {
	client *Client
	pretty bool
	name   string
}

// NewXPackSecurityEnableUserService creates a new XPackSecurityEnableUserService.
func NewXPackSecurityEnableUserService(client *Client) *XPackSecurityEnableUserService {
	return &XPackSecurityEnableUserService{
		client: client,
	}
}

// Name is name of the user to retrieve.
func (s *XPackSecurityEnableUserService) Name(name string) *XPackSecurityEnableUserService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *XPackSecurityEnableUserService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_security/user/{name}/_enable", map[string]string{
		"name": s.name,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *XPackSecurityEnableUserService) Validate() error {
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
func (s *XPackSecurityEnableUserService) Do(ctx context.Context) error {
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
	return err
}
