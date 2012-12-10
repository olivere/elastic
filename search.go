// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// Search types
type SearchType int

const (
	DfsQueryThenFetch SearchType = iota // 0
	QueryThenFetch                      // 1 (=default)
	DfsQueryAndFetch                    // 2
	QueryAndFetch                       // 3
	Scan                                // 4
	Count                               // 5
)

// Search for documents in ElasticSearch.
type SearchService struct {
	client     *Client
	searchType SearchType
	indices    []string
	queryHint  string
	routing    string
	preference string
	types      []string
	timeout    *time.Duration
	query      Query
	filters    []Filter
	minScore   *float64
	from       *int
	size       *int
	explain    *bool
	version    *bool
	sorts      map[string]bool
	fields     []string
	facets     []Facet
	debug      bool
}

func NewSearchService(client *Client) *SearchService {
	builder := &SearchService{
		client:     client,
		searchType: QueryThenFetch,
		filters:    make([]Filter, 0),
		sorts:      make(map[string]bool),
		fields:     make([]string, 0),
		facets:     make([]Facet, 0),
		debug:      false,
	}
	return builder
}

func (s *SearchService) Debug(debug bool) *SearchService {
	s.debug = debug
	return s
}

func (s *SearchService) Index(index string) *SearchService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, index)
	return s
}

func (s *SearchService) Indices(indices ...string) *SearchService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, indices...)
	return s
}

func (s *SearchService) Type(_type string) *SearchService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, _type)
	return s
}

func (s *SearchService) Types(types ...string) *SearchService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, types...)
	return s
}

func (s *SearchService) SearchType(searchType SearchType) *SearchService {
	s.searchType = searchType
	return s
}

func (s *SearchService) Routing(routing string) *SearchService {
	s.routing = routing
	return s
}

func (s *SearchService) Preference(preference string) *SearchService {
	s.preference = preference
	return s
}

func (s *SearchService) QueryHint(queryHint string) *SearchService {
	s.queryHint = queryHint
	return s
}

func (s *SearchService) Timeout(timeout *time.Duration) *SearchService {
	s.timeout = timeout
	return s
}

func (s *SearchService) Query(query Query) *SearchService {
	s.query = query
	return s
}

func (s *SearchService) AddFilter(filter Filter) *SearchService {
	s.filters = append(s.filters, filter)
	return s
}

func (s *SearchService) AddFacet(facet Facet) *SearchService {
	s.facets = append(s.facets, facet)
	return s
}

func (s *SearchService) MinScore(minScore float64) *SearchService {
	s.minScore = &minScore
	return s
}

func (s *SearchService) From(from int) *SearchService {
	s.from = &from
	return s
}

func (s *SearchService) Size(size int) *SearchService {
	s.size = &size
	return s
}

func (s *SearchService) Explain(explain bool) *SearchService {
	s.explain = &explain
	return s
}

func (s *SearchService) Version(version bool) *SearchService {
	s.version = &version
	return s
}

func (s *SearchService) AddSort(field string, ascending bool) *SearchService {
	s.sorts[field] = ascending
	return s
}

func (s *SearchService) Fields(fields ...string) *SearchService {
	s.fields = append(s.fields, fields...)
	return s
}

func (s *SearchService) Do() (*SearchResult, error) {
	// Build url
	urls := "/"

	// Indices part
	indexPart := make([]string, 0)
	for _, index := range s.indices {
		indexPart = append(indexPart, cleanPathString(index))
	}
	urls += strings.Join(indexPart, ",")

	// TODO Types part

	// Search
	urls += "/_search"

	// Set up a new request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	body := make(map[string]interface{})
	if s.query != nil {
		body["query"] = s.query.Source()
	}
	if s.from != nil && *s.from > 0 {
		body["from"] = *s.from
	}
	if s.size != nil && *s.size > 0 {
		body["size"] = *s.size
	}
	if len(s.filters) == 1 {
		body["filter"] = s.filters[0].Source()
	} else if len(s.filters) > 1 {
		f := make(map[string]interface{})
		andedFilters := make([]interface{}, 0)
		for _, filter := range s.filters {
			andedFilters = append(andedFilters, filter.Source())
		}
		f["and"] = andedFilters
		body["filter"] = f
	}
	req.SetBodyJson(body)

	if s.debug {
		out, _ := httputil.DumpRequestOut((*http.Request)(req), true)
		fmt.Printf("%s\n", string(out))
	}

	// Get response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if s.debug {
		out, _ := httputil.DumpResponse(res, true)
		fmt.Printf("%s\n", string(out))
	}

	ret := new(SearchResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type SearchResult struct {
	TookInMillis int64       `json:"took"`
	ScrollId     string      `json:"_scroll_id,omitempty"`
	Hits         *SearchHits `json:"hits"`
	Facets       *Facets     `json:"facets"`
	TimedOut     bool        `json:"timed_out"`
}

type SearchHits struct {
	TotalHits int64        `json:"total"`
	MaxScore  *float64     `json:"max_score,omitempty"`
	Hits      []*SearchHit `json:"hits"`
}

type SearchHit struct {
	Score   float64          `json:"_score"`
	Index   string           `json:"_index"`
	Id      string           `json:"_id"`
	Type    string           `json:"_type"`
	Version int64            `json:"_version"`
	Source  *json.RawMessage `json:"_source"`

	// Explanation
	// Shard
	// HighlightFields
	// SortValues
	// MatchedFilters
}
