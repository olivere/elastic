package elastic

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/atlassian/elastic/uritemplates"
	"net/http"
)

const (
	ClusterLevel = "cluster"
	IndicesLevel = "indices"
)

// FieldStatsService allows finding statistical properties of a field without executing a search,
// but looking up measurements that are natively available in the Lucene index.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/search-field-stats.html
// for details
type FieldStatsService struct {
	client  *Client
	level   string
	indices []string
	body    FieldStatsRequest
}

// NewFieldStatsService creates a new FieldStatsService
func NewFieldStatsService(client *Client) *FieldStatsService {
	return &FieldStatsService{
		client:  client,
		level:   ClusterLevel,
		indices: make([]string, 0),
		body: FieldStatsRequest{
			Fields: make([]string, 0),
		},
	}
}

// Indices sets the names of the indices to get stats for
func (s *FieldStatsService) Indices(indices ...string) *FieldStatsService {
	s.indices = append(s.indices, indices...)
	return s
}

// Level sets if stats should be returned on a per index level or on a cluster wide level;
// should be one of 'cluster' or 'indices'; defaults to former
func (s *FieldStatsService) Level(level string) *FieldStatsService {
	s.level = level
	return s
}

// Fields to compute and return field stats for
func (s *FieldStatsService) Fields(fields ...string) *FieldStatsService {
	s.body.Fields = append(s.body.Fields, fields...)
	return s
}

// IndexConstraints adds a field-level constraint; can be called multiple times
func (s *FieldStatsService) IndexConstraints(field string, constraints FieldConstraints) *FieldStatsService {
	if s.body.IndexConstraints == nil {
		s.body.IndexConstraints = make(map[string]FieldConstraints)
	}

	s.body.IndexConstraints[field] = constraints
	return s
}

func (s *FieldStatsService) buildURL() (string, url.Values, error) {
	// Build URL
	indices := strings.Join(s.indices, ",")
	path, err := uritemplates.Expand("{indices}/_field_stats", map[string]string{
		"indices": indices,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	if indices != "" {
		path = "/" + path
	}

	// Add query string parameters
	params := url.Values{}
	params.Set("level", s.level)

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *FieldStatsService) Validate() error {
	var invalid []string
	if s.level != IndicesLevel && s.level != ClusterLevel {
		invalid = append(invalid, "level")
	}

	if len(s.body.Fields) == 0 {
		invalid = append(invalid, "fields")
	}

	if len(invalid) != 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}

	return nil
}

// Do executes the operation.
func (s *FieldStatsService) Do() (*FieldStatsResponse, error) {
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
	res, err := s.client.PerformRequest("POST", path, params, s.body, http.StatusNotFound)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return &FieldStatsResponse{make(map[string]IndexFieldStats)}, nil
	}

	// Return operation response
	ret := new(FieldStatsResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// FieldStatsRequest is the request body content
type FieldStatsRequest struct {
	Fields           []string                    `json:"fields"`
	IndexConstraints map[string]FieldConstraints `json:"index_constraints,omitempty"`
}

// FieldConstraints is a constraint on a field
type FieldConstraints struct {
	Min *Comparisons `json:"min_value,omitempty"`
	Max *Comparisons `json:"max_value,omitempty"`
}

// Comparisons contain all comparison operations
type Comparisons struct {
	Lte string `json:"lte,omitempty"`
	Lt  string `json:"lt,omitempty"`
	Gte string `json:"gte,omitempty"`
	Gt  string `json:"gt,omitempty"`
}

// FieldStatsResponse is the response body content
type FieldStatsResponse struct {
	Indices map[string]IndexFieldStats `json:"indices,omitempty"`
}

// IndexFieldStats contains field stats for an index
type IndexFieldStats struct {
	Fields map[string]FieldStats `json:"fields,omitempty"`
}

// FieldStats contains stats of an individual  field
type FieldStats struct {
	MaxDoc                int64  `json:"max_doc"`
	DocCount              int64  `json:"doc_count"`
	Density               int64  `json:"density"`
	SumDocFrequeny        int64  `json:"sum_doc_freq"`
	SumTotalTermFrequency int64  `json:"sum_total_term_freq"`
	MinValue              string `json:"min_value_as_string"`
	MaxValue              string `json:"max_value_as_string"`
}
