package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatMasterService allows to get the master node of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-master.html
// for details.
type CatMasterService struct {
	client        *Client
	format        string
	pretty        bool
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatMasterService creates a new CatMasterService.
func NewCatMasterService(client *Client) *CatMasterService {
	return &CatMasterService{
		client: client,
	}
}

// Format indicates that the JSON response be indented and human readable.
func (s *CatMasterService) Format(format string) *CatMasterService {
	if format != "" {
		s.format = format
	} else {
		s.format = "json"
	}
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatMasterService) Pretty(pretty bool) *CatMasterService {
	s.pretty = pretty
	return s
}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatMasterService) Local(local bool) *CatMasterService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatMasterService) MasterTimeout(masterTimeout string) *CatMasterService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatMasterService) Timeout(timeout string) *CatMasterService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatMasterService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string

	path = "/_cat/master"

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
func (s *CatMasterService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatMasterService) Do(ctx context.Context) (*CatMasterResponse, error) {
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
	ret := new(CatMasterResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Masters); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatMasterResponse is the response of CatMasterService.Do
type CatMasterResponse struct {
	Masters []*masterRecord
}

type masterRecord struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Ip   string `json:"ip"`
	Node string `json:"node"`
}
