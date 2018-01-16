// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/olivere/elastic.v5/uritemplates"
)

// IndicesSegmentsService is documented at https://www.elastic.co/guide/en/elasticsearch/reference/5.x/indices-segments.html.
type IndicesSegmentsService struct {
	client             *Client
	pretty             bool
	index              []string
	allowNoIndices     *bool
	expandWildcards    string
	ignoreUnavailable  *bool
	operationThreading interface{}
	verbose            *bool
}

// NewIndicesSegmentsService creates a new IndicesSegmentsService.
func NewIndicesSegmentsService(client *Client) *IndicesSegmentsService {
	return &IndicesSegmentsService{
		client: client,
		index:  make([]string, 0),
	}
}

// Index is documented as: A comma-separated list of index names; use `_all` or empty string to perform the operation on all indices.
func (s *IndicesSegmentsService) Index(indices ...string) *IndicesSegmentsService {
	s.index = append(s.index, indices...)
	return s
}

// AllowNoIndices is documented as: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes `_all` string or when no indices have been specified).
func (s *IndicesSegmentsService) AllowNoIndices(allowNoIndices bool) *IndicesSegmentsService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards is documented as: Whether to expand wildcard expression to concrete indices that are open, closed or both..
func (s *IndicesSegmentsService) ExpandWildcards(expandWildcards string) *IndicesSegmentsService {
	s.expandWildcards = expandWildcards
	return s
}

// IgnoreUnavailable is documented as: Whether specified concrete indices should be ignored when unavailable (missing or closed).
func (s *IndicesSegmentsService) IgnoreUnavailable(ignoreUnavailable bool) *IndicesSegmentsService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// OperationThreading is documented as: TODO: ?.
func (s *IndicesSegmentsService) OperationThreading(operationThreading interface{}) *IndicesSegmentsService {
	s.operationThreading = operationThreading
	return s
}

// Verbose is documented as: Includes detailed memory usage by Lucene..
func (s *IndicesSegmentsService) Verbose(verbose bool) *IndicesSegmentsService {
	s.verbose = &verbose
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *IndicesSegmentsService) Pretty(pretty bool) *IndicesSegmentsService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *IndicesSegmentsService) buildURL() (string, url.Values, error) {
	var err error
	var path string

	if len(s.index) > 0 {
		path, err = uritemplates.Expand("/{index}/_segments", map[string]string{
			"index": strings.Join(s.index, ","),
		})
	} else {
		path = "/_segments"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.allowNoIndices != nil {
		params.Set("allow_no_indices", fmt.Sprintf("%v", *s.allowNoIndices))
	}
	if s.expandWildcards != "" {
		params.Set("expand_wildcards", s.expandWildcards)
	}
	if s.ignoreUnavailable != nil {
		params.Set("ignore_unavailable", fmt.Sprintf("%v", *s.ignoreUnavailable))
	}
	if s.operationThreading != nil {
		params.Set("operation_threading", fmt.Sprintf("%v", s.operationThreading))
	}
	if s.verbose != nil {
		params.Set("verbose", fmt.Sprintf("%v", *s.verbose))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *IndicesSegmentsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *IndicesSegmentsService) Do(ctx context.Context) (*IndicesSegmentsResponse, error) {
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
	res, err := s.client.PerformRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(IndicesSegmentsResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// IndicesSegmentsResponse is the response of IndicesSegmentsService.Do.
type IndicesSegmentsResponse struct {
	// Shards provides information returned from shards.
	Shards shardsInfo `json:"_shards"`

	// Indices provides a map into the stats of an index. The key of the
	// map is the index name.
	Indices map[string]*IndexSegments `json:"indices,omitempty"`
}

type IndexSegments struct {
	// Shards provides a map into the shard related information of an index. The key of the
	// map is the number of a specific shard.
	Shards map[string][]*IndexSegmentsShards `json:"shards,omitempty"`
}

type IndexSegmentsShards struct {
	Routing              *IndexSegmentsRouting `json:"routing,omitempty"`
	NumCommittedSegments int64                 `json:"num_committed_segments,omitempty"`
	NumSearchSegments    int64                 `json:"num_search_segments"`

	// Segments provides a map into the segment related information of a shard. The key of the
	// map is the specific lucene segment id.
	Segments map[string]*IndexSegmentsDetails
}

type IndexSegmentsRouting struct {
	State   string `json:"state,omitempty"`
	Primary bool   `json:"primary,omitempty"`
	Node    string `json:"node,omitempty"`
}

type IndexSegmentsDetails struct {
	Generation    int64  `json:"generation,omitempty"`
	NumDocs       int64  `json:"num_docs,omitempty"`
	DeletedDocs   int64  `json:"deleted_docs,omitempty"`
	SizeInBytes   int64  `json:"size_in_bytes,omitempty"`
	MemoryInBytes int64  `json:"memory_in_bytes,omitempty"`
	Committed     bool   `json:"committed,omitempty"`
	Search        bool   `json:"search,omitempty"`
	Version       string `json:"version,omitempty"`
	Compound      bool   `json:"compound,omitempty"`
}
