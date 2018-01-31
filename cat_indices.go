package elastic

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// CatIndicesService allows to get the master node of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-master.html
// for details.
type CatIndicesService struct {
	client        *Client
	format        string
	pretty        bool
	indices       []string
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatIndicesService creates a new CatIndicesService.
func NewCatIndicesService(client *Client) *CatIndicesService {
	return &CatIndicesService{
		client:  client,
		indices: make([]string, 0),
		format:  "json",
	}
}

// Index limits the information returned to specific indices.
func (s *CatIndicesService) Index(indices ...string) *CatIndicesService {
	s.indices = append(s.indices, indices...)
	return s
}

// Format indicates that the JSON response be indented and human readable.
//func (s *CatIndicesService) Format(format string) *CatIndicesService {
//	if format != "" {
//		s.format = format
//	} else {
//		s.format = "json"
//	}
//	return s
//}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatIndicesService) Pretty(pretty bool) *CatIndicesService {
	s.pretty = pretty
	return s
}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatIndicesService) Local(local bool) *CatIndicesService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatIndicesService) MasterTimeout(masterTimeout string) *CatIndicesService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatIndicesService) Timeout(timeout string) *CatIndicesService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatIndicesService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string

	if len(s.indices) > 0 {
		path, err = uritemplates.Expand("/_cat/indices/{index}", map[string]string{
			"index": strings.Join(s.indices, ","),
		})
	} else {
		path = "/_cat/indices"
	}

	if err != nil {
		return "", url.Values{}, err
	}

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
func (s *CatIndicesService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatIndicesService) Do(ctx context.Context) (*CatIndicesResponse, error) {
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
	ret := new(CatIndicesResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Indices); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatIndicesResponse is the response of CatIndicesService.Do
type CatIndicesResponse struct {
	Indices []*indicesRecord
}

type indicesRecord struct {
	Health    string `json:"health"`
	Status    string `json:"status"`
	Index     string `json:"index"`
	Pri       string `json:"pri"`
	Rep       string `json:"rep"`
	Count     string `json:"docs.count"`
	Deleted   string `json:"docs.deleted"`
	Size      string `json:"store.size"`
	StoreSize string `json:"pri.store.size"`
}
