// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Search for documents in Elasticsearch.
type SearchService struct {
	client       *Client
	searchSource *SearchSource
	pretty       bool
	searchType   string
	indices      []string
	queryHint    string
	routing      string
	preference   string
	types        []string
	debug        bool
}

func NewSearchService(client *Client) *SearchService {
	builder := &SearchService{
		client:       client,
		searchSource: NewSearchSource(),
		debug:        false,
		pretty:       false,
	}
	return builder
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

func (s *SearchService) Type(typ string) *SearchService {
	if s.types == nil {
		s.types = []string{typ}
	} else {
		s.types = append(s.types, typ)
	}
	return s
}

func (s *SearchService) Types(types ...string) *SearchService {
	if s.types == nil {
		s.types = make([]string, len(types))
	}
	s.types = append(s.types, types...)
	return s
}

func (s *SearchService) Pretty(pretty bool) *SearchService {
	s.pretty = pretty
	return s
}

func (s *SearchService) Debug(debug bool) *SearchService {
	s.debug = debug
	return s
}

func (s *SearchService) Timeout(timeout string) *SearchService {
	s.searchSource = s.searchSource.Timeout(timeout)
	return s
}

func (s *SearchService) TimeoutInMillis(timeoutInMillis int) *SearchService {
	s.searchSource = s.searchSource.TimeoutInMillis(timeoutInMillis)
	return s
}

func (s *SearchService) SearchType(searchType string) *SearchService {
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

func (s *SearchService) Query(query Query) *SearchService {
	s.searchSource = s.searchSource.Query(query)
	return s
}

// PostFilter is executed as the last filter. It only affects the
// search hits but not facets.
func (s *SearchService) PostFilter(postFilter Filter) *SearchService {
	s.searchSource = s.searchSource.PostFilter(postFilter)
	return s
}

func (s *SearchService) Highlight(highlight *Highlight) *SearchService {
	s.searchSource = s.searchSource.Highlight(highlight)
	return s
}

func (s *SearchService) GlobalSuggestText(globalText string) *SearchService {
	s.searchSource = s.searchSource.GlobalSuggestText(globalText)
	return s
}

func (s *SearchService) Suggester(suggester Suggester) *SearchService {
	s.searchSource = s.searchSource.Suggester(suggester)
	return s
}

func (s *SearchService) Facet(name string, facet Facet) *SearchService {
	s.searchSource = s.searchSource.Facet(name, facet)
	return s
}

func (s *SearchService) Aggregation(name string, aggregation Aggregation) *SearchService {
	s.searchSource = s.searchSource.Aggregation(name, aggregation)
	return s
}

func (s *SearchService) MinScore(minScore float64) *SearchService {
	s.searchSource = s.searchSource.MinScore(minScore)
	return s
}

func (s *SearchService) From(from int) *SearchService {
	s.searchSource = s.searchSource.From(from)
	return s
}

func (s *SearchService) Size(size int) *SearchService {
	s.searchSource = s.searchSource.Size(size)
	return s
}

func (s *SearchService) Explain(explain bool) *SearchService {
	s.searchSource = s.searchSource.Explain(explain)
	return s
}

func (s *SearchService) Version(version bool) *SearchService {
	s.searchSource = s.searchSource.Version(version)
	return s
}

func (s *SearchService) Sort(field string, ascending bool) *SearchService {
	s.searchSource = s.searchSource.Sort(field, ascending)
	return s
}

func (s *SearchService) SortWithInfo(info SortInfo) *SearchService {
	s.searchSource = s.searchSource.SortWithInfo(info)
	return s
}

func (s *SearchService) Fields(fields ...string) *SearchService {
	s.searchSource = s.searchSource.Fields(fields...)
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

	// Types part
	if len(s.types) > 0 {
		typesPart := make([]string, 0)
		for _, typ := range s.types {
			typesPart = append(typesPart, cleanPathString(typ))
		}
		urls += "/"
		urls += strings.Join(typesPart, ",")
	}

	// Search
	urls += "/_search"

	// Parameters
	params := make(url.Values)
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if s.searchType != "" {
		params.Set("search_type", s.searchType)
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	req.SetBodyJson(s.searchSource.Source())

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
	TookInMillis int64                      `json:"took"`
	ScrollId     string                     `json:"_scroll_id"`
	Hits         *SearchHits                `json:"hits"`
	Suggest      SearchSuggest              `json:"suggest"`
	Facets       SearchFacets               `json:"facets"`
	Aggregations map[string]json.RawMessage `json:"aggregations"` // see search_aggs.go
	TimedOut     bool                       `json:"timed_out"`
	Error        string                     `json:"error,omitempty"` // used in MultiSearch only
}

// GetAggregation returns the aggregation with the specified name.
func (res *SearchResult) GetAggregation(name string) (*SearchAggregation, bool) {
	agg, found := res.Aggregations[name]
	if !found {
		return nil, false
	}
	return NewSearchAggregation(name, agg), true
}

type SearchHits struct {
	TotalHits int64        `json:"total"`
	MaxScore  *float64     `json:"max_score"`
	Hits      []*SearchHit `json:"hits"`
}

type SearchHit struct {
	Score     *float64               `json:"_score"`
	Index     string                 `json:"_index"`
	Id        string                 `json:"_id"`
	Type      string                 `json:"_type"`
	Version   *int64                 `json:"_version"`
	Sort      *[]interface{}         `json:"sort"`
	Highlight SearchHitHighlight     `json:"highlight"`
	Source    *json.RawMessage       `json:"_source"`
	Fields    map[string]interface{} `json:"fields"`

	// Explanation
	// Shard
	// HighlightFields
	// SortValues
	// MatchedFilters
}

// Suggest

type SearchSuggest map[string][]SearchSuggestion

type SearchSuggestion struct {
	Text    string                   `json:"text"`
	Offset  int                      `json:"offset"`
	Length  int                      `json:"length"`
	Options []SearchSuggestionOption `json:"options"`
}

type SearchSuggestionOption struct {
	Text    string      `json:"text"`
	Score   float32     `json:"score"`
	Freq    int         `json:"freq"`
	Payload interface{} `json:"payload"`
}

// Facets

type SearchFacets map[string]*SearchFacet

type SearchFacet struct {
	Type    string             `json:"_type"`
	Missing int                `json:"missing"`
	Total   int                `json:"total"`
	Other   int                `json:"other"`
	Terms   []searchFacetTerm  `json:"terms"`
	Ranges  []searchFacetRange `json:"ranges"`
	Entries []searchFacetEntry `json:"entries"`
}

type searchFacetTerm struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type searchFacetRange struct {
	From       *float64 `json:"from"`
	FromStr    *string  `json:"from_str"`
	To         *float64 `json:"to"`
	ToStr      *string  `json:"to_str"`
	Count      int      `json:"count"`
	Min        *float64 `json:"min"`
	Max        *float64 `json:"max"`
	TotalCount int      `json:"total_count"`
	Total      *float64 `json:"total"`
	Mean       *float64 `json:"mean"`
}

type searchFacetEntry struct {
	// Key for this facet, e.g. in histograms
	Key interface{} `json:"key"`
	// Date histograms contain the number of milliseconds as date:
	// If e.Time = 1293840000000, then: Time.at(1293840000000/1000) => 2011-01-01
	Time int64 `json:"time"`
	// Number of hits for this facet
	Count int `json:"count"`
	// Min is either a string like "Infinity" or a float64.
	// This is returned with some DateHistogram facets.
	Min interface{} `json:"min,omitempty"`
	// Max is either a string like "-Infinity" or a float64
	// This is returned with some DateHistogram facets.
	Max interface{} `json:"max,omitempty"`
	// Total is the sum of all entries on the recorded Time
	// This is returned with some DateHistogram facets.
	Total float64 `json:"total,omitempty"`
	// TotalCount is the number of entries for Total
	// This is returned with some DateHistogram facets.
	TotalCount int `json:"total_count,omitempty"`
	// Mean is the mean value
	// This is returned with some DateHistogram facets.
	Mean float64 `json:"mean,omitempty"`
}

// Aggregations (see search_aggs.go)

// Highlighting

type SearchHitHighlight map[string][]string
