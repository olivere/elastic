// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// ClusterHealthService allows to get the status of the cluster.
// It is documented at http://www.elasticsearch.org/guide/en/elasticsearch/reference/1.4/cluster-health.html.
type ClusterHealthService struct {
	client                  *Client
	debug                   bool
	pretty                  bool
	indices                 []string
	waitForStatus           string
	level                   string
	local                   *bool
	masterTimeout           string
	timeout                 string
	waitForActiveShards     *int
	waitForNodes            string
	waitForRelocatingShards *int
}

// NewClusterHealthService creates a new ClusterHealthService.
func NewClusterHealthService(client *Client) *ClusterHealthService {
	return &ClusterHealthService{client: client, indices: make([]string, 0)}
}

// Index is documented as: Limit the information returned to a specific index.
func (s *ClusterHealthService) Index(index string) *ClusterHealthService {
	s.indices = make([]string, 0)
	s.indices = append(s.indices, index)
	return s
}

// Indices is documented as: Limit the information returned to a specific index.
func (s *ClusterHealthService) Indices(indices ...string) *ClusterHealthService {
	s.indices = make([]string, 0)
	s.indices = append(s.indices, indices...)
	return s
}

// MasterTimeout is documented as: Explicit operation timeout for connection to master node.
func (s *ClusterHealthService) MasterTimeout(masterTimeout string) *ClusterHealthService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout is documented as: Explicit operation timeout.
func (s *ClusterHealthService) Timeout(timeout string) *ClusterHealthService {
	s.timeout = timeout
	return s
}

// WaitForActiveShards is documented as: Wait until the specified number of shards is active.
func (s *ClusterHealthService) WaitForActiveShards(waitForActiveShards int) *ClusterHealthService {
	s.waitForActiveShards = &waitForActiveShards
	return s
}

// WaitForNodes is documented as: Wait until the specified number of nodes is available.
func (s *ClusterHealthService) WaitForNodes(waitForNodes string) *ClusterHealthService {
	s.waitForNodes = waitForNodes
	return s
}

// WaitForRelocatingShards is documented as: Wait until the specified number of relocating shards is finished.
func (s *ClusterHealthService) WaitForRelocatingShards(waitForRelocatingShards int) *ClusterHealthService {
	s.waitForRelocatingShards = &waitForRelocatingShards
	return s
}

// WaitForStatus is documented as: Wait until cluster is in a specific state.
func (s *ClusterHealthService) WaitForStatus(waitForStatus string) *ClusterHealthService {
	s.waitForStatus = waitForStatus
	return s
}

// Level is documented as: Specify the level of detail for returned information.
func (s *ClusterHealthService) Level(level string) *ClusterHealthService {
	s.level = level
	return s
}

// Local is documented as: Return local information, do not retrieve the state from master node (default: false).
func (s *ClusterHealthService) Local(local bool) *ClusterHealthService {
	s.local = &local
	return s
}

// buildURL builds the URL for the operation.
func (s *ClusterHealthService) buildURL() (string, error) {
	// Build URL
	urls, err := uritemplates.Expand("/_cluster/health/{index}", map[string]string{
		"index": strings.Join(s.indices, ","),
	})
	if err != nil {
		return "", err
	}

	// Add query string parameters
	params := url.Values{}
	if s.waitForRelocatingShards != nil {
		params.Set("waitForRelocatingShards", fmt.Sprintf("%d", *s.waitForRelocatingShards))
	}
	if s.waitForStatus != "" {
		params.Set("waitForStatus", s.waitForStatus)
	}
	if s.level != "" {
		params.Set("level", s.level)
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	if s.masterTimeout != "" {
		params.Set("masterTimeout", s.masterTimeout)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.waitForActiveShards != nil {
		params.Set("waitForActiveShards", fmt.Sprintf("%d", *s.waitForActiveShards))
	}
	if s.waitForNodes != "" {
		params.Set("waitForNodes", s.waitForNodes)
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	return urls, nil
}

// Validate checks if the operation is valid.
func (s *ClusterHealthService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *ClusterHealthService) Do() (*ClusterHealthResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	urls, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request
	req, err := s.client.NewRequest("GET", urls)
	if err != nil {
		return nil, err
	}

	// Debug output?
	if s.debug {
		out, _ := httputil.DumpRequestOut((*http.Request)(req), true)
		log.Printf("%s\n", string(out))
	}

	// Get HTTP response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Debug output?
	if s.debug {
		out, _ := httputil.DumpResponse(res, true)
		log.Printf("%s\n", string(out))
	}
	// Return operation response
	resp := new(ClusterHealthResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ClusterHealthResponse is the response of ClusterHealthService.Do.
type ClusterHealthResponse struct {
	ClusterName         string `json:"cluster_name"`
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializedShards   int    `json:"initialized_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}
