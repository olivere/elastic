// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// CatIndicesService returns the list of indices plus some additional
// information about them.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.2/cat-indices.html
// for details.
type CatIndicesService struct {
	client        *Client
	pretty        bool
	index         string
	bytes         string // b, k, m, or g
	local         *bool
	masterTimeout string
	columns       []string
	health        string   // green, yellow, or red
	primaryOnly   *bool    // true for primary shards only
	sort          []string // list of columns for sort order
}

// NewCatIndicesService creates a new CatIndicesService.
func NewCatIndicesService(client *Client) *CatIndicesService {
	return &CatIndicesService{
		client: client,
	}
}

// Index is the name of the index to list (by default all indices are returned).
func (s *CatIndicesService) Index(index string) *CatIndicesService {
	s.index = index
	return s
}

// Bytes represents the unit in which to display byte values.
// Valid values are: "b", "k", "m", or "g".
func (s *CatIndicesService) Bytes(bytes string) *CatIndicesService {
	s.bytes = bytes
	return s
}

// Local indicates to return local information, i.e. do not retrieve
// the state from master node (default: false).
func (s *CatIndicesService) Local(local bool) *CatIndicesService {
	s.local = &local
	return s
}

// MasterTimeout is the explicit operation timeout for connection to master node.
func (s *CatIndicesService) MasterTimeout(masterTimeout string) *CatIndicesService {
	s.masterTimeout = masterTimeout
	return s
}

// Columns to return in the response.
// To get a list of all possible columns to return, run the following command
// in your terminal:
//
// Example:
//   curl 'http://localhost:9200/_cat/indices?help'
//
// You can use Columns("*") to return all possible columns. That might take
// a little longer than the default set of columns.
func (s *CatIndicesService) Columns(columns ...string) *CatIndicesService {
	s.columns = columns
	return s
}

// Health filters indices by their health status.
// Valid values are: "green", "yellow", or "red".
func (s *CatIndicesService) Health(healthState string) *CatIndicesService {
	s.health = healthState
	return s
}

// PrimaryOnly when set to true returns stats only for primary shards (default: false).
func (s *CatIndicesService) PrimaryOnly(primaryOnly bool) *CatIndicesService {
	s.primaryOnly = &primaryOnly
	return s
}

// Sort is a list of fields to sort by.
func (s *CatIndicesService) Sort(fields ...string) *CatIndicesService {
	s.sort = fields
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatIndicesService) Pretty(pretty bool) *CatIndicesService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *CatIndicesService) buildURL() (string, url.Values, error) {
	// Build URL
	var (
		path string
		err  error
	)

	if s.index != "" {
		path, err = uritemplates.Expand("/_cat/indices/{index}", map[string]string{
			"index": s.index,
		})
	} else {
		path = "/_cat/indices"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{
		"format": []string{"json"}, // always returns as JSON
	}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.bytes != "" {
		params.Set("bytes", s.bytes)
	}
	if v := s.local; v != nil {
		params.Set("local", fmt.Sprint(*v))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if len(s.columns) > 0 {
		params.Set("h", strings.Join(s.columns, ","))
	}
	if s.health != "" {
		params.Set("health", s.health)
	}
	if v := s.primaryOnly; v != nil {
		params.Set("pri", fmt.Sprint(*v))
	}
	if len(s.sort) > 0 {
		params.Set("s", strings.Join(s.sort, ","))
	}
	return path, params, nil
}

// Do executes the operation.
func (s *CatIndicesService) Do(ctx context.Context) (CatIndicesResponse, error) {
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
	var ret CatIndicesResponse
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a get request.

// CatIndicesResponse is the outcome of CatIndicesService.Do.
type CatIndicesResponse []CatIndicesResponseRow

// CatIndicesResponseRow specifies the data returned for one index
// of a CatIndicesResponse. Notice that not all of these fields might
// be filled; that depends on the number of columns chose in the
// request (see CatIndicesService.Columns).
type CatIndicesResponseRow struct {
	Health                       string `json:"health"`                              // "green", "yellow", or "red"
	Status                       string `json:"status"`                              // "open" or "closed"
	Index                        string `json:"index"`                               // index name
	UUID                         string `json:"uuid"`                                // index uuid
	Pri                          int    `json:"pri,string"`                          // number of primary shards
	Rep                          int    `json:"rep,string"`                          // number of replica shards
	DocsCount                    int    `json:"docs.count,string"`                   // number of available documents
	DocsDeleted                  int    `json:"docs.deleted,string"`                 // number of deleted documents
	CreationDate                 int64  `json:"creation.date,string"`                // index creation date (millisecond value), e.g. 1527077221644
	CreationDateString           string `json:"creation.date.string"`                // index creation date (as string), e.g. "2018-05-23T12:07:01.644Z"
	StoreSize                    string `json:"store.size"`                          // store size of primaries & replicas, e.g. "4.6kb"
	PriStoreSize                 string `json:"pri.store.size"`                      // store size of primaries, e.g. "230b"
	CompletionSize               string `json:"completion.size"`                     // size of completion on primaries & replicas
	PriCompletionSize            string `json:"pri.completion.size"`                 // size of completion on primaries
	FielddataMemorySize          string `json:"fielddata.memory_size"`               // used fielddata cache on primaries & replicas
	PriFielddataMemorySize       string `json:"pri.fielddata.memory_size"`           // used fielddata cache on primaries
	FielddataEvictions           int    `json:"fielddata.evictions,string"`          // fielddata evictions on primaries & replicas
	PriFielddataEvictions        int    `json:"pri.fielddata.evictions,string"`      // fielddata evictions on primaries
	QueryCacheMemorySize         string `json:"query_cache.memory_size"`             // used query cache on primaries & replicas
	PriQueryCacheMemorySize      string `json:"pri.query_cache.memory_size"`         // used query cache on primaries
	QueryCacheEvictions          int    `json:"query_cache.evictions,string"`        // query cache evictions on primaries & replicas
	PriQueryCacheEvictions       int    `json:"pri.query_cache.evictions,string"`    // query cache evictions on primaries
	RequestCacheMemorySize       string `json:"request_cache.memory_size"`           // used request cache on primaries & replicas
	PriRequestCacheMemorySize    string `json:"pri.request_cache.memory_size"`       // used request cache on primaries
	RequestCacheEvictions        int    `json:"request_cache.evictions,string"`      // request cache evictions on primaries & replicas
	PriRequestCacheEvictions     int    `json:"pri.request_cache.evictions,string"`  // request cache evictions on primaries
	RequestCacheHitCount         int    `json:"request_cache.hit_count,string"`      // request cache hit count on primaries & replicas
	PriRequestCacheHitCount      int    `json:"pri.request_cache.hit_count,string"`  // request cache hit count on primaries
	RequestCacheMissCount        int    `json:"request_cache.miss_count,string"`     // request cache miss count on primaries & replicas
	PriRequestCacheMissCount     int    `json:"pri.request_cache.miss_count,string"` // request cache miss count on primaries
	FlushTotal                   int    `json:"flush.total"`                         // number of flushes on primaries & replicas
	PriFlushTotal                int    `json:"pri.flush.total"`                     // number of flushes on primaries
	FlushTotalTime               string `json:"flush.total_time"`                    // time spent in flush on primaries & replicas
	PriFlushTotalTime            string `json:"pri.flush.total_time"`                // time spent in flush on primaries
	GetCurrent                   int    `json:"get.current,string"`                  // number of current get ops on primaries & replicas
	PriGetCurrent                int    `json:"pri.get.current,string"`              // number of current get ops on primaries
	GetTime                      string `json:"get.time"`                            // time spent in get on primaries & replicas
	PriGetTime                   string `json:"pri.get.time"`                        // time spent in get on primaries
	GetTotal                     int    `json:"get.total,string"`                    // number of get ops on primaries & replicas
	PriGetTotal                  int    `json:"pri.get.total,string"`                // number of get ops on primaries
	GetExistsTime                string `json:"get.exists_time"`                     // time spent in successful gets on primaries & replicas
	PriGetExistsTime             string `json:"pri.get.exists_time"`                 // time spent in successful gets on primaries
	GetExistsTotal               int    `json:"get.exists_total,string"`             // number of successful gets on primaries & replicas
	PriGetExistsTotal            int    `json:"pri.get.exists_total,string"`         // number of successful gets on primaries
	GetMissingTime               string `json:"get.missing_time"`                    // time spent in failed gets on primaries & replicas
	PriGetMissingTime            string `json:"pri.get.missing_time"`                // time spent in failed gets on primaries
	GetMissingTotal              int    `json:"get.missing_total,string"`            // number of failed gets on primaries & replicas
	PriGetMissingTotal           int    `json:"pri.get.missing_total,string"`        // number of failed gets on primaries
	IndexingDeleteCurrent        int    `json:"indexing.delete_current,string"`      // number of current deletions on primaries & replicas
	PriIndexingDeleteCurrent     int    `json:"pri.indexing.delete_current,string"`  // number of current deletions on primaries
	IndexingDeleteTime           string `json:"indexing.delete_time"`                // time spent in deletions on primaries & replicas
	PriIndexingDeleteTime        string `json:"pri.indexing.delete_time"`            // time spent in deletions on primaries
	IndexingDeleteTotal          int    `json:"indexing.delete_total,string"`        // number of delete ops on primaries & replicas
	PriIndexingDeleteTotal       int    `json:"pri.indexing.delete_total,string"`    // number of delete ops on primaries
	IndexingIndexCurrent         int    `json:"indexing.index_current,string"`       // number of current indexing on primaries & replicas
	PriIndexingIndexCurrent      int    `json:"pri.indexing.index_current,string"`   // number of current indexing on primaries
	IndexingIndexTime            string `json:"indexing.index_time"`                 // time spent in indexing on primaries & replicas
	PriIndexingIndexTime         string `json:"pri.indexing.index_time"`             // time spent in indexing on primaries
	IndexingIndexTotal           int    `json:"indexing.index_total,string"`         // number of index ops on primaries & replicas
	PriIndexingIndexTotal        int    `json:"pri.indexing.index_total,string"`     // number of index ops on primaries
	IndexingIndexFailed          int    `json:"indexing.index_failed,string"`        // number of failed indexing ops on primaries & replicas
	PriIndexingIndexFailed       int    `json:"pri.indexing.index_failed,string"`    // number of failed indexing ops on primaries
	MergesCurrent                int    `json:"merges.current,string"`               // number of current merges on primaries & replicas
	PriMergesCurrent             int    `json:"pri.merges.current,string"`           // number of current merges on primaries
	MergesCurrentDocs            int    `json:"merges.current_docs,string"`          // number of current merging docs on primaries & replicas
	PriMergesCurrentDocs         int    `json:"pri.merges.current_docs,string"`      // number of current merging docs on primaries
	MergesCurrentSize            string `json:"merges.current_size"`                 // size of current merges on primaries & replicas
	PriMergesCurrentSize         string `json:"pri.merges.current_size"`             // size of current merges on primaries
	MergesTotal                  int    `json:"merges.total,string"`                 // number of completed merge ops on primaries & replicas
	PriMergesTotal               int    `json:"pri.merges.total,string"`             // number of completed merge ops on primaries
	MergesTotalDocs              int    `json:"merges.total_docs,string"`            // docs merged on primaries & replicas
	PriMergesTotalDocs           int    `json:"pri.merges.total_docs,string"`        // docs merged on primaries
	MergesTotalSize              string `json:"merges.total_size"`                   // size merged on primaries & replicas
	PriMergesTotalSize           string `json:"pri.merges.total_size"`               // size merged on primaries
	MergesTotalTime              string `json:"merges.total_time"`                   // time spent in merges on primaries & replicas
	PriMergesTotalTime           string `json:"pri.merges.total_time"`               // time spent in merges on primaries
	RefreshTotal                 int    `json:"refresh.total,string"`                // total refreshes on primaries & replicas
	PriRefreshTotal              int    `json:"pri.refresh.total,string"`            // total refreshes on primaries
	RefreshTime                  string `json:"refresh.time"`                        // time spent in refreshes on primaries & replicas
	PriRefreshTime               string `json:"pri.refresh.time"`                    // time spent in refreshes on primaries
	RefreshListeners             int    `json:"refresh.listeners,string"`            // number of pending refresh listeners on primaries & replicas
	PriRefreshListeners          int    `json:"pri.refresh.listeners,string"`        // number of pending refresh listeners on primaries
	SearchFetchCurrent           int    `json:"search.fetch_current,string"`         // current fetch phase ops on primaries & replicas
	PriSearchFetchCurrent        int    `json:"pri.search.fetch_current,string"`     // current fetch phase ops on primaries
	SearchFetchTime              string `json:"search.fetch_time"`                   // time spent in fetch phase on primaries & replicas
	PriSearchFetchTime           string `json:"pri.search.fetch_time"`               // time spent in fetch phase on primaries
	SearchFetchTotal             int    `json:"search.fetch_total,string"`           // total fetch ops on primaries & replicas
	PriSearchFetchTotal          int    `json:"pri.search.fetch_total,string"`       // total fetch ops on primaries
	SearchOpenContexts           int    `json:"search.open_contexts,string"`         // open search contexts on primaries & replicas
	PriSearchOpenContexts        int    `json:"pri.search.open_contexts,string"`     // open search contexts on primaries
	SearchQueryCurrent           int    `json:"search.query_current,string"`         // current query phase ops on primaries & replicas
	PriSearchQueryCurrent        int    `json:"pri.search.query_current,string"`     // current query phase ops on primaries
	SearchQueryTime              string `json:"search.query_time"`                   // time spent in query phase on primaries & replicas, e.g. "0s"
	PriSearchQueryTime           string `json:"pri.search.query_time"`               // time spent in query phase on primaries, e.g. "0s"
	SearchQueryTotal             int    `json:"search.query_total,string"`           // total query phase ops on primaries & replicas
	PriSearchQueryTotal          int    `json:"pri.search.query_total,string"`       // total query phase ops on primaries
	SearchScrollCurrent          int    `json:"search.scroll_current,string"`        // open scroll contexts on primaries & replicas
	PriSearchScrollCurrent       int    `json:"pri.search.scroll_current,string"`    // open scroll contexts on primaries
	SearchScrollTime             string `json:"search.scroll_time"`                  // time scroll contexts held open on primaries & replicas, e.g. "0s"
	PriSearchScrollTime          string `json:"pri.search.scroll_time"`              // time scroll contexts held open on primaries, e.g. "0s"
	SearchScrollTotal            int    `json:"search.scroll_total,string"`          // completed scroll contexts on primaries & replicas
	PriSearchScrollTotal         int    `json:"pri.search.scroll_total,string"`      // completed scroll contexts on primaries
	SegmentsCount                int    `json:"segments.count,string"`               // number of segments on primaries & replicas
	PriSegmentsCount             int    `json:"pri.segments.count,string"`           // number of segments on primaries
	SegmentsMemory               string `json:"segments.memory"`                     // memory used by segments on primaries & replicas, e.g. "1.3kb"
	PriSegmentsMemory            string `json:"pri.segments.memory"`                 // memory used by segments on primaries, e.g. "1.3kb"
	SegmentsIndexWriterMemory    string `json:"segments.index_writer_memory"`        // memory used by index writer on primaries & replicas, e.g. "0b"
	PriSegmentsIndexWriterMemory string `json:"pri.segments.index_writer_memory"`    // memory used by index writer on primaries, e.g. "0b"
	SegmentsVersionMapMemory     string `json:"segments.version_map_memory"`         // memory used by version map on primaries & replicas, e.g. "0b"
	PriSegmentsVersionMapMemory  string `json:"pri.segments.version_map_memory"`     // memory used by version map on primaries, e.g. "0b"
	SegmentsFixedBitsetMemory    string `json:"segments.fixed_bitset_memory"`        // memory used by fixed bit sets for nested object field types and type filters for types referred in _parent fields on primaries & replicas, e.g. "0b"
	PriSegmentsFixedBitsetMemory string `json:"pri.segments.fixed_bitset_memory"`    // memory used by fixed bit sets for nested object field types and type filters for types referred in _parent fields on primaries, e.g. "0b"
	WarmerCurrent                int    `json:"warmer.count,string"`                 // current warmer ops on primaries & replicas
	PriWarmerCurrent             int    `json:"pri.warmer.count,string"`             // current warmer ops on primaries
	WarmerTotal                  int    `json:"warmer.total,string"`                 // total warmer ops on primaries & replicas
	PriWarmerTotal               int    `json:"pri.warmer.total,string"`             // total warmer ops on primaries
	WarmerTotalTime              string `json:"warmer.total_time"`                   // time spent in warmers on primaries & replicas, e.g. "47s"
	PriWarmerTotalTime           string `json:"pri.warmer.total_time"`               // time spent in warmers on primaries, e.g. "47s"
	SuggestCurrent               int    `json:"suggest.current,string"`              // number of current suggest ops on primaries & replicas
	PriSuggestCurrent            int    `json:"pri.suggest.current,string"`          // number of current suggest ops on primaries
	SuggestTime                  string `json:"suggest.time"`                        // time spend in suggest on primaries & replicas, "31s"
	PriSuggestTime               string `json:"pri.suggest.time"`                    // time spend in suggest on primaries, e.g. "31s"
	SuggestTotal                 int    `json:"suggest.total,string"`                // number of suggest ops on primaries & replicas
	PriSuggestTotal              int    `json:"pri.suggest.total,string"`            // number of suggest ops on primaries
	MemoryTotal                  string `json:"memory.total"`                        // total user memory on primaries & replicas, e.g. "1.5kb"
	PriMemoryTotal               string `json:"pri.memory.total"`                    // total user memory on primaries, e.g. "1.5kb"
}
