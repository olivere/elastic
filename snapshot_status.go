package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// SnapshotStatusService lists the snapshots on a repository
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/modules-snapshots.html
// for details.
type SnapshotStatusService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	snapshot          []string
	repository        string
	masterTimeout     string
	ignoreUnavailable *bool
}

// NewSnapshotStatusService creates a new SnapshotStatusService.
func NewSnapshotStatusService(client *Client) *SnapshotStatusService {
	return &SnapshotStatusService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *SnapshotStatusService) Pretty(pretty bool) *SnapshotStatusService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SnapshotStatusService) Human(human bool) *SnapshotStatusService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SnapshotStatusService) ErrorTrace(errorTrace bool) *SnapshotStatusService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SnapshotStatusService) FilterPath(filterPath ...string) *SnapshotStatusService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SnapshotStatusService) Header(name string, value string) *SnapshotStatusService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SnapshotStatusService) Headers(headers http.Header) *SnapshotStatusService {
	s.headers = headers
	return s
}

// Repository is the repository name.
func (s *SnapshotStatusService) Repository(repository string) *SnapshotStatusService {
	s.repository = repository
	return s
}

// Snapshot is the list of snapshot names. If not set, defaults to all snapshots.
func (s *SnapshotStatusService) Snapshot(snapshots ...string) *SnapshotStatusService {
	s.snapshot = append(s.snapshot, snapshots...)
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *SnapshotStatusService) MasterTimeout(masterTimeout string) *SnapshotStatusService {
	s.masterTimeout = masterTimeout
	return s
}

// buildURL builds the URL for the operation.
func (s *SnapshotStatusService) buildURL() (string, url.Values, error) {
	var err error
	var path string

	if s.repository != "" {
		if len(s.snapshot) > 0 {
			path, err = uritemplates.Expand("/_snapshot/{repository}/{snapshot}/_status", map[string]string{
				"repository": s.repository,
				"snapshot":   strings.Join(s.snapshot, ","),
			})
		} else {
			path, err = uritemplates.Expand("/_snapshot/{repository}/_status", map[string]string{
				"repository": s.repository,
			})
		}
	} else {
		path, err = uritemplates.Expand("/_snapshot/_status", nil)
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if v := s.ignoreUnavailable; v != nil {
		params.Set("ignore_unavailable", fmt.Sprint(*v))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
// Validation only fails if snapshot names were provided but no repository was
// provided.
// All other cases are Okay
func (s *SnapshotStatusService) Validate() error {
	if len(s.snapshot) > 0 && s.repository == "" {
		return fmt.Errorf("snapshots were specified but repository is missing")
	}
	return nil
}

// Do executes the operation.
func (s *SnapshotStatusService) Do(ctx context.Context) (*SnapshotStatusResponse, error) {
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
		Method:  "GET",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SnapshotStatusResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type SnapshotStatusResponse struct {
	Snapshots []SnapshotStatus `json:"snapshots"`
}

type SnapshotStatus struct {
	Snapshot           string                 `json:"snapshot"`
	Repository         string                 `json:"repository"`
	UUID               string                 `json:"uuid"`
	State              string                 `json:"state"`
	IncludeGlobalState bool                   `json:"include_global_state"`
	ShardsStats        ShardsStats            `json:"shards_stats"`
	Stats              Stats                  `json:"stats"`
	Indices            map[string]IndiceStats `json:"indices"`
}

type ShardsStats struct {
	Initializing int `json:"initializing"`
	Started      int `json:"started"`
	Finalizing   int `json:"finalizing"`
	Done         int `json:"done"`
	Failed       int `json:"failed"`
	Total        int `json:"total"`
}

type incremental struct {
	FileCount   int `json:"file_count"`
	SizeInBytes int `json:"size_in_bytes"`
}

type total struct {
	FileCount   int `json:"file_count"`
	SizeInBytes int `json:"size_in_bytes"`
}

type Stats struct {
	Incremental          incremental `json:"incremental"`
	Total                total       `json:"total"`
	StartTimeInMillis    int         `json:"start_time_in_millis"`
	TimeInMillis         int64       `json:"time_in_millis"`
	NumberOfFiles        int         `json:"number_of_files"`
	ProcessedFiles       int         `json:"processed_files"`
	TotalSizeInBytes     int         `json:"total_size_in_bytes"`
	ProcessedSizeInBytes int         `json:"processed_size_in_bytes"`
}

type Shard struct {
	Stage string `json:"stage"`
	Stats Stats  `json:"stats"`
	Node  string `json:"node"`
}

type IndiceStats struct {
	ShardsStats ShardsStats      `json:"shards_stats"`
	Stats       Stats            `json:"stats"`
	Shards      map[string]Shard `json:"shards"`
}
