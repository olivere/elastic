package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatHealthService allows to get the health of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-health.html
// for details.
type CatHealthService struct {
	client        *Client
	format        string
	pretty        bool
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatHealthService creates a new CatHealthService.
func NewCatHealthService(client *Client) *CatHealthService {
	return &CatHealthService{
		client: client,
	}
}

// Format health that the JSON response be indented and human readable.
func (s *CatHealthService) Format(format string) *CatHealthService {
	if format != "" {
		s.format = format
	} else {
		s.format = "json"
	}
	return s
}

// Pretty health that the JSON response be indented and human readable.
func (s *CatHealthService) Pretty(pretty bool) *CatHealthService {
	s.pretty = pretty
	return s
}

// Local health whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatHealthService) Local(local bool) *CatHealthService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatHealthService) MasterTimeout(masterTimeout string) *CatHealthService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatHealthService) Timeout(timeout string) *CatHealthService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatHealthService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string

	path = "/_cat/health"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
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
func (s *CatHealthService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatHealthService) Do(ctx context.Context) (*CatHealthResponse, error) {
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
		Method: "GET",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(CatHealthResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Healths); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatHealthResponse is the response of CatHealthService.Do
type CatHealthResponse struct {
	Healths []*healthRecord
}

type healthRecord struct {
	ActiveShardsPercent string `json:"active_shards_percent"`
	Cluster             string `json:"cluster"`
	Epoch               string `json:"epoch"`
	Init                string `json:"init"`
	MaxTaskWaitTime     string `json:"max_task_wait_time"`
	NodeData            string `json:"node.data"`
	NodeTotal           string `json:"node.total"`
	PendingTasks        string `json:"pending_tasks"`
	Pri                 string `json:"pri"`
	Relo                string `json:"relo"`
	Shards              string `json:"shards"`
	Status              string `json:"status"`
	Timestamp           string `json:"timestamp"`
	Unassign            string `json:"unassign"`
}
