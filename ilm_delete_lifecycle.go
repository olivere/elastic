// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// See the documentation at
// https://www.elastic.co/guide/en/elasticsearch/reference/6.7/ilm-get-lifecycle.html.
type IlmDeleteLifecycleService struct {
	client        *Client
	policy        string
	pretty        bool
	timeout       string
	masterTimeout string
	flatSettings  *bool
	local         *bool
}

// NewIlmDeleteLifecycleService creates a new IlmDeleteLifecycleService.
func NewIlmDeleteLifecycleService(client *Client) *IlmDeleteLifecycleService {
	return &IlmDeleteLifecycleService{
		client: client,
	}
}

// Policy is the name of the index lifecycle policy.
func (s *IlmDeleteLifecycleService) Policy(policy string) *IlmDeleteLifecycleService {
	s.policy = policy
	return s
}

// Timeout is an explicit operation timeout.
func (s *IlmDeleteLifecycleService) Timeout(timeout string) *IlmDeleteLifecycleService {
	s.timeout = timeout
	return s
}

// MasterTimeout specifies the timeout for connection to master.
func (s *IlmDeleteLifecycleService) MasterTimeout(masterTimeout string) *IlmDeleteLifecycleService {
	s.masterTimeout = masterTimeout
	return s
}

// FlatSettings is returns settings in flat format (default: false).
func (s *IlmDeleteLifecycleService) FlatSettings(flatSettings bool) *IlmDeleteLifecycleService {
	s.flatSettings = &flatSettings
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *IlmDeleteLifecycleService) Pretty(pretty bool) *IlmDeleteLifecycleService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *IlmDeleteLifecycleService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string
	if s.policy != "" {
		path, err = uritemplates.Expand("/_ilm/policy/{policy}", map[string]string{
			"policy": s.policy,
		})
	} else {
		path = "/_template"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.flatSettings != nil {
		params.Set("flat_settings", fmt.Sprintf("%v", *s.flatSettings))
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *IlmDeleteLifecycleService) Validate() error {
	var invalid []string
	if s.policy == "" {
		invalid = append(invalid, "Policy")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *IlmDeleteLifecycleService) Do(ctx context.Context) (*IlmDeleteLifecycleResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Delete URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Delete HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "DELETE",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(IlmDeleteLifecycleResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// IlmDeleteLifecycleResponse is the response of IlmDeleteLifecycleService.Do.
type IlmDeleteLifecycleResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
