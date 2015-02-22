// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	// Version is the current version of Elastic.
	Version = "1.5.0.dev"

	// DefaultUrl is the default endpoint of Elasticsearch on the local machine.
	// It is used e.g. when initializing a new Client without a specific URL.
	DefaultURL = "http://127.0.0.1:9200"

	// DefaultScheme is the default protocol scheme to use when sniffing
	// the Elasticsearch cluster.
	DefaultScheme = "http"

	// DefaultHealthcheckSchedule is the default duration between
	// two health checks of the nodes in the cluster.
	DefaultHealthcheckSchedule = 60 * time.Second

	// DefaultSnifferSchedule is the duration between two sniffing procedures,
	// i.e. the lookup of all nodes in the cluster and their addition/removal
	// from the list of actual connections.
	DefaultSnifferSchedule = 15 * time.Minute

	// DefaultSnifferTimeout is the default timeout after which the
	// sniffing process times out.
	DefaultSnifferTimeout = 1 * time.Second

	// DefaultMaxRetries is the number of retries for a single request after
	// Elastic will give up and return an error. It is zero by default, so
	// retry is disabled by default.
	DefaultMaxRetries = 0
)

var (
	// ErrNoClient is raised when no Elasticsearch node is available.
	ErrNoClient = errors.New("no Elasticsearch node available")

	// ErrRetry is raised when a request cannot be executed after the configured
	// number of retries.
	ErrRetry = errors.New("cannot connect after several retries")
)

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	c *http.Client // net/http Client to use for requests

	connsMu sync.RWMutex // connsMu guards the next block
	conns   []*conn      // all connections
	cindex  int          // index into conns

	configMu                  sync.RWMutex  // guards the next block
	running                   bool          // true if the client's background processes are running
	logger                    *log.Logger   // standard log
	tracer                    *log.Logger   // trace log
	maxRetries                int           // max. number of retries
	scheme                    string        // http or https
	healthcheckSchedule       time.Duration // schedule for healthcheck of all nodes
	healthcheckScheduleUpdate chan bool     // notify healthchecker about updated schedule
	healthcheckStop           chan bool     // notify healthchecker to stop
	healthcheckStopped        chan bool     // notification from healthchecker that it is stopped
	snifferTimeout            time.Duration // time the sniffer waits for a response from nodes info API
	snifferSchedule           time.Duration // schedule for sniffing process
	snifferScheduleUpdate     chan bool     // notify sniffer about updated schedule
	snifferStop               chan bool     // notify sniffer to stop
	snifferStopped            chan bool     // notification from sniffer that it is stopped
	decoder                   Decoder       // used to decode data sent from Elasticsearch

	mu        sync.RWMutex // mutex for the next two fields
	activeUrl string       // currently active connection url
	hasActive bool         // true if we have an active connection
}

// NewClient creates a new client to work with Elasticsearch.
//
// The caller can specify zero or more URLs of nodes in a cluster. If the
// caller no URL, the default URL of http://127.0.0.1:9200 is used.
//
// The new client then sniffes the cluster via the Nodes Info API
// (see http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/cluster-nodes-info.html#cluster-nodes-info).
// It uses the URLs specified by the caller. The caller is responsible
// to only pass a list of URLs of nodes that belong to the same cluster.
// This sniffing process is run periodically. Use SetSnifferSchedule to
// set its schedule (default is 15 minutes). In other words: By default,
// the client will find new nodes in the cluster and remove those that are
// no longer available every 15 minutes.
//
// The list of nodes found in the sniffing process will be used to make
// connections to the REST API of Elasticsearch. These nodes are also
// periodically checked in a shorter time frame. This process is called
// a health check. By default, a health check is done every 60 seconds.
// You can set a shorter or longer schedule by SetHealthcheckSchedule.
//
// Connections are automatically marked as dead or healthy while
// making requests to Elasticsearch. By default, retries are disabled.
// If you want to enable retries, set a maximum number of retries with
// SetMaxRetries.
//
// An error is returned when no client is passed as a first argument.
// You can pass http.DefaultClient if you don't use your own client or
// transport.
//
// An error is also returned when the new client cannot sniff the cluster,
// i.e. it does not find any nodes.
func NewClient(client *http.Client, urls ...string) (*Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Set up the client and initialize the URLs
	c := &Client{
		conns:                     make([]*conn, 0),
		cindex:                    -1,
		c:                         client,
		scheme:                    DefaultScheme,
		decoder:                   &DefaultDecoder{},
		maxRetries:                DefaultMaxRetries,
		healthcheckSchedule:       DefaultHealthcheckSchedule,
		healthcheckScheduleUpdate: make(chan bool),
		healthcheckStop:           make(chan bool),
		healthcheckStopped:        make(chan bool),
		snifferSchedule:           DefaultSnifferSchedule,
		snifferScheduleUpdate:     make(chan bool),
		snifferStop:               make(chan bool),
		snifferStopped:            make(chan bool),
		snifferTimeout:            DefaultSnifferTimeout,
	}

	if len(urls) == 0 {
		urls = []string{DefaultURL}
	}
	urls = canonicalize(urls...)

	// Sniff the cluster initially
	if err := c.sniff(urls...); err != nil {
		return nil, err
	}

	// Perform an initial health check
	c.healthcheck()

	go c.sniffer()       // periodically update cluster information
	go c.healthchecker() // start goroutine periodically ping all nodes of the cluster

	c.configMu.Lock()
	c.running = true
	c.configMu.Unlock()

	c.logf("Client started on %s", c.String())

	return c, nil
}

// String returns a string representation of the client status.
func (c *Client) String() string {
	c.connsMu.Lock()
	conns := c.conns
	c.connsMu.Unlock()

	var buf bytes.Buffer
	for i, conn := range conns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(conn.String())
	}
	return buf.String()
}

// IsRunning returns true if the background processes of the client are
// running, false otherwise.
func (c *Client) IsRunning() bool {
	c.configMu.RLock()
	defer c.configMu.RUnlock()
	return c.running
}

// Start starts the background processes like sniffing the cluster and
// periodic health checks. You don't need to run Start when creating a
// client with NewClient; the background processes are run by default.
//
// If the background processes are already running, then this is a no-op.
func (c *Client) Start() {
	c.configMu.RLock()
	if c.running {
		c.configMu.RUnlock()
		return
	}
	c.configMu.RUnlock()

	go c.sniffer()
	go c.healthchecker()

	c.configMu.Lock()
	c.running = true
	c.configMu.Unlock()

	c.logf("Client started")
}

// Stop stops the background processes that the client is running,
// i.e. sniffing the cluster periodically and running health checks
// on the nodes.
//
// If the background processes are not running, then this is a no-op.
func (c *Client) Stop() {
	c.configMu.RLock()
	if !c.running {
		c.configMu.RUnlock()
		return
	}
	c.configMu.RUnlock()

	c.healthcheckStop <- true
	<-c.healthcheckStopped

	c.snifferStop <- true
	<-c.snifferStopped

	c.configMu.Lock()
	c.running = false
	c.configMu.Unlock()

	c.logf("Client stopped")
}

// SetLogger sets the logger for output from Elastic.
// If you set it to nil (default), it will not print anything.
func (c *Client) SetLogger(logger *log.Logger) {
	c.configMu.Lock()
	c.logger = logger
	c.configMu.Unlock()
}

// SetTracer sets the tracer to log HTTP requests to and responses from Elastic.
// If you set it to nil (default), it will not print anything.
func (c *Client) SetTracer(tracer *log.Logger) {
	c.configMu.Lock()
	c.tracer = tracer
	c.configMu.Unlock()
}

// SetMaxRetries sets the maximum number a request is retried.
// If it is <= 0, retrying is disabled (the default).
func (c *Client) SetMaxRetries(maxRetries int) {
	c.configMu.Lock()
	c.maxRetries = maxRetries
	c.configMu.Unlock()
}

// SetSnifferSchedule sets the duration between two sniffer procedures.
// The default duration is 15 minutes.
func (c *Client) SetSnifferSchedule(schedule time.Duration) {
	c.configMu.Lock()
	c.snifferSchedule = schedule
	c.configMu.Unlock()

	// Notify to pick up new setting
	c.snifferScheduleUpdate <- true
}

// SetSnifferTimeout sets the timeout for the sniffer that finds the
// nodes in a cluster. The default is 1 second.
func (c *Client) SetSnifferTimeout(timeout time.Duration) {
	c.configMu.Lock()
	c.snifferTimeout = timeout
	c.configMu.Unlock()
}

// SetHealthcheckSchedule sets the duration between two health checks.
// The default duration is 60 seconds.
func (c *Client) SetHealthcheckSchedule(schedule time.Duration) {
	c.configMu.Lock()
	c.healthcheckSchedule = schedule
	c.configMu.Unlock()

	// Notify to pick up new setting
	c.healthcheckScheduleUpdate <- true
}

// SetDecoder sets the interface to be used for decoding data from
// Elasticsearch. The default is DefaultDecoder.
func (c *Client) SetDecoder(decoder Decoder) {
	c.configMu.Lock()
	if decoder != nil {
		c.decoder = decoder
	} else {
		c.decoder = &DefaultDecoder{}
	}
	c.configMu.Unlock()
}

// logf is a helper to log standard output.
func (c *Client) logf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(format, args...)
	}
}

// tracef is a helper to trace e.g. HTTP requests/responses.
func (c *Client) tracef(format string, args ...interface{}) {
	if c.tracer != nil {
		c.tracer.Printf(format, args...)
	}
}

// dumpRequest dumps the given HTTP request.
func (c *Client) dumpRequest(r *http.Request) {
	if c.tracer != nil {
		out, err := httputil.DumpRequestOut(r, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// dumpResponse dumps the given HTTP response.
func (c *Client) dumpResponse(resp *http.Response) {
	if c.tracer != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			c.tracef("%s\n", string(out))
		}
	}
}

// sniffer periodically runs sniff.
func (c *Client) sniffer() {
	for {
		c.configMu.RLock()
		ticker := time.NewTicker(c.snifferSchedule)
		c.configMu.RUnlock()

		select {
		case <-c.snifferScheduleUpdate:
			// config setting changed
			break
		case <-c.snifferStop:
			// we are asked to stop, so we signal back that we're stopping now
			c.snifferStopped <- true
			return
		case <-ticker.C:
			c.sniff()
		}
	}
}

// sniff uses the Node Info API to return the list of nodes in the cluster
// that host is part of. It uses the list of URLs passed and returns the
// nodes of the first cluster found. It is the responsibility of the caller
// to ensure that the snifferURLs belong to the same cluster.
func (c *Client) sniff(snifferURLs ...string) error {
	if len(snifferURLs) == 0 {
		// Okay, start with the nodes we already have.
		snifferURLs = make([]string, 0)
		c.connsMu.RLock()
		for _, conn := range c.conns {
			if !conn.IsDead() {
				snifferURLs = append(snifferURLs, conn.URL())
			}
		}
		c.connsMu.RUnlock()

		// Still no luck
		if len(snifferURLs) == 0 {
			return ErrNoClient
		}
	}

	ch := make(chan []*conn, len(snifferURLs))

	// Sniff each URL provided, in parallel.
	for _, url := range snifferURLs {
		go func(url string) { ch <- c.sniffNode(url) }(url)
	}

	c.configMu.RLock()
	timeout := c.snifferTimeout
	c.configMu.RUnlock()

	// Wait for the results to come back, or the process times out.
	for {
		select {
		case nodes := <-ch:
			if len(nodes) > 0 {
				c.updateConns(nodes)
				return nil
			}
		case <-time.After(timeout):
			// We get here if no cluster responds in time
			c.logf("Sniffer timeout when reloading connections")
			return ErrNoClient
		}
	}

	return ErrNoClient
}

// updateConns updates the clients' connections with new information
// gather by a sniff operation.
func (c *Client) updateConns(conns []*conn) {
	c.connsMu.Lock()

	newConns := make([]*conn, 0)

	// Build up new connections:
	// If we find an existing connection, take that (including no. of failures etc.).
	// If we find a new connection, use it.
	for _, conn := range conns {
		var found bool
		for _, oldConn := range c.conns {
			if oldConn.NodeID() == conn.NodeID() {
				// Take over the old connection
				newConns = append(newConns, oldConn)
				found = true
				break
			}
		}
		if !found {
			// New connection not there previously
			newConns = append(newConns, conn)
		}
	}

	c.conns = newConns
	c.cindex = -1
	c.connsMu.Unlock()
}

// reSniffHostAndPort is used to extract hostname and port from a result
// from a Nodes Info API.
var reSniffHostAndPort = regexp.MustCompile(`\/([^:]*):([0-9]+)\]`)

// sniffNode sniffs a single node. This method is run as a goroutine
// in sniff. If successful, it returns the list of node URLs extracted
// from the result of calling Nodes Info API. Otherwise, an empty array
// is returned.
func (c *Client) sniffNode(url string) []*conn {
	nodes := make([]*conn, 0)

	// Call the Nodes Info API at /_nodes/http
	req, err := NewRequest("GET", url+"/_nodes/http")
	if err != nil {
		c.logf("Sniffing node %s failed: %v", url, err)
		return nodes
	}

	res, err := c.c.Do((*http.Request)(req))
	if err != nil {
		c.logf("Sniffing node %s failed: %v", url, err)
		return nodes
	}
	if res == nil {
		return nodes
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	var info NodesInfoResponse
	if err := json.NewDecoder(res.Body).Decode(&info); err == nil {
		if len(info.Nodes) > 0 {
			switch c.scheme {
			case "https":
				for nodeID, node := range info.Nodes {
					m := reSniffHostAndPort.FindStringSubmatch(node.HTTPSAddress)
					if len(m) == 3 {
						url := fmt.Sprintf("https://%s:%s", m[1], m[2])
						nodes = append(nodes, newConn(nodeID, url))
					}
				}
			default:
				for nodeID, node := range info.Nodes {
					m := reSniffHostAndPort.FindStringSubmatch(node.HTTPAddress)
					if len(m) == 3 {
						url := fmt.Sprintf("http://%s:%s", m[1], m[2])
						nodes = append(nodes, newConn(nodeID, url))
					}
				}
			}
		}
	}

	return nodes
}

// healthchecker periodically runs healthcheck.
func (c *Client) healthchecker() {
	for {
		c.configMu.RLock()
		ticker := time.NewTicker(c.healthcheckSchedule)
		c.configMu.RUnlock()

		select {
		case <-c.healthcheckScheduleUpdate:
			// pick up new configuration setting
			break
		case <-c.healthcheckStop:
			// we are asked to stop, so we signal back that we're stopping now
			c.healthcheckStopped <- true
			return
		case <-ticker.C:
			c.healthcheck()
		}
	}
}

// healthcheck does a health check on all nodes in the cluster. Depending on
// the node state, it marks connections as dead, sets them alive etc.
func (c *Client) healthcheck() {
	c.connsMu.RLock()
	conns := c.conns
	c.connsMu.RUnlock()

	for _, conn := range conns {
		params := make(url.Values)
		params.Set("timeout", "1")
		req, err := NewRequest("HEAD", conn.URL()+"/?"+params.Encode())
		if err == nil {
			res, err := c.c.Do((*http.Request)(req))
			if err == nil {
				if res.Body != nil {
					defer res.Body.Close()
				}
				if res.StatusCode >= 200 && res.StatusCode < 300 {
					conn.MarkAsAlive()
				} else {
					conn.MarkAsDead()
					c.logf("Mark %s as dead [status=%d]", conn.URL(), res.StatusCode)
				}
			} else {
				c.logf("Mark %s as dead: %v", conn.URL(), err)
				conn.MarkAsDead()
			}
		} else {
			c.logf("Mark %s as dead: %v", conn.URL(), err)
			conn.MarkAsDead()
		}
	}
}

// next returns the next available connection, or ErrNoClient.
func (c *Client) next() (*conn, error) {
	// We do round-robin here.
	// TODO: This should be a pluggable strategy, like the Selector in the official clients.
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	i := 0
	numConns := len(c.conns)
	for {
		i += 1
		if i > numConns {
			break
		}

		c.cindex += 1
		if c.cindex >= numConns {
			c.cindex = 0
		}
		conn := c.conns[c.cindex]
		if !conn.IsDead() {
			return conn, nil
		}
	}

	// We tried hard, but there is no node available
	return nil, ErrNoClient
}

// PerformRequest does a HTTP request to Elasticsearch while logging, tracing,
// marking dead connections, retrying, and reloading connections.
//
// It returns a response and an error on failure.
func (c *Client) PerformRequest(method, path string, params url.Values, body interface{}) (*Response, error) {
	start := time.Now().UTC()

	c.configMu.RLock()
	retries := c.maxRetries
	c.configMu.RUnlock()

	var err error
	var conn *conn
	var req *Request
	var resp *Response

	// Maybe make this configurable
	sleepBetweenTimeouts := 100 * time.Millisecond

	for {
		pathWithParams := path
		if len(params) > 0 {
			pathWithParams += "?" + params.Encode()
		}

		// Get a connection
		conn, err = c.next()
		if err != nil {
			c.logf("Cannot get new connection from pool")
			return nil, err // only retry failed HTTP requests
		}

		// Set up a new request
		req, err = NewRequest(method, conn.URL()+pathWithParams)
		if err != nil {
			c.logf("Error creating new request: %s %s: %v", strings.ToUpper(method), conn.URL()+pathWithParams, err)
			return nil, err // only retry failed HTTP requests
		}

		// Set body
		if body != nil {
			switch b := body.(type) {
			case string:
				req.SetBodyString(b)
				break
			default:
				req.SetBodyJson(body)
				break
			}
		}

		// Tracing
		c.dumpRequest((*http.Request)(req))

		// Get response
		res, err := c.c.Do((*http.Request)(req))
		if err != nil {
			retries -= 1
			if retries <= 0 {
				c.logf("Mark %s as dead", conn.URL())
				conn.MarkAsDead() // mark connection as dead
				return nil, err
			}
			time.Sleep(sleepBetweenTimeouts)
			continue // try again
		}
		if res.Body != nil {
			defer res.Body.Close()
		}

		// Check for errors
		if err := checkResponse(res); err != nil {
			retries -= 1
			if retries <= 0 {
				return nil, err
			}
			time.Sleep(sleepBetweenTimeouts)
			continue // try again
		}

		// Tracing
		c.dumpResponse(res)

		conn.MarkAsHealthy()

		resp, err = c.newResponse(res)
		if err != nil {
			return nil, err
		}

		break
	}

	if c.logger != nil {
		duration := time.Now().UTC().Sub(start)
		c.logf("%s %s [status:%d, request:%.3fs]",
			strings.ToUpper(method),
			req.URL,
			resp.StatusCode,
			float64(int64(duration/time.Millisecond))/1000)
	}

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
func (c *Client) GetTemplate() *GetTemplateService {
	return NewGetTemplateService(c)
}

// PutTemplate creates or updates a search template.
func (c *Client) PutTemplate() *PutTemplateService {
	return NewPutTemplateService(c)
}

// DeleteTemplate deletes a search template.
func (c *Client) DeleteTemplate() *DeleteTemplateService {
	return NewDeleteTemplateService(c)
}

// ClusterHealth retrieves the health of the cluster.
func (c *Client) ClusterHealth() *ClusterHealthService {
	return NewClusterHealthService(c)
}

// ClusterState retrieves the state of the cluster.
func (c *Client) ClusterState() *ClusterStateService {
	return NewClusterStateService(c)
}

// NodesInfo retrieves one or more or all of the cluster nodes information.
func (c *Client) NodesInfo() *NodesInfoService {
	return NewNodesInfoService(c)
}
