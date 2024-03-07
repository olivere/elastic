package opensearch

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ClusterGetSettingService allows to get a very simple status on the health of the cluster.
//
// See http://www.opensearch.co/guide/en/opensearchsearch/reference/7.0/cluster-health.html
// for details.
type ClusterGetSettingService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	masterTimeout string
	timeout       string
}

// NewClusterGetSettingService creates a new ClusterGetSettingService.
func NewClusterGetSettingService(client *Client) *ClusterGetSettingService {
	return &ClusterGetSettingService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *ClusterGetSettingService) Pretty(pretty bool) *ClusterGetSettingService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *ClusterGetSettingService) Human(human bool) *ClusterGetSettingService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *ClusterGetSettingService) ErrorTrace(errorTrace bool) *ClusterGetSettingService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *ClusterGetSettingService) FilterPath(filterPath ...string) *ClusterGetSettingService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *ClusterGetSettingService) Header(name string, value string) *ClusterGetSettingService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *ClusterGetSettingService) Headers(headers http.Header) *ClusterGetSettingService {
	s.headers = headers
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *ClusterGetSettingService) MasterTimeout(masterTimeout string) *ClusterGetSettingService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *ClusterGetSettingService) Timeout(timeout string) *ClusterGetSettingService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *ClusterGetSettingService) buildURL() (string, url.Values, error) {
	// Build URL
	path := "/_cluster/settings"

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
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *ClusterGetSettingService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *ClusterGetSettingService) Do(ctx context.Context) (*ClusterGetSettingResponse, error) {
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
	ret := new(ClusterGetSettingResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// ClusterGetSettingResponse is the response of ClusterGetSettingService.Do.
type ClusterGetSettingResponse struct {
	ClusterName                    string  `json:"cluster_name"`
	Status                         string  `json:"status"`
	TimedOut                       bool    `json:"timed_out"`
	NumberOfNodes                  int     `json:"number_of_nodes"`
	NumberOfDataNodes              int     `json:"number_of_data_nodes"`
	ActivePrimaryShards            int     `json:"active_primary_shards"`
	ActiveShards                   int     `json:"active_shards"`
	RelocatingShards               int     `json:"relocating_shards"`
	InitializingShards             int     `json:"initializing_shards"`
	UnassignedShards               int     `json:"unassigned_shards"`
	DelayedUnassignedShards        int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks           int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch          int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitTimeInQueue         string  `json:"task_max_waiting_in_queue"`        // "0s"
	TaskMaxWaitTimeInQueueInMillis int     `json:"task_max_waiting_in_queue_millis"` // 0
	ActiveShardsPercent            string  `json:"active_shards_percent"`            // "100.0%"
	ActiveShardsPercentAsNumber    float64 `json:"active_shards_percent_as_number"`  // 100.0

	// Index name -> index health
	Indices map[string]*ClusterIndexHealth `json:"indices"`
}
