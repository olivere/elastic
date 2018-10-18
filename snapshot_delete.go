// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
)

// SnapshotDeleteService is documented at https://www.elastic.co/guide/en/elasticsearch/reference/6.2/modules-snapshots.html.
type SnapshotDeleteService struct {
	client            *Client
	pretty            bool
	repository        string
	snapshot          string
	masterTimeout     string
	waitForCompletion *bool
	bodyJson          interface{}
	bodyString        string
}

func (c *Client) SnapshotDelete(repository string, snapshot string) *SnapshotDeleteService {
	return NewSnapshotDeleteService(c).Repository(repository).Snapshot(snapshot)
}

// NewSnapshotDeleteService creates a new SnapshotDeleteService.
func NewSnapshotDeleteService(client *Client) *SnapshotDeleteService {
	return &SnapshotDeleteService{
		client: client,
	}
}

// Repository is the repository name.
func (s *SnapshotDeleteService) Repository(repository string) *SnapshotDeleteService {
	s.repository = repository
	return s
}

// Snapshot is the snapshot name.
func (s *SnapshotDeleteService) Snapshot(snapshot string) *SnapshotDeleteService {
	s.snapshot = snapshot
	return s
}

// BodyJson is documented as: The snapshot definition.
func (s *SnapshotDeleteService) BodyJson(body interface{}) *SnapshotDeleteService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: The snapshot definition.
func (s *SnapshotDeleteService) BodyString(body string) *SnapshotDeleteService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SnapshotDeleteService) buildURL() string {
	// Build URL
	path := fmt.Sprintf("/_snapshot/%s/%s", s.repository, s.snapshot)

	return path
}

// Validate checks if the operation is valid.
func (s *SnapshotDeleteService) Validate() error {
	var invalid []string
	if s.repository == "" {
		invalid = append(invalid, "Repository")
	}
	if s.snapshot == "" {
		invalid = append(invalid, "Snapshot")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SnapshotDeleteService) Do(ctx context.Context) (*SnapshotDeleteResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path := s.buildURL()

	// Setup HTTP request body
	var body interface{}
	if s.bodyJson != nil {
		body = s.bodyJson
	} else {
		body = s.bodyString
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "DELETE",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SnapshotDeleteResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SnapshotDeleteResponse is the response of SnapshotDeleteService.Do.
type SnapshotDeleteResponse struct {
	// Accepted indicates whether the delete operation was successful.
	// It's available when waitForCompletion is false.
	Accepted *bool `json:"acknowledged"`
}
