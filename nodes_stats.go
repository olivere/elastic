// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"gopkg.in/olivere/elastic.v3/uritemplates"
)

type Level int64

const (
	NODE Level = iota
	INDICES
	SHARDS
)

var (
	_        = fmt.Print
	_        = log.Print
	_        = strings.Index
	_        = uritemplates.Expand
	_        = url.Parse
	levelMap = map[Level]string{NODE: "NODE", INDICES: "INDICES", SHARDS: "SHARDS"}
)

// NodesStatsService allows to retrieve one or more or all of the
// cluster nodes information.
// It is documented at https://www.elastic.co/guide/en/elasticsearch/reference/1.7/cluster-nodes-stats.html
type NodesStatsService struct {
	client           *Client
	pretty           bool
	nodeId           []string
	metric           []string
	indexMetric      []string
	completionFields []string
	fielddataFields  []string
	fields           []string
	groups           []string
	human            *bool
	level            *Level
	types            []string
	timeout          string
}

// NewNodesStatsService creates a new NodesStatsService.
func NewNodesStatsService(client *Client) *NodesStatsService {
	return &NodesStatsService{
		client:      client,
		nodeId:      []string{"_all"},
		metric:      []string{"_all"},
		indexMetric: []string{"_all"},
	}
}

// NodeId is a list of node IDs or names to limit the returned information.
// Use "_local" to return information from the node you're connecting to,
// leave empty to get information from all nodes.
func (s *NodesStatsService) NodeId(nodeId ...string) *NodesStatsService {
	s.nodeId = make([]string, 0)
	s.nodeId = append(s.nodeId, nodeId...)
	return s
}

// Metric is a list of metrics you wish returned. Leave empty to return all.
// Valid metrics are: breaker, fs, http, indices, jvm, network, os, process,
// thread_pool, and transport.
func (s *NodesStatsService) Metric(metric ...string) *NodesStatsService {
	s.metric = make([]string, 0)
	s.metric = append(s.metric, metric...)
	return s
}

// IndexMetric is a list of metrics you wish returned for the indices metric.
// Leave empty to return all.
// Valid metrics are: completion, docs, fielddata, filter_cache, flush, get,
// id_cache, indexing, merge, percolate, query_cache, request_cache, refresh,
// search, segments, store, warmer, and suggest.
func (s *NodesStatsService) IndexMetric(indexMetric ...string) *NodesStatsService {
	s.indexMetric = make([]string, 0)
	s.indexMetric = append(s.indexMetric, indexMetric...)
	return s
}

// CompletionFields is a comma-separated list of fields for `fielddata` and
// `suggest` index metric (supports wildcards).
func (s *NodesStatsService) CompletionFields(completionFields ...string) *NodesStatsService {
	s.completionFields = make([]string, 0)
	s.completionFields = append(s.completionFields, completionFields...)
	return s
}

// FielddataFields is a comma-separated list of fields for `fielddata` index
// metric (supports wildcards).
func (s *NodesStatsService) FielddataFields(fielddataFields ...string) *NodesStatsService {
	s.fielddataFields = make([]string, 0)
	s.fielddataFields = append(s.fielddataFields, fielddataFields...)
	return s
}

// Fields is a comma-separated list of fields for `fielddata` and `completion`
// index metric (supports wildcards).
func (s *NodesStatsService) Fields(fields ...string) *NodesStatsService {
	s.fields = make([]string, 0)
	s.fields = append(s.fields, fields...)
	return s
}

// Groups is a comma-separated list of search groups for `search` index metric.
func (s *NodesStatsService) Groups(groups ...string) *NodesStatsService {
	s.groups = make([]string, 0)
	s.groups = append(s.groups, groups...)
	return s
}

// Human indicates whether to return time and byte values in human-readable format (Default is FALSE).
func (s *NodesStatsService) Human(human bool) *NodesStatsService {
	s.human = &human
	return s
}

// Return indices stats aggregated at node, index or shard level.
func (s *NodesStatsService) Level(level Level) *NodesStatsService {
	s.level = &level
	return s
}

// Types is A comma-separated list of document types for the `indexing` index metric.
func (s *NodesStatsService) Types(types ...string) *NodesStatsService {
	s.types = make([]string, 0)
	s.types = append(s.types, types...)
	return s
}

// Timeout is documented as: Explicit operation timeout.
func (s *NodesStatsService) Timeout(timeout string) *NodesStatsService {
	s.timeout = timeout
	return s
}

// Pretty indicates whether to indent the returned JSON.
func (s *NodesStatsService) Pretty(pretty bool) *NodesStatsService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *NodesStatsService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_nodes/{node_id}/stats/{metric}/{index_metric}", map[string]string{
		"node_id":      strings.Join(s.nodeId, ","),
		"metric":       strings.Join(s.metric, ","),
		"index_metric": strings.Join(s.indexMetric, ","),
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if len(s.completionFields) > 0 {
		params.Set("completion_fields", strings.Join(s.completionFields, ","))
	}
	if len(s.fielddataFields) > 0 {
		params.Set("fielddata_fields", strings.Join(s.fielddataFields, ","))
	}
	if len(s.fields) > 0 {
		params.Set("fields", strings.Join(s.fields, ","))
	}
	if len(s.groups) > 0 {
		params.Set("groups", strings.Join(s.groups, ","))
	}
	if s.human != nil {
		params.Set("human", fmt.Sprintf("%v", *s.human))
	}
	if s.level != nil {
		params.Set("level", fmt.Sprintf("%v", levelMap[*s.level]))
	}
	if len(s.types) > 0 {
		params.Set("types", strings.Join(s.types, ","))
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.pretty {
		params.Set("pretty", "1")
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *NodesStatsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *NodesStatsService) Do() (*NodesStatsResponse, error) {
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
	res, err := s.client.PerformRequest("GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(NodesStatsResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// NodesStatsResponse is the response of NodesStatsService.Do.
type NodesStatsResponse struct {
	ClusterName string                 `json:"cluster_name"`
	Nodes       map[string]*NodesStats `json:"nodes"`
}

type NodesStats struct {
	Timestamp        int64                          `json:"timestamp"`
	Name             string                         `json:"name"`
	TransportAddress string                         `json:"transport_address"`
	Host             string                         `json:"host"`
	IP               []string                       `json:"ip"`
	Indices          *NodesStatsIndices             `json:"indices"`
	OS               *NodesStatsOS                  `json:"os"`
	Process          *NodesStatsProcess             `json:"process"`
	JVM              *NodesStatsJVM                 `json:"jvm"`
	ThreadPool       *NodesStatsThreadPool          `json:"thread_pool"`
	FS               *NodesStatsFS                  `json:"fs"`
	Transport        *NodesStatsTransport           `json:"transport"`
	HTTP             *NodesStatsHTTP                `json:"http"`
	Breakers         map[string]*NodesStatsBreakers `json:"breakers"`
	Script           *NodesStatsScript              `json:"breakers"`
}

type NodesStatsIndices struct {
	Docs         *NodesStatsIndicesDocs         `json:"docs"`
	Store        *NodesStatsIndicesStore        `json:"store"`
	Indexing     *NodesStatsIndicesIndexing     `json:"indexing"`
	Get          *NodesStatsIndicesGet          `json:"get"`
	Search       *NodesStatsIndicesSearch       `json:"search"`
	Merges       *NodesStatsIndicesMerges       `json:"merges"`
	Refresh      *NodesStatsIndicesRefresh      `json:"refresh"`
	Flush        *NodesStatsIndicesFlush        `json:"flush"`
	Warmer       *NodesStatsIndicesWarmer       `json:"warmer"`
	QueryCache   *NodesStatsIndicesQueryCache   `json:"query_cache"`
	FieldData    *NodesStatsIndicesFieldData    `json:"fielddata"`
	Percolate    *NodesStatsIndicesPercolate    `json:"percolate"`
	Completion   *NodesStatsIndicesCompletion   `json:"completion"`
	Segments     *NodesStatsIndicesSegments     `json:"segments"`
	Translog     *NodesStatsIndicesTranslog     `json:"translog"`
	Suggest      *NodesStatsIndicesSuggest      `json:"suggest"`
	RequestCache *NodesStatsIndicesRequestCache `json:"request_cache"`
	Recover      *NodesStatsIndicesRecover      `json:"recover"`
}

type NodesStatsOS struct {
	Timestamp   int64             `json:"timestamp"`
	CPUPercent  int64             `json:"cpu_percent"`
	LoadAverage float64           `json:"load_average"`
	Mem         *NodesStatsOSMem  `json:"mem"`
	Swap        *NodesStatsOSSwap `json:"swap"`
}

type NodesStatsProcess struct {
	Timestamp           int64                 `json:"timestamp"`
	OpenFileDescriptors int64                 `json:"open_file_descriptors"`
	MaxFileDescriptors  int64                 `json:"max_file_descriptors"`
	CPU                 *NodesStatsProcessCPU `json:"cpu"`
	Mem                 *NodesStatsProcessMem `json:"mem"`
}

type NodesStatsJVM struct {
	Timestamp   int64                                `json:"timestamp"`
	UptimeInMs  int64                                `json:"uptime_in_millis"`
	Mem         *NodesStatsJVMMem                    `json:"mem"`
	Threads     *NodesStatsJVMThreads                `json:"threads"`
	GC          *NodesStatsJVMGC                     `json:"gc"`
	BufferPools map[string]*NodesStatsJVMBufferPools `json:"buffer_pools"`
}

type NodesStatsThreadPool struct {
	Percolate  *NodesStatsThreadPoolSection `json:"percolate"`
	Listener   *NodesStatsThreadPoolSection `json:"listener"`
	Index      *NodesStatsThreadPoolSection `json:"index"`
	Refresh    *NodesStatsThreadPoolSection `json:"refresh"`
	Suggest    *NodesStatsThreadPoolSection `json:"suggest"`
	Generic    *NodesStatsThreadPoolSection `json:"generic"`
	Warmer     *NodesStatsThreadPoolSection `json:"warmer"`
	Search     *NodesStatsThreadPoolSection `json:"search"`
	Flush      *NodesStatsThreadPoolSection `json:"flush"`
	Optimize   *NodesStatsThreadPoolSection `json:"optimize"`
	Management *NodesStatsThreadPoolSection `json:"management"`
	Get        *NodesStatsThreadPoolSection `json:"get"`
	Merge      *NodesStatsThreadPoolSection `json:"merge"`
	Bulk       *NodesStatsThreadPoolSection `json:"bulk"`
	Snapshot   *NodesStatsThreadPoolSection `json:"snapshot"`
}

type NodesStatsFS struct {
	Timestamp int64               `json:"timestamp"`
	Total     *NodesStatsFSTotal  `json:"total"`
	Data      []*NodesStatsFSData `json:"data"`
}

type NodesStatsTransport struct {
	ServerOpen    int64 `json:"server_open"`
	RXCount       int64 `json:"rx_count"`
	RXSizeInBytes int64 `json:"rx_size_in_bytes"`
	TXCount       int64 `json:"tx_count"`
	TXSizeInBytes int64 `json:"tx_size_in_bytes"`
}

type NodesStatsHTTP struct {
	CurrentOpen int64 `json:"current_open"`
	TotalOpened int64 `json:"total_opened"`
}

type NodesStatsBreakers struct {
	LimitSizeInBytes     int64   `json:"limit_size_in_bytes"`
	LimitSize            string  `json:"limit_size"`
	EstimatedSizeInBytes int64   `json:"estimated_size_in_bytes"`
	EstimatedSize        string  `json:"estimated_size"`
	Overhead             float64 `json:"overhead"`
	Tripped              float64 `json:"tripped"`
}

type NodesStatsScript struct {
	Compilations   int64 `json:"timestamp"`
	CacheEvictions int64 `json:"cache_evictions"`
}

type NodesStatsIndicesDocs struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type NodesStatsIndicesStore struct {
	SizeInBytes      int64 `json:"size_in_bytes"`
	ThrottleTimeInMs int64 `json:"throttle_time_in_millis"`
}

type NodesStatsIndicesIndexing struct {
	IndexTotal       int64 `json:"index_total"`
	IndexTimeInMs    int64 `json:"index_time_in_millis"`
	IndexCurrent     int64 `json:"index_current"`
	IndexFailed      int64 `json:"index_failed"`
	DeleteTotal      int64 `json:"delete_total"`
	DeleteTotalInMs  int64 `json:"delete_time_in_millis"`
	DeleteCurrent    int64 `json:"delete_current"`
	NoopUpdateTotal  int64 `json:"noop_update_total"`
	IsThrottled      bool  `json:"is_throttled"`
	ThrottleTimeInMs int64 `json:"throttle_time_in_millis"`
}

type NodesStatsIndicesGet struct {
	Total           int64 `json:"total"`
	TimeInMs        int64 `json:"time_in_millis"`
	ExistsTotal     int64 `json:"exists_total"`
	ExistsTimeInMs  int64 `json:"exists_time_in_millis"`
	MissingTotal    int64 `json:"missing_total"`
	MissingTimeInMs int64 `json:"missing_time_in_millis"`
	Current         int64 `json:"current"`
}

type NodesStatsIndicesSearch struct {
	OpenContexts   int64 `json:"open_contexts"`
	QueryTotal     int64 `json:"query_total"`
	QueryTimeInMs  int64 `json:"query_time_in_millis"`
	QueryCurrent   int64 `json:"query_current"`
	FetchTotal     int64 `json:"fetch_total"`
	FetchTimeInMs  int64 `json:"fetch_time_in_millis"`
	FetchCurrent   int64 `json:"fetch_current"`
	ScrollTotal    int64 `json:"scroll_total"`
	ScrollTimeInMs int64 `json:"scroll_time_in_millis"`
	ScrollCurrent  int64 `json:"scroll_current"`
}

type NodesStatsIndicesMerges struct {
	Current                  int64 `json:"current"`
	CurrentDocs              int64 `json:"current_docs"`
	CurrentSizeInBytes       int64 `json:"current_size_in_bytes"`
	Total                    int64 `json:"total"`
	TotalTimeInMs            int64 `json:"total_time_in_millis"`
	TotalDocs                int64 `json:"total_docs"`
	TotalSizeInBytes         int64 `json:"total_size_in_bytes"`
	TotalStoppedTimeInMs     int64 `json:"total_stopped_time_in_millis"`
	TotalThrottledTimeInMs   int64 `json:"total_throttled_time_in_millis"`
	TotalAutoThrottleInBytes int64 `json:"total_auto_throttle_in_bytes"`
}

type NodesStatsIndicesRefresh struct {
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}

type NodesStatsIndicesFlush struct {
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}

type NodesStatsIndicesWarmer struct {
	Current       int64 `json:"current"`
	Total         int64 `json:"total"`
	TotalTimeInMs int64 `json:"total_time_in_millis"`
}

type NodesStatsIndicesQueryCache struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	TotalCount        int64 `json:"total_count"`
	HitCount          int64 `json:"hit_count"`
	MissCount         int64 `json:"miss_count"`
	CacheSize         int64 `json:"cache_size"`
	CacheCount        int64 `json:"cache_count"`
	Evictions         int64 `json:"evictions"`
}

type NodesStatsIndicesFieldData struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	Evictions         int64 `json:"evictions"`
}

type NodesStatsIndicesPercolate struct {
	Total             int64  `json:"total"`
	TimeInMs          int64  `json:"time_in_millis"`
	Current           int64  `json:"current"`
	MemorySizeInBytes int64  `json:"memory_size_in_bytes"`
	MemorySize        string `json:"memory_size"`
	Queries           int64  `json:"queries"`
}

type NodesStatsIndicesCompletion struct {
	SizeInBytes int64 `json:"size_in_bytes"`
}

type NodesStatsIndicesSegments struct {
	Count                       int64 `json:"count"`
	MemoryInBytes               int64 `json:"memory_in_bytes"`
	TermsMemoryInBytes          int64 `json:"terms_memory_in_bytes"`
	StoredFieldsMemoryInBytes   int64 `json:"stored_fields_mem"`
	TermVectorsMemoryInBytes    int64 `json:"term_vectors_memory_in_bytes"`
	NormsMemoryInBytes          int64 `json:"norms_memory_in_bytes"`
	DocValuesMemoryInBytes      int64 `json:"doc_values_memory_in_bytes"`
	IndexWriterMemoryInBytes    int64 `json:"index_writer_memory_in_bytes"`
	IndexWriterMaxMemoryInBytes int64 `json:"index_writer_max_memory_in_bytes"`
	VersionMapMemoryInBytes     int64 `json:"version_map_memory_in_bytes"`
	FixedBitSetMemoryInBytes    int64 `json:"fixed_bit_set_memory_in_bytes"`
}

type NodesStatsIndicesTranslog struct {
	Operations  int64 `json:"operations"`
	SizeInBytes int64 `json:"size_in_bytes"`
}

type NodesStatsIndicesSuggest struct {
	Total    int64 `json:"total"`
	TimeInMs int64 `json:"time_in_millis"`
	Current  int64 `json:"current"`
}

type NodesStatsIndicesRequestCache struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	Evictions         int64 `json:"evictions"`
	HitCount          int64 `json:"hit_count"`
	MissCount         int64 `json:"miss_count"`
}

type NodesStatsIndicesRecover struct {
	CurrentAsSource  int64 `json:"current_as_source"`
	CurrentAsTarget  int64 `json:"current_as_target"`
	ThrottleTimeInMs int64 `json:"throttle_time_in_millis"`
}

type NodesStatsOSMem struct {
	TotalInBytes int64 `json:"total_in_bytes"`
	FreeInBytes  int64 `json:"free_in_bytes"`
	UsedInBytes  int64 `json:"used_in_bytes"`
	FreePercent  int64 `json:"free_percent"`
	UsedPercent  int64 `json:"used_percent"`
}

type NodesStatsOSSwap struct {
	TotalInBytes int64 `json:"total_in_bytes"`
	FreeInBytes  int64 `json:"free_in_bytes"`
	UsedInBytes  int64 `json:"used_in_bytes"`
}

type NodesStatsProcessCPU struct {
	Percent   int64 `json:"percent"`
	TotalInMs int64 `json:"total_in_millis"`
}

type NodesStatsProcessMem struct {
	TotalVirtualInBytes int64 `json:"total_virtual_in_bytes"`
}

type NodesStatsJVMMem struct {
	HeapUsedInBytes         int64                             `json:"heap_used_in_bytes"`
	HeapUsedPercent         int64                             `json:"heap_used_percent"`
	HeapCommittedInBytes    int64                             `json:"heap_committed_in_bytes"`
	HeapMaxInBytes          int64                             `json:"heap_max_in_bytes"`
	NonHeapUsedInBytes      int64                             `json:"non_heap_used_in_bytes"`
	NonHeapCommittedInBytes int64                             `json:"non_heap_committed_in_bytes"`
	Pools                   map[string]*NodesStatsJVMMemPools `json:"pools"`
}

type NodesStatsJVMThreads struct {
	Count     int64 `json:"count"`
	PeakCount int64 `json:"peak_count"`
}

type NodesStatsJVMGC struct {
	Collectors map[string]*NodesStatsJVMGCCollectors `json:"collectors"`
}

type NodesStatsJVMBufferPools struct {
	Count                int64 `json:"count"`
	UsedInBytes          int64 `json:"used_in_bytes"`
	TotalCapacityInBytes int64 `total_capacity_in_bytes:"pid"`
}

type NodesStatsJVMMemPools struct {
	UsedInBytes     int64 `json:"used_in_bytes"`
	MaxInBytes      int64 `json:"max_in_bytes"`
	PeakUsedInBytes int64 `json:"peak_used_in_bytes"`
	PeakMaxInBytes  int64 `json:"peak_max_in_bytes"`
}

type NodesStatsJVMGCCollectors struct {
	CollectionCount    int64 `json:"collection_count"`
	CollectionTimeInMs int64 `json:"collection_time_in_millis"`
}

type NodesStatsThreadPoolSection struct {
	Threads   int64 `json:"threads"`
	Queue     int64 `json:"queue"`
	Active    int64 `json:"active"`
	Rejected  int64 `json:"rejected"`
	Largest   int64 `json:"largest"`
	Completed int64 `json:"completed"`
}

type NodesStatsFSTotal struct {
	TotalInBytes     int64  `json:"total_in_bytes"`
	FreeInBytes      int64  `json:"free_in_bytes"`
	AvailableInBytes int64  `json:"available_in_bytes"`
	Spins            string `json:"spins"`
}

type NodesStatsFSData struct {
	Path             string `json:"path"`
	Mount            string `json:"mount"`
	Type             string `json:"type"`
	Dev              string `json:"dev"`
	TotalInBytes     int64  `json:"total_in_bytes"`
	FreeInBytes      int64  `json:"free_in_bytes"`
	AvailableInBytes int64  `json:"available_in_bytes"`
	Spins            string `json:"spins"`
}
