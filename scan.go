// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

const (
	defaultKeepAlive = "5m"
)

var (
	// End of stream (or scan)
	EOS = errors.New("EOS")

	// No ScrollId
	ErrNoScrollId = errors.New("no scrollId")
)

// ScanService manages a cursor through documents in Elasticsearch.
type ScanService struct {
	client    *Client
	indices   []string
	types     []string
	keepAlive string
	fields    []string
	query     Query
	sorts     []SortInfo
	size      *int
	pretty    bool
}

func NewScanService(client *Client) *ScanService {
	builder := &ScanService{
		client: client,
		query:  NewMatchAllQuery(),
	}
	return builder
}

func (s *ScanService) Index(index string) *ScanService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, index)
	return s
}

func (s *ScanService) Indices(indices ...string) *ScanService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, indices...)
	return s
}

func (s *ScanService) Type(typ string) *ScanService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, typ)
	return s
}

func (s *ScanService) Types(types ...string) *ScanService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, types...)
	return s
}

// Scroll is an alias for KeepAlive, the time to keep
// the cursor alive (e.g. "5m" for 5 minutes).
func (s *ScanService) Scroll(keepAlive string) *ScanService {
	s.keepAlive = keepAlive
	return s
}

// KeepAlive sets the maximum time the cursor will be
// available before expiration (e.g. "5m" for 5 minutes).
func (s *ScanService) KeepAlive(keepAlive string) *ScanService {
	s.keepAlive = keepAlive
	return s
}

// Fields specifies the fields the scan query should load.
// By default fields is nil so _source is loaded
func (s *ScanService) Fields(fields ...string) *ScanService {
	s.fields = fields
	return s
}

func (s *ScanService) Query(query Query) *ScanService {
	s.query = query
	return s
}

// Sort the results by the given field, in the given order.
// Use the alternative SortWithInfo to use a struct to define the sorting.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-request-sort.html
// for detailed documentation of sorting.
func (s *ScanService) Sort(field string, ascending bool) *ScanService {
	s.sorts = append(s.sorts, SortInfo{Field: field, Ascending: ascending})
	return s
}

// SortWithInfo defines how to sort results.
// Use the Sort func for a shortcut.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-request-sort.html
// for detailed documentation of sorting.
func (s *ScanService) SortWithInfo(info SortInfo) *ScanService {
	s.sorts = append(s.sorts, info)
	return s
}

func (s *ScanService) Pretty(pretty bool) *ScanService {
	s.pretty = pretty
	return s
}

func (s *ScanService) Size(size int) *ScanService {
	s.size = &size
	return s
}

func (s *ScanService) Do() (*ScanCursor, error) {
	// Build url
	path := "/"

	// Indices part
	indexPart := make([]string, 0)
	for _, index := range s.indices {
		index, err := uritemplates.Expand("{index}", map[string]string{
			"index": index,
		})
		if err != nil {
			return nil, err
		}
		indexPart = append(indexPart, index)
	}
	if len(indexPart) > 0 {
		path += strings.Join(indexPart, ",")
	}

	// Types
	typesPart := make([]string, 0)
	for _, typ := range s.types {
		typ, err := uritemplates.Expand("{type}", map[string]string{
			"type": typ,
		})
		if err != nil {
			return nil, err
		}
		typesPart = append(typesPart, typ)
	}
	if len(typesPart) > 0 {
		path += "/" + strings.Join(typesPart, ",")
	}

	// Search
	path += "/_search"

	// Parameters
	params := make(url.Values)
	if len(s.sorts) == 0 {
		params.Set("search_type", "scan")
	}
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if s.keepAlive != "" {
		params.Set("scroll", s.keepAlive)
	} else {
		params.Set("scroll", defaultKeepAlive)
	}
	if s.size != nil && *s.size > 0 {
		params.Set("size", fmt.Sprintf("%d", *s.size))
	}
	if s.fields != nil {
		params.Set("fields", strings.Join(s.fields, ","))
	}

	// Set body
	body := make(map[string]interface{})
	if s.query != nil {
		body["query"] = s.query.Source()
	}
	if len(s.sorts) > 0 {
		sortarr := make([]interface{}, 0)
		for _, sort := range s.sorts {
			sortarr = append(sortarr, sort.Source())
		}
		body["sort"] = sortarr
	}

	// Get response
	res, err := s.client.PerformRequest("POST", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return result
	searchResult := new(SearchResult)
	if err := json.Unmarshal(res.Body, searchResult); err != nil {
		return nil, err
	}

	cursor := NewScanCursor(s.client, s.keepAlive, s.pretty, searchResult)

	return cursor, nil
}

// scanCursor represents a single page of results from
// an Elasticsearch Scan operation.
type ScanCursor struct {
	Results *SearchResult

	client      *Client
	keepAlive   string
	pretty      bool
	currentPage int
}

// newScanCursor returns a new initialized instance
// of scanCursor.
func NewScanCursor(client *Client, keepAlive string, pretty bool, searchResult *SearchResult) *ScanCursor {
	return &ScanCursor{
		client:    client,
		keepAlive: keepAlive,
		pretty:    pretty,
		Results:   searchResult,
	}
}

// TotalHits is a convenience method that returns the number
// of hits the cursor will iterate through.
func (c *ScanCursor) TotalHits() int64 {
	if c.Results.Hits == nil {
		return 0
	}
	return c.Results.Hits.TotalHits
}

// Next returns the next search result or nil when all
// documents have been scanned.
//
// Usage:
//
//   for {
//     res, err := cursor.Next()
//     if err == elastic.EOS {
//       // End of stream (or scan)
//       break
//     }
//     if err != nil {
//       // Handle error
//     }
//     // Work with res
//   }
//
func (c *ScanCursor) Next() (*SearchResult, error) {
	if c.currentPage > 0 {
		if c.Results.Hits == nil || len(c.Results.Hits.Hits) == 0 || c.Results.Hits.TotalHits == 0 {
			return nil, EOS
		}
	}
	if c.Results.ScrollId == "" {
		return nil, EOS
	}

	// Build url
	path := "/_search/scroll"

	// Parameters
	params := make(url.Values)
	if c.pretty {
		params.Set("pretty", fmt.Sprintf("%v", c.pretty))
	}
	if c.keepAlive != "" {
		params.Set("scroll", c.keepAlive)
	} else {
		params.Set("scroll", defaultKeepAlive)
	}

	// Set body
	body := c.Results.ScrollId

	// Get response
	res, err := c.client.PerformRequest("POST", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return result
	if err := json.Unmarshal(res.Body, c.Results); err != nil {
		return nil, err
	}

	c.currentPage += 1

	return c.Results, nil
}
