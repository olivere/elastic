package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v5/uritemplates"
)

// SnapshotCreateService creates a snapshot.
// See https://www.elastic.co/guide/en/elasticsearch/reference/5.3/modules-snapshots.html
// for details.
type SnapshotCreateService struct {
	includeGlobalState *bool
	ignoreUnavailable  bool
	partial            bool
	waitForCompletion  bool
	pretty             bool
	client             *Client
	masterTimeout      string
	timeout            string
	indices            string
	name               string
	repository         string
}

// WaitForCompletion indicates whether to wait until the operation has completed before returning.
func (s *SnapshotCreateService) WaitForCompletion(waitForCompletion bool) *SnapshotCreateService {
	s.waitForCompletion = waitForCompletion
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *SnapshotCreateService) MasterTimeout(masterTimeout string) *SnapshotCreateService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout is an explicit operation timeout.
func (s *SnapshotCreateService) Timeout(timeout string) *SnapshotCreateService {
	s.timeout = timeout
	return s
}

// IgnoreUnavailable is documented as: Setting it to true will cause indices that do not exist to be
// ignored during snapshot creation. By default, when ignore_unavailable option is not set and an
// index is missing the snapshot request will fail.
func (s *SnapshotCreateService) IgnoreUnavailable(ignoreUnavailable bool) *SnapshotCreateService {
	s.ignoreUnavailable = ignoreUnavailable
	return s
}

// IncludeGlobalState is documented as: By setting include_global_state to false it’s possible to
// prevent the cluster global state to be stored as part of the snapshot.
func (s *SnapshotCreateService) IncludeGlobalState(includeGlobalState bool) *SnapshotCreateService {
	s.includeGlobalState = &includeGlobalState
	return s
}

// Partial is documented as: By default, the entire snapshot will fail if one or more
// indices participating in the snapshot don’t have all primary shards available. This behaviour can
// be changed by setting partial to true.
func (s *SnapshotCreateService) Partial(partial bool) *SnapshotCreateService {
	s.partial = partial
	return s
}

// Indices is documented as: The list of indices to create snapshot
func (s *SnapshotCreateService) Indices(indices string) *SnapshotCreateService {
	s.indices = indices
	return s
}

// Repository is the repository name.
func (s *SnapshotCreateService) Repository(repository string) *SnapshotCreateService {
	s.repository = repository
	return s
}

// Name is documented as: The snapshot name.
func (s *SnapshotCreateService) Name(name string) *SnapshotCreateService {
	s.name = name
	return s
}

// NewSnapshotCreateService creates a new SnapshotCreateService.
func NewSnapshotCreateService(client *Client) *SnapshotCreateService {
	return &SnapshotCreateService{
		client: client,
	}
}

// Validate checks if the operation is valid.
func (s *SnapshotCreateService) Validate() error {
	var invalid []string
	if s.repository == "" {
		invalid = append(invalid, "Repository")
	}
	if s.name == "" {
		invalid = append(invalid, "Name")
	}
	if s.indices == "" {
		invalid = append(invalid, "Indices")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *SnapshotCreateService) Pretty(pretty bool) *SnapshotCreateService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *SnapshotCreateService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_snapshot/{repository}/{name}", map[string]string{
		"name":       s.name,
		"repository": s.repository,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	params.Set("wait_for_completion", strconv.FormatBool(s.waitForCompletion))
	if s.waitForCompletion {
		params.Set("wait_for_completion", "true")
	}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	return path, params, nil
}

// buildBody builds the body for the operation.
func (s *SnapshotCreateService) buildBody() (interface{}, error) {
	body := map[string]interface{}{
		"indices":            s.indices,
		"ignore_unavailable": s.ignoreUnavailable,
		"partial":            s.partial,
	}
	if s.includeGlobalState != nil {
		body["include_global_state"] = *s.includeGlobalState
	}
	return body, nil
}

// Do creates a snapshot in repository.
func (s *SnapshotCreateService) Do(
	ctx context.Context,
) (*SnapshotCreateResponse, error) {
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
	body, err := s.buildBody()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, http.MethodPut, path, params, body)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SnapshotCreateResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SnapshotCreateResponse is the response of SnapshotCreateResponse.Do.
type SnapshotCreateResponse struct {
	// Accepted indicates whether the request was accepted by elasticsearch.
	// It's available when waitForCompletion is false.
	Accepted *bool `json:"accepted"`
	// Snapshot is available when waitForCompletion is true.
	Snapshot *struct {
		Snapshot          string        `json:"snapshot"`
		VersionID         int           `json:"version_id"`
		Version           string        `json:"version"`
		Indices           []string      `json:"indices"`
		State             string        `json:"state"`
		StartTime         time.Time     `json:"start_time"`
		StartTimeInMillis int64         `json:"start_time_in_millis"`
		EndTime           time.Time     `json:"end_time"`
		EndTimeInMillis   int64         `json:"end_time_in_millis"`
		DurationInMillis  int           `json:"duration_in_millis"`
		Failures          []interface{} `json:"failures"`
		Shards            struct {
			Total      int `json:"total"`
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
		} `json:"shards"`
	} `json:"snapshot"`
}
