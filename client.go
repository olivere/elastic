// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// Version is the current version of Elastic.
	Version = "2.0.26"

	// DefaultUrl is the default endpoint of Elasticsearch on the local machine.
	DefaultURL = "http://127.0.0.1:9200"

	// DefaultSendGetBodyAs is the HTTP method to use when elastic is sending
	// a GET request with a body.
	DefaultSendGetBodyAs = "GET"

	// DefaultGzipEnabled specifies if gzip compression is enabled by default.
	DefaultGzipEnabled = false
)

var (
	// ErrNoClient is raised when no Elasticsearch node is available.
	ErrNoClient = errors.New("no Elasticsearch node available")

	// ErrRetry is raised when a request cannot be executed after the configured
	// number of retries.
	ErrRetry = errors.New("cannot connect after several retries")

	// ErrTimeout is raised when a request timed out, e.g. when WaitForStatus
	// didn't return in time.
	ErrTimeout = errors.New("timeout")

	// ErrNoURLSpecified is raised when no Elasticsearch URL has been specified.
	ErrNoURLSpecified = errors.New("no Elasticsearch URL specified")

	// ErrURLAlreadySpecified is raised when no Elasticsearch URL has been specified.
	ErrURLAlreadySpecified = errors.New("Elasticsearch URL specified more than once")

	// ErrInvalidURLSpecified is raised when no Elasticsearch URL has been specified.
	ErrInvalidURLSpecified = errors.New("invalid Elasticsearch URL specified")
)

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*Client) error

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	c *http.Client // net/http Client to use for requests

	mu                sync.RWMutex // guards the next block
	url               string
	errorlog          *log.Logger // error log for critical messages
	infolog           *log.Logger // information log for e.g. response times
	tracelog          *log.Logger // trace log for debugging
	scheme            string      // http or https
	decoder           Decoder     // used to decode data sent from Elasticsearch
	basicAuth         bool        // indicates whether to send HTTP Basic Auth credentials
	basicAuthUsername string      // username for HTTP Basic Auth
	basicAuthPassword string      // password for HTTP Basic Auth
	sendGetBodyAs     string      // override for when sending a GET with a body
	gzipEnabled       bool        // gzip compression enabled or disabled (default)
}

// NewClient creates a new client to work with Elasticsearch.
//
// The caller can configure the new client by passing configuration options
// to the func.
//
// Example:
//
//   client, err := elastic.NewClient(
//     elastic.SetURL("http://localhost:9200", "http://localhost:9201"),
//     elastic.SetMaxRetries(10),
//     elastic.SetBasicAuth("user", "secret"))
//
// If no URL is configured, Elastic uses DefaultURL by default.
//
// If the sniffer is enabled (the default), the new client then sniffes
// the cluster via the Nodes Info API
// (see http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/cluster-nodes-info.html#cluster-nodes-info).
// It uses the URLs specified by the caller. The caller is responsible
// to only pass a list of URLs of nodes that belong to the same cluster.
// This sniffing process is run on startup and periodically.
// Use SnifferInterval to set the interval between two sniffs (default is
// 15 minutes). In other words: By default, the client will find new nodes
// in the cluster and remove those that are no longer available every
// 15 minutes. Disable the sniffer by passing SetSniff(false) to NewClient.
//
// The list of nodes found in the sniffing process will be used to make
// connections to the REST API of Elasticsearch. These nodes are also
// periodically checked in a shorter time frame. This process is called
// a health check. By default, a health check is done every 60 seconds.
// You can set a shorter or longer interval by SetHealthcheckInterval.
// Disabling health checks is not recommended, but can be done by
// SetHealthcheck(false).
//
// Connections are automatically marked as dead or healthy while
// making requests to Elasticsearch. When a request fails, Elastic will
// retry up to a maximum number of retries configured with SetMaxRetries.
// Retries are disabled by default.
//
// If no HttpClient is configured, then http.DefaultClient is used.
// You can use your own http.Client with some http.Transport for
// advanced scenarios.
//
// An error is also returned when some configuration option is invalid or
// the new client cannot sniff the cluster (if enabled).
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		c:             http.DefaultClient,
		decoder:       &DefaultDecoder{},
		sendGetBodyAs: DefaultSendGetBodyAs,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(c.url)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, ErrInvalidURLSpecified
	}
	u.Fragment = ""
	u.Path = ""
	u.RawQuery = ""
	c.url = u.String()

	return c, nil
}

// SetHttpClient can be used to specify the http.Client to use when making
// HTTP requests to Elasticsearch.
func SetHttpClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		if httpClient != nil {
			c.c = httpClient
		} else {
			c.c = http.DefaultClient
		}
		return nil
	}
}

// SetBasicAuth can be used to specify the HTTP Basic Auth credentials to
// use when making HTTP requests to Elasticsearch.
func SetBasicAuth(username, password string) ClientOptionFunc {
	return func(c *Client) error {
		c.basicAuthUsername = username
		c.basicAuthPassword = password
		c.basicAuth = c.basicAuthUsername != "" || c.basicAuthPassword != ""
		return nil
	}
}

// SetURL defines the URL endpoints of the Elasticsearch nodes. Notice that
// when sniffing is enabled, these URLs are used to initially sniff the
// cluster on startup.
func SetURL(url string) ClientOptionFunc {
	return func(c *Client) error {
		if len(c.url) != 0 {
			return ErrURLAlreadySpecified
		}

		c.url = url
		return nil
	}
}

// SetGzip enables or disables gzip compression (disabled by default).
func SetGzip(enabled bool) ClientOptionFunc {
	return func(c *Client) error {
		c.gzipEnabled = enabled
		return nil
	}
}

// SetDecoder sets the Decoder to use when decoding data from Elasticsearch.
// DefaultDecoder is used by default.
func SetDecoder(decoder Decoder) func(*Client) error {
	return func(c *Client) error {
		if decoder != nil {
			c.decoder = decoder
		} else {
			c.decoder = &DefaultDecoder{}
		}
		return nil
	}
}

// SetErrorLog sets the logger for critical messages like nodes joining
// or leaving the cluster or failing requests. It is nil by default.
func SetErrorLog(logger *log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.errorlog = logger
		return nil
	}
}

// SetInfoLog sets the logger for informational messages, e.g. requests
// and their response times. It is nil by default.
func SetInfoLog(logger *log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.infolog = logger
		return nil
	}
}

// SetTraceLog specifies the log.Logger to use for output of HTTP requests
// and responses which is helpful during debugging. It is nil by default.
func SetTraceLog(logger *log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.tracelog = logger
		return nil
	}
}

// SendGetBodyAs specifies the HTTP method to use when sending a GET request
// with a body. It is GET by default.
func SetSendGetBodyAs(httpMethod string) func(*Client) error {
	return func(c *Client) error {
		c.sendGetBodyAs = httpMethod
		return nil
	}
}

// String returns a string representation of the client status.
func (c *Client) String() string {
	return fmt.Sprintf("ElasticsearchClient(%s)", c.url)
}

// errorf logs to the error log.
func (c *Client) errorf(format string, args ...interface{}) {
	if c.errorlog != nil {
		c.errorlog.Printf(format, args...)
	}
}

// infof logs informational messages.
func (c *Client) infof(format string, args ...interface{}) {
	if c.infolog != nil {
		c.infolog.Printf(format, args...)
	}
}

// tracef logs to the trace log.
func (c *Client) tracef(format string, args ...interface{}) {
	if c.tracelog != nil {
		c.tracelog.Printf(format, args...)
	}
}

// dumpRequest dumps the given HTTP request to the trace log.
func (c *Client) dumpRequest(r *http.Request) {
	if c.tracelog != nil {
		out, err := httputil.DumpRequestOut(r, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// dumpResponse dumps the given HTTP response to the trace log.
func (c *Client) dumpResponse(resp *http.Response) {
	if c.tracelog != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// PerformRequest does a HTTP request to Elasticsearch.
// It returns a response and an error on failure.
func (c *Client) PerformRequest(method, path string, params url.Values, body interface{}) (*Response, error) {
	start := time.Now().UTC()

	c.mu.RLock()
	basicAuth := c.basicAuth
	basicAuthUsername := c.basicAuthUsername
	basicAuthPassword := c.basicAuthPassword
	sendGetBodyAs := c.sendGetBodyAs
	gzipEnabled := c.gzipEnabled
	c.mu.RUnlock()

	var err error
	var req *Request
	var resp *Response

	// Change method if sendGetBodyAs is specified.
	if method == "GET" && body != nil && sendGetBodyAs != "GET" {
		method = sendGetBodyAs
	}

	pathWithParams := path
	if len(params) > 0 {
		pathWithParams += "?" + params.Encode()
	}

	req, err = NewRequest(method, c.url+pathWithParams)
	if err != nil {
		c.errorf("elastic: cannot create request for %s %s: %v", strings.ToUpper(method), c.url+pathWithParams, err)
		return nil, err
	}

	if basicAuth {
		req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
	}

	// Set body
	if body != nil {
		req.SetBody(body, gzipEnabled)
	}

	// Tracing
	c.dumpRequest((*http.Request)(req))

	// Get response
	res, err := c.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Check for errors
	if err := checkResponse(res); err != nil {
		return nil, err
	}

	//TODO: move this before response check for consistent logs?
	// Tracing
	c.dumpResponse(res)

	resp, err = c.newResponse(res)
	if err != nil {
		return nil, err
	}

	duration := time.Now().UTC().Sub(start)
	c.infof("%s %s [status:%d, request:%.3fs]",
		strings.ToUpper(method),
		req.URL,
		resp.StatusCode,
		float64(int64(duration/time.Millisecond))/1000)

	return resp, nil
}

// ElasticsearchVersion returns the version number of Elasticsearch
// running on the given URL.
func (c *Client) ElasticsearchVersion(url string) (string, error) {
	res, _, err := c.Ping().URL(url).Do()
	if err != nil {
		return "", err
	}
	return res.Version.Number, nil
}

// IndexNames returns the names of all indices in the cluster.
func (c *Client) IndexNames() ([]string, error) {
	res, err := c.IndexGetSettings().Index("_all").Do()
	if err != nil {
		return nil, err
	}
	var names []string
	for name, _ := range res {
		names = append(names, name)
	}
	return names, nil
}

// Ping checks if a given node in a cluster exists and (optionally)
// returns some basic information about the Elasticsearch server,
// e.g. the Elasticsearch version number.
func (c *Client) Ping() *PingService {
	return NewPingService(c)
}

// CreateIndex returns a service to create a new index.
func (c *Client) CreateIndex(name string) *CreateIndexService {
	builder := NewCreateIndexService(c)
	builder.Index(name)
	return builder
}

// DeleteIndex returns a service to delete an index.
func (c *Client) DeleteIndex(name string) *DeleteIndexService {
	builder := NewDeleteIndexService(c)
	builder.Index(name)
	return builder
}

// IndexExists allows to check if an index exists.
func (c *Client) IndexExists(name string) *IndexExistsService {
	builder := NewIndexExistsService(c)
	builder.Index(name)
	return builder
}

// TypeExists allows to check if one or more types exist in one or more indices.
func (c *Client) TypeExists() *IndicesExistsTypeService {
	return NewIndicesExistsTypeService(c)
}

// IndexStats provides statistics on different operations happining
// in one or more indices.
func (c *Client) IndexStats(indices ...string) *IndicesStatsService {
	builder := NewIndicesStatsService(c)
	builder = builder.Index(indices...)
	return builder
}

// OpenIndex opens an index.
func (c *Client) OpenIndex(name string) *OpenIndexService {
	builder := NewOpenIndexService(c)
	builder.Index(name)
	return builder
}

// CloseIndex closes an index.
func (c *Client) CloseIndex(name string) *CloseIndexService {
	builder := NewCloseIndexService(c)
	builder.Index(name)
	return builder
}

// Index a document.
func (c *Client) Index() *IndexService {
	builder := NewIndexService(c)
	return builder
}

// IndexGet retrieves information about one or more indices.
// IndexGet is only available for Elasticsearch 1.4 or later.
func (c *Client) IndexGet() *IndicesGetService {
	builder := NewIndicesGetService(c)
	return builder
}

// IndexGetSettings retrieves settings about one or more indices.
func (c *Client) IndexGetSettings() *IndicesGetSettingsService {
	builder := NewIndicesGetSettingsService(c)
	return builder
}

// Update a document.
func (c *Client) Update() *UpdateService {
	builder := NewUpdateService(c)
	return builder
}

// Delete a document.
func (c *Client) Delete() *DeleteService {
	builder := NewDeleteService(c)
	return builder
}

// DeleteByQuery deletes documents as found by a query.
func (c *Client) DeleteByQuery() *DeleteByQueryService {
	builder := NewDeleteByQueryService(c)
	return builder
}

// Get a document.
func (c *Client) Get() *GetService {
	builder := NewGetService(c)
	return builder
}

// MultiGet retrieves multiple documents in one roundtrip.
func (c *Client) MultiGet() *MultiGetService {
	builder := NewMultiGetService(c)
	return builder
}

// Exists checks if a document exists.
func (c *Client) Exists() *ExistsService {
	builder := NewExistsService(c)
	return builder
}

// Count documents.
func (c *Client) Count(indices ...string) *CountService {
	builder := NewCountService(c)
	builder.Indices(indices...)
	return builder
}

// Search is the entry point for searches.
func (c *Client) Search(indices ...string) *SearchService {
	builder := NewSearchService(c)
	builder.Indices(indices...)
	return builder
}

// Percolate allows to send a document and return matching queries.
// See http://www.elastic.co/guide/en/elasticsearch/reference/current/search-percolate.html.
func (c *Client) Percolate() *PercolateService {
	builder := NewPercolateService(c)
	return builder
}

// MultiSearch is the entry point for multi searches.
func (c *Client) MultiSearch() *MultiSearchService {
	return NewMultiSearchService(c)
}

// Suggest returns a service to return suggestions.
func (c *Client) Suggest(indices ...string) *SuggestService {
	builder := NewSuggestService(c)
	builder.Indices(indices...)
	return builder
}

// Scan through documents. Use this to iterate inside a server process
// where the results will be processed without returning them to a client.
func (c *Client) Scan(indices ...string) *ScanService {
	builder := NewScanService(c)
	builder.Indices(indices...)
	return builder
}

// Scroll through documents. Use this to efficiently scroll through results
// while returning the results to a client. Use Scan when you don't need
// to return requests to a client (i.e. not paginating via request/response).
func (c *Client) Scroll(indices ...string) *ScrollService {
	builder := NewScrollService(c)
	builder.Indices(indices...)
	return builder
}

// ClearScroll can be used to clear search contexts manually.
func (c *Client) ClearScroll() *ClearScrollService {
	builder := NewClearScrollService(c)
	return builder
}

// Optimize asks Elasticsearch to optimize one or more indices.
func (c *Client) Optimize(indices ...string) *OptimizeService {
	builder := NewOptimizeService(c)
	builder.Indices(indices...)
	return builder
}

// Refresh asks Elasticsearch to refresh one or more indices.
func (c *Client) Refresh(indices ...string) *RefreshService {
	builder := NewRefreshService(c)
	builder.Indices(indices...)
	return builder
}

// Flush asks Elasticsearch to free memory from the index and
// flush data to disk.
func (c *Client) Flush() *FlushService {
	builder := NewFlushService(c)
	return builder
}

// Explain computes a score explanation for a query and a specific document.
func (c *Client) Explain(index, typ, id string) *ExplainService {
	builder := NewExplainService(c)
	builder = builder.Index(index).Type(typ).Id(id)
	return builder
}

// Bulk is the entry point to mass insert/update/delete documents.
func (c *Client) Bulk() *BulkService {
	builder := NewBulkService(c)
	return builder
}

// Alias enables the caller to add and/or remove aliases.
func (c *Client) Alias() *AliasService {
	builder := NewAliasService(c)
	return builder
}

// Aliases returns aliases by index name(s).
func (c *Client) Aliases() *AliasesService {
	builder := NewAliasesService(c)
	return builder
}

// GetTemplate gets a search template.
// Use IndexXXXTemplate funcs to manage index templates.
func (c *Client) GetTemplate() *GetTemplateService {
	return NewGetTemplateService(c)
}

// PutTemplate creates or updates a search template.
// Use IndexXXXTemplate funcs to manage index templates.
func (c *Client) PutTemplate() *PutTemplateService {
	return NewPutTemplateService(c)
}

// DeleteTemplate deletes a search template.
// Use IndexXXXTemplate funcs to manage index templates.
func (c *Client) DeleteTemplate() *DeleteTemplateService {
	return NewDeleteTemplateService(c)
}

// IndexGetTemplate gets an index template.
// Use XXXTemplate funcs to manage search templates.
func (c *Client) IndexGetTemplate(names ...string) *IndicesGetTemplateService {
	builder := NewIndicesGetTemplateService(c)
	builder = builder.Name(names...)
	return builder
}

// IndexTemplateExists gets check if an index template exists.
// Use XXXTemplate funcs to manage search templates.
func (c *Client) IndexTemplateExists(name string) *IndicesExistsTemplateService {
	builder := NewIndicesExistsTemplateService(c)
	builder = builder.Name(name)
	return builder
}

// IndexPutTemplate creates or updates an index template.
// Use XXXTemplate funcs to manage search templates.
func (c *Client) IndexPutTemplate(name string) *IndicesPutTemplateService {
	builder := NewIndicesPutTemplateService(c)
	builder = builder.Name(name)
	return builder
}

// IndexDeleteTemplate deletes an index template.
// Use XXXTemplate funcs to manage search templates.
func (c *Client) IndexDeleteTemplate(name string) *IndicesDeleteTemplateService {
	builder := NewIndicesDeleteTemplateService(c)
	builder = builder.Name(name)
	return builder
}

// GetMapping gets a mapping.
func (c *Client) GetMapping() *GetMappingService {
	return NewGetMappingService(c)
}

// PutMapping registers a mapping.
func (c *Client) PutMapping() *PutMappingService {
	return NewPutMappingService(c)
}

// DeleteMapping deletes a mapping.
func (c *Client) DeleteMapping() *DeleteMappingService {
	return NewDeleteMappingService(c)
}

// GetWarmer gets one or more warmers by name.
func (c *Client) GetWarmer() *IndicesGetWarmerService {
	return NewIndicesGetWarmerService(c)
}

// PutWarmer registers a warmer.
func (c *Client) PutWarmer() *IndicesPutWarmerService {
	return NewIndicesPutWarmerService(c)
}

// DeleteWarmer deletes one or more warmers.
func (c *Client) DeleteWarmer() *IndicesDeleteWarmerService {
	return NewIndicesDeleteWarmerService(c)
}

// ClusterHealth retrieves the health of the cluster.
func (c *Client) ClusterHealth() *ClusterHealthService {
	return NewClusterHealthService(c)
}

// ClusterState retrieves the state of the cluster.
func (c *Client) ClusterState() *ClusterStateService {
	return NewClusterStateService(c)
}

// ClusterStats retrieves cluster statistics.
func (c *Client) ClusterStats() *ClusterStatsService {
	return NewClusterStatsService(c)
}

// NodesInfo retrieves one or more or all of the cluster nodes information.
func (c *Client) NodesInfo() *NodesInfoService {
	return NewNodesInfoService(c)
}

// Reindex returns a service that will reindex documents from a source
// index into a target index. See
// http://www.elastic.co/guide/en/elasticsearch/guide/current/reindex.html
// for more information about reindexing.
func (c *Client) Reindex(sourceIndex, targetIndex string) *Reindexer {
	return NewReindexer(c, sourceIndex, CopyToTargetIndex(targetIndex))
}

// WaitForStatus waits for the cluster to have the given status.
// This is a shortcut method for the ClusterHealth service.
//
// WaitForStatus waits for the specified timeout, e.g. "10s".
// If the cluster will have the given state within the timeout, nil is returned.
// If the request timed out, ErrTimeout is returned.
func (c *Client) WaitForStatus(status string, timeout string) error {
	health, err := c.ClusterHealth().WaitForStatus(status).Timeout(timeout).Do()
	if err != nil {
		return err
	}
	if health.TimedOut {
		return ErrTimeout
	}
	return nil
}

// WaitForGreenStatus waits for the cluster to have the "green" status.
// See WaitForStatus for more details.
func (c *Client) WaitForGreenStatus(timeout string) error {
	return c.WaitForStatus("green", timeout)
}

// WaitForYellowStatus waits for the cluster to have the "yellow" status.
// See WaitForStatus for more details.
func (c *Client) WaitForYellowStatus(timeout string) error {
	return c.WaitForStatus("yellow", timeout)
}

// TermVector returns information and statistics on terms in the fields
// of a particular document.
func (c *Client) TermVector(index, typ string) *TermvectorService {
	builder := NewTermvectorService(c)
	builder = builder.Index(index).Type(typ)
	return builder
}
