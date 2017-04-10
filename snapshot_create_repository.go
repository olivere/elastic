// Copyright 2012-2017 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"gopkg.in/olivere/elastic.v5/uritemplates"
)

// SnapshotCreateRepositoryService is documented at https://www.elastic.co/guide/en/elasticsearch/reference/5.x/modules-snapshots.html.
type SnapshotCreateRepositoryService struct {
	client        *Client
	pretty        bool
	repository    string
	masterTimeout string
	timeout       string
	verify        *bool
	bodyJson      interface{}
	bodyString    string
}

// NewSnapshotCreateRepositoryService creates a new SnapshotCreateRepositoryService.
func NewSnapshotCreateRepositoryService(client *Client) *SnapshotCreateRepositoryService {
	return &SnapshotCreateRepositoryService{
		client: client,
	}
}

// Repository is documented as: A repository name.
func (s *SnapshotCreateRepositoryService) Repository(repository string) *SnapshotCreateRepositoryService {
	s.repository = repository
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *SnapshotCreateRepositoryService) MasterTimeout(masterTimeout string) *SnapshotCreateRepositoryService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout is documented as: Explicit operation timeout.
func (s *SnapshotCreateRepositoryService) Timeout(timeout string) *SnapshotCreateRepositoryService {
	s.timeout = timeout
	return s
}

// Verify is documented as: Whether to verify the repository after creation.
func (s *SnapshotCreateRepositoryService) Verify(verify bool) *SnapshotCreateRepositoryService {
	s.verify = &verify
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *SnapshotCreateRepositoryService) Pretty(pretty bool) *SnapshotCreateRepositoryService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: The repository definition.
func (s *SnapshotCreateRepositoryService) BodyJson(body interface{}) *SnapshotCreateRepositoryService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: The repository definition.
func (s *SnapshotCreateRepositoryService) BodyString(body string) *SnapshotCreateRepositoryService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SnapshotCreateRepositoryService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_snapshot/{repository}", map[string]string{
		"repository": s.repository,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.verify != nil {
		params.Set("verify", fmt.Sprintf("%v", *s.verify))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *SnapshotCreateRepositoryService) Validate() error {
	var invalid []string
	if s.repository == "" {
		invalid = append(invalid, "Repository")
	}
	if s.bodyString == "" && s.bodyJson == nil {
		invalid = append(invalid, "BodyJson")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SnapshotCreateRepositoryService) Do(ctx context.Context) (*SnapshotCreateRepositoryResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request body
	var body interface{}
	if s.bodyJson != nil {
		body = s.bodyJson
	} else {
		body = s.bodyString
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, "PUT", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SnapshotCreateRepositoryResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SnapshotCreateRepositoryResponse is the response of SnapshotCreateRepositoryService.Do.
type SnapshotCreateRepositoryResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
