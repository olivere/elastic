package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatNodesService allows to get nodes of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-master.html
// for details.
type CatNodesService struct {
	client        *Client
	format        string
	local         *bool
	masterTimeout string
}

// NewCatNodesService creates a new CatNodesService.
func NewCatNodesService(client *Client) *CatNodesService {
	return &CatNodesService{
		client: client,
		format: "json",
	}
}

// Format indicates that the JSON response be indented and human readable.
//func (s *CatNodesService) Format(format string) *CatNodesService {
//	if format != "" {
//		s.format = format
//	} else {
//		s.format = "json"
//	}
//	return s
//}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatNodesService) Local(local bool) *CatNodesService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatNodesService) MasterTimeout(masterTimeout string) *CatNodesService {
	s.masterTimeout = masterTimeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatNodesService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string
	path = "/_cat/nodes"

	// Add query string parameters
	params := url.Values{}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *CatNodesService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatNodesService) Do(ctx context.Context) (*CatNodesResponse, error) {
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
	ret := new(CatNodesResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Nodes); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatNodesResponse is the response of CatNodesService.Do
type CatNodesResponse struct {
	Nodes []*nodesRecord
}

type nodesRecord struct {
	Host        string `json:"host"`
	Ip          string `json:"ip"`
	HeapPercent string `json:"heap.percent"`
	RamPercent  string `json:"ram.percent"`
	Load        string `json:"load"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
	Name        string `json:"name"`
}
