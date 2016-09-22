// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"gopkg.in/olivere/elastic.v5/uritemplates"
)

// DeleteByQueryService deletes documents that match a query.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/docs-delete-by-query.html.
type DeleteByQueryService struct {
	client            *Client
	index             []string
	typ               []string
	analyzer          string
	consistency       string
	defaultOper       string
	df                string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
	replication       string
	routing           string
	timeout           string
	pretty            bool
	q                 string
	query             Query
}

// NewDeleteByQueryService creates a new DeleteByQueryService.
// You typically use the client's DeleteByQuery to get a reference to
// the service.
func NewDeleteByQueryService(client *Client) *DeleteByQueryService {
	builder := &DeleteByQueryService{
		client: client,
	}
	return builder
}

// Index sets the indices on which to perform the delete operation.
func (s *DeleteByQueryService) Index(index ...string) *DeleteByQueryService {
	s.index = append(s.index, index...)
	return s
}

// Type limits the delete operation to the given types.
func (s *DeleteByQueryService) Type(typ ...string) *DeleteByQueryService {
	s.typ = append(s.typ, typ...)
	return s
}

// Analyzer to use for the query string.
func (s *DeleteByQueryService) Analyzer(analyzer string) *DeleteByQueryService {
	s.analyzer = analyzer
	return s
}

// Consistency represents the specific write consistency setting for the operation.
// It can be one, quorum, or all.
func (s *DeleteByQueryService) Consistency(consistency string) *DeleteByQueryService {
	s.consistency = consistency
	return s
}

// DefaultOperator for query string query (AND or OR).
func (s *DeleteByQueryService) DefaultOperator(defaultOperator string) *DeleteByQueryService {
	s.defaultOper = defaultOperator
	return s
}

// DF is the field to use as default where no field prefix is given in the query string.
func (s *DeleteByQueryService) DF(defaultField string) *DeleteByQueryService {
	s.df = defaultField
	return s
}

// DefaultField is the field to use as default where no field prefix is given in the query string.
// It is an alias to the DF func.
func (s *DeleteByQueryService) DefaultField(defaultField string) *DeleteByQueryService {
	s.df = defaultField
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *DeleteByQueryService) IgnoreUnavailable(ignore bool) *DeleteByQueryService {
	s.ignoreUnavailable = &ignore
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices (including the _all string
// or when no indices have been specified).
func (s *DeleteByQueryService) AllowNoIndices(allow bool) *DeleteByQueryService {
	s.allowNoIndices = &allow
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both. It can be "open" or "closed".
func (s *DeleteByQueryService) ExpandWildcards(expand string) *DeleteByQueryService {
	s.expandWildcards = expand
	return s
}

// Replication sets a specific replication type (sync or async).
func (s *DeleteByQueryService) Replication(replication string) *DeleteByQueryService {
	s.replication = replication
	return s
}

// Q specifies the query in Lucene query string syntax. You can also use
// Query to programmatically specify the query.
func (s *DeleteByQueryService) Q(query string) *DeleteByQueryService {
	s.q = query
	return s
}

// QueryString is an alias to Q. Notice that you can also use Query to
// programmatically set the query.
func (s *DeleteByQueryService) QueryString(query string) *DeleteByQueryService {
	s.q = query
	return s
}

// Routing sets a specific routing value.
func (s *DeleteByQueryService) Routing(routing string) *DeleteByQueryService {
	s.routing = routing
	return s
}

// Timeout sets an explicit operation timeout, e.g. "1s" or "10000ms".
func (s *DeleteByQueryService) Timeout(timeout string) *DeleteByQueryService {
	s.timeout = timeout
	return s
}

// Pretty indents the JSON output from Elasticsearch.
func (s *DeleteByQueryService) Pretty(pretty bool) *DeleteByQueryService {
	s.pretty = pretty
	return s
}

// Query sets the query programmatically.
func (s *DeleteByQueryService) Query(query Query) *DeleteByQueryService {
	s.query = query
	return s
}

// buildURL builds the URL for the operation.
func (s *DeleteByQueryService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string
	if len(s.typ) > 0 {
		path, err = uritemplates.Expand("/{index}/{type}/_delete_by_query", map[string]string{
			"index": strings.Join(s.index, ","),
			"type":  strings.Join(s.typ, ","),
		})
	} else {
		path, err = uritemplates.Expand("/{index}/_delete_by_query", map[string]string{
			"index": strings.Join(s.index, ","),
		})
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.analyzer != "" {
		params.Set("analyzer", s.analyzer)
	}
	if s.consistency != "" {
		params.Set("consistency", s.consistency)
	}
	if s.defaultOper != "" {
		params.Set("default_operator", s.defaultOper)
	}
	if s.df != "" {
		params.Set("df", s.df)
	}
	if s.ignoreUnavailable != nil {
		params.Set("ignore_unavailable", fmt.Sprintf("%v", *s.ignoreUnavailable))
	}
	if s.allowNoIndices != nil {
		params.Set("allow_no_indices", fmt.Sprintf("%v", *s.allowNoIndices))
	}
	if s.expandWildcards != "" {
		params.Set("expand_wildcards", s.expandWildcards)
	}
	if s.replication != "" {
		params.Set("replication", s.replication)
	}
	if s.routing != "" {
		params.Set("routing", s.routing)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if s.q != "" {
		params.Set("q", s.q)
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *DeleteByQueryService) Validate() error {
	var invalid []string
	if len(s.index) == 0 {
		invalid = append(invalid, "Index")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the delete-by-query operation.
func (s *DeleteByQueryService) Do(ctx context.Context) (*BulkIndexByScrollResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Set body if there is a query set
	var body interface{}
	if s.query != nil {
		src, err := s.query.Source()
		if err != nil {
			return nil, err
		}
		query := make(map[string]interface{})
		query["query"] = src
		body = query
	}

	// Get response
	res, err := s.client.PerformRequest(ctx, "POST", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return result
	ret := new(BulkIndexByScrollResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// BulkIndexByScrollResponse is the outcome of executing Do with
// DeleteByQueryService and UpdateByQueryService.
type BulkIndexByScrollResponse struct {
	Took             int64 `json:"took"`
	TimedOut         bool  `json:"timed_out"`
	Total            int64 `json:"total"`
	Updated          int64 `json:"updated"`
	Created          int64 `json:"created"`
	Deleted          int64 `json:"deleted"`
	Batches          int64 `json:"batches"`
	VersionConflicts int64 `json:"version_conflicts"`
	Noops            int64 `json:"noops"`
	Retries          struct {
		Bulk   int64 `json:"bulk"`
		Search int64 `json:"search"`
	} `json:"retries"`
	Throttled            string                             `json:"throttled"`
	ThrottledMillis      int64                              `json:"throttled_millis"`
	RequestsPerSecond    float64                            `json:"requests_per_second"`
	Canceled             string                             `json:"canceled"`
	ThrottledUntil       string                             `json:"throttled_until"`
	ThrottledUntilMillis int64                              `json:"throttled_until_millis"`
	Failures             []bulkIndexByScrollResponseFailure `json:"failures"`
}

type bulkIndexByScrollResponseFailure struct {
	Index  string `json:"index,omitempty"`
	Type   string `json:"type,omitempty"`
	Id     string `json:"id,omitempty"`
	Status int    `json:"status,omitempty"`
	Shard  int    `json:"shard,omitempty"`
	Node   int    `json:"node,omitempty"`
	// TOOD "cause" contains exception details
	// TOOD "reason" contains exception details
}
