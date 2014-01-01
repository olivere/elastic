// Copyright 2012 Oliver Eilhard. All rights reserved.
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

// Information about sorting a field.
type SortInfo struct {
	Field          string
	Ascending      bool
	Missing        *interface{}
	IgnoreUnmapped bool
}

// Search for documents in ElasticSearch.
type SearchService struct {
	client            *Client
	pretty            bool
	searchType        string
	indices           []string
	queryHint         string
	routing           string
	preference        string
	types             []string
	timeout           string
	query             Query
	filters           []Filter
	highlight         *Highlight
	globalSuggestText string
	suggesters        []Suggester
	minScore          *float64
	from              *int
	size              *int
	explain           *bool
	version           *bool
	sorts             []SortInfo
	fields            []string
	facets            map[string]Facet
	debug             bool
}

func NewSearchService(client *Client) *SearchService {
	builder := &SearchService{
		client:  client,
		filters: make([]Filter, 0),
		sorts:   make([]SortInfo, 0),
		fields:  make([]string, 0),
		facets:  make(map[string]Facet, 0),
		debug:   false,
		pretty:  false,
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
	s.timeout = timeout
	return s
}

func (s *SearchService) TimeoutInMillis(timeoutInMillis int) *SearchService {
	s.timeout = fmt.Sprintf("%dms", timeoutInMillis)
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
	s.query = query
	return s
}

func (s *SearchService) Filter(filter Filter) *SearchService {
	s.filters = append(s.filters, filter)
	return s
}

func (s *SearchService) Highlight(highlight *Highlight) *SearchService {
	s.highlight = highlight
	return s
}

func (s *SearchService) GlobalSuggestText(globalText string) *SearchService {
	s.globalSuggestText = globalText
	return s
}

func (s *SearchService) Suggester(suggester Suggester) *SearchService {
	s.suggesters = append(s.suggesters, suggester)
	return s
}

func (s *SearchService) Facet(name string, facet Facet) *SearchService {
	s.facets[name] = facet
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

func (s *SearchService) Sort(field string, ascending bool) *SearchService {
	s.sorts = append(s.sorts, SortInfo{Field: field, Ascending: ascending})
	return s
}

func (s *SearchService) SortWithInfo(info SortInfo) *SearchService {
	s.sorts = append(s.sorts, info)
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
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
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
	body := make(map[string]interface{})

	// Query
	if s.query != nil {
		body["query"] = s.query.Source()
	}

	// Filters
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

	// Highlight
	if s.highlight != nil {
		body["highlight"] = s.highlight.Source()
	}

	// Suggesters
	if len(s.suggesters) > 0 {
		suggesters := make(map[string]interface{})

		for _, s := range s.suggesters {
			suggesters[s.Name()] = s.Source(false)
		}

		if s.globalSuggestText != "" {
			suggesters["text"] = s.globalSuggestText
		}

		body["suggest"] = suggesters
	}

	// Facets
	if len(s.facets) >= 1 {
		// "facets" : {
		//   "manufacturer" : {
		//     "terms" : { ... }
		//   },
		//   "price" : {
		//     "range" : { ... }
		//   }
		// }
		facetsMap := make(map[string]interface{})
		body["facets"] = facetsMap

		for field, facet := range s.facets {
			facetsMap[field] = facet.Source()
		}
	}

	// Limit/Offset
	if s.from != nil && *s.from > 0 {
		body["from"] = *s.from
	}
	if s.size != nil && *s.size > 0 {
		body["size"] = *s.size
	}

	// Sort
	if len(s.sorts) > 0 {
		sortSlice := make([]interface{}, 0)
		for _, info := range s.sorts {
			sortProp := make(map[string]interface{})
			if info.Ascending {
				sortProp["order"] = "asc"
			} else {
				sortProp["order"] = "desc"
			}
			sortElem := make(map[string]interface{})
			sortElem[info.Field] = sortProp
			sortSlice = append(sortSlice, sortElem)
		}
		body["sort"] = sortSlice
	}

	// Fields
	if len(s.fields) > 0 {
		body["fields"] = s.fields
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
	TookInMillis int64         `json:"took"`
	ScrollId     string        `json:"_scroll_id"`
	Hits         *SearchHits   `json:"hits"`
	Suggest      SearchSuggest `json:"suggest"`
	Facets       SearchFacets  `json:"facets"`
	TimedOut     bool          `json:"timed_out"`
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
	Text  string  `json:"text"`
	Score float32 `json:"score"`
	Freq  int     `json:"freq"`
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

type SearchHitHighlight map[string][]string
