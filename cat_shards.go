package elastic

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// CatShardsService allows to get the master node of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-master.html
// for details.
type CatShardsService struct {
	client        *Client
	format        string
	pretty        bool
	indices       []string
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatShardsService creates a new CatShardsService.
func NewCatShardsService(client *Client) *CatShardsService {
	return &CatShardsService{
		client:  client,
		indices: make([]string, 0),
		format:  "json",
	}
}

// Index limits the information returned to specific indices.
func (s *CatShardsService) Index(indices ...string) *CatShardsService {
	s.indices = append(s.indices, indices...)
	return s
}

// Format indicates that the JSON response be indented and human readable.
//func (s *CatShardsService) Format(format string) *CatShardsService {
//	if format != "" {
//		s.format = format
//	} else {
//		s.format = "json"
//	}
//	return s
//}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatShardsService) Pretty(pretty bool) *CatShardsService {
	s.pretty = pretty
	return s
}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatShardsService) Local(local bool) *CatShardsService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatShardsService) MasterTimeout(masterTimeout string) *CatShardsService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatShardsService) Timeout(timeout string) *CatShardsService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatShardsService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string

	if len(s.indices) > 0 {
		path, err = uritemplates.Expand("/_cat/shards/{index}", map[string]string{
			"index": strings.Join(s.indices, ","),
		})
	} else {
		path = "/_cat/shards"
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
func (s *CatShardsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatShardsService) Do(ctx context.Context) (*CatShardsResponse, error) {
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
	ret := new(CatShardsResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Shards); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatShardsResponse is the response of CatShardsService.Do
type CatShardsResponse struct {
	Shards []*shardsRecord
}

type shardsRecord struct {
	Index  string `json:"index"`
	Shard  string `json:"shard"`
	Prirep string `json:"prirep"`
	State  string `json:"state"`
	Docs   string `json:"docs"`
	Store  string `json:"store"`
	Ip     string `json:"ip"`
	Node   string `json:"node"`
}
