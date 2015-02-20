// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
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

	// defaultUrl to be used as base for Elasticsearch requests.
	defaultUrl = "http://127.0.0.1:9200"

	// pingDuration is the time to periodically check the Elasticsearch URLs.
	pingDuration = 60 * time.Second
)

var (
	// ErrNoClient is raised when no active Elasticsearch client is available.
	ErrNoClient = errors.New("no active client")

	// ErrRetry is raised when a request cannot be executed after the defined
	// number of retries.
	ErrRetry = errors.New("cannot connect after several retries")
)

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	urls []string // urls is a list of all clients for Elasticsearch queries

	c *http.Client // c is the net/http Client to use for requests

	logger *log.Logger // standard log
	tracer *log.Logger // trace log

	maxRetries int // max. number of retries

	protocol       string        // http or https
	snifferTimeout time.Duration // time the sniffer waits for a response from nodes info API

	decoder Decoder // used to decode data sent from Elasticsearch

	mu        sync.RWMutex // mutex for the next two fields
	activeUrl string       // currently active connection url
	hasActive bool         // true if we have an active connection
}

// NewClient creates a new client to work with Elasticsearch.
func NewClient(client *http.Client, urls ...string) (*Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	c := &Client{c: client, protocol: "http", decoder: &DefaultDecoder{}}
	switch len(urls) {
	case 0:
		c.urls = make([]string, 1)
		c.urls[0] = defaultUrl
	case 1:
		c.urls = make([]string, 1)
		c.urls[0] = urls[0]
	default:
		c.urls = urls
	}

	if len(c.urls) == 1 {
		// If we specify just one URL, we assume that sniffing the cluster
		// with the nodes info API is ok. This is what the offical clients do.
		urls := c.sniff()
		if len(urls) == 0 {
			return nil, errors.New("no nodes found in cluster")
		}
		c.urls = urls
	} else if len(c.urls) >= 2 {
		// If we provide 2 or more URLs, then we assume that the caller knows
		// what he does and doesn't want us to sniff.
	}

	c.pingUrls()
	go c.pinger() // start goroutine periodically ping all clients
	return c, nil
}

// SetLogger sets the logger for output from Elastic.
// If you set it to nil (default), it will not print anything.
func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
}

// SetTracer sets the tracer to log HTTP requests to and responses from Elastic.
// If you set it to nil (default), it will not print anything.
func (c *Client) SetTracer(tracer *log.Logger) {
	c.tracer = tracer
}

// SetMaxRetries sets the maximum number a request is retried.
// If it is <= 0, retrying is disabled (default).
func (c *Client) SetMaxRetries(maxRetries int) {
	c.maxRetries = maxRetries
}

// SetSnifferTimeout sets the timeout for the sniffer that finds the
// nodes in a cluster. The default is 1 second.
func (c *Client) SetSnifferTimeout(timeout time.Duration) {
	if int64(timeout*time.Second) < 1 {
		c.snifferTimeout = 1 * time.Second
	} else {
		c.snifferTimeout = timeout
	}
}

// SetDecoder sets the interface to be used for decoding data from
// Elasticsearch. The default is DefaultDecoder.
func (c *Client) SetDecoder(decoder Decoder) {
	if decoder != nil {
		c.decoder = decoder
	} else {
		c.decoder = &DefaultDecoder{}
	}
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

// NewRequest creates a new request with the given method and prepends
// the base URL to the path. If no active connection to Elasticsearch
// is available, ErrNoClient is returned.
func (c *Client) NewRequest(method, path string) (*Request, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.hasActive {
		return nil, ErrNoClient
	}
	url := c.activeUrl // c.selector.Select() // get the next URL to use
	return NewRequest(method, url+path)
}

// pinger periodically runs pingUrls.
func (c *Client) pinger() {
	ticker := time.NewTicker(pingDuration)
	for {
		select {
		case <-ticker.C:
			c.pingUrls()
		}
	}
}

// sniff uses the Node Info API to return the list of nodes in the cluster
// that host is part of. It returns the
func (c *Client) sniff() []string {
	timeout := c.snifferTimeout
	if int64(timeout*time.Second) < 1 {
		timeout = 1 * time.Second
	}

	sniffCh := make(chan sniffResult, 1)

	// Sniff each url provided, in parallel.
	for _, url := range c.urls {
		go func() { sniffCh <- c.sniffNode(url) }()
	}

	select {
	case res := <-sniffCh:
		if len(res.URLs) > 0 {
			return res.URLs
		}
		break
	case <-time.After(timeout):
		break
	}

	// We get here if no cluster responds in time
	return []string{}
}

type sniffResult struct {
	URLs []string
}

// sniffNode sniffs a single node.
func (c *Client) sniffNode(url string) sniffResult {
	re := regexp.MustCompile(`\/([^:]*):([0-9]+)\]`)

	req, err := NewRequest("GET", url+"/_nodes/http")
	if err == nil {
		res, err := c.c.Do((*http.Request)(req))
		if err == nil && res != nil {
			if res.Body != nil {
				defer res.Body.Close()
			}
			var info NodesInfoResponse
			if err := json.NewDecoder(res.Body).Decode(&info); err == nil {
				if len(info.Nodes) > 0 {
					var urls []string
					switch c.protocol {
					case "https":
						for _, node := range info.Nodes {
							m := re.FindStringSubmatch(node.HTTPSAddress)
							if len(m) == 3 {
								urls = append(urls, fmt.Sprintf("https://%s:%s", m[1], m[2]))
							}
						}
						break
					default:
						for _, node := range info.Nodes {
							m := re.FindStringSubmatch(node.HTTPAddress)
							if len(m) == 3 {
								urls = append(urls, fmt.Sprintf("http://%s:%s", m[1], m[2]))
							}
						}
						break
					}
					return sniffResult{URLs: urls}
				}
			}
		}
	}
	return sniffResult{}
}

// pingUrls iterates through all client URLs. It checks if the client
// is available. It takes the first one available and saves its URL
// in activeUrl. If no client is available, hasActive is set to false
// and NewRequest will fail.
func (c *Client) pingUrls() {
	for _, url_ := range c.urls {
		params := make(url.Values)
		params.Set("timeout", "1")
		req, err := NewRequest("HEAD", url_+"/?"+params.Encode())
		if err == nil {
			res, err := c.c.Do((*http.Request)(req))
			if err == nil {
				if res.Body != nil {
					defer res.Body.Close()
				}
				if res.StatusCode == http.StatusOK {
					// Everything okay: Update activeUrl and set hasActive to true.
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.activeUrl != "" && c.activeUrl != url_ {
						log.Printf("elastic: switched connection from %s to %s", c.activeUrl, url_)
					}
					c.activeUrl = url_
					c.hasActive = true
					return
				}
			} else {
				c.logf("elastic: %v", err)
			}
		} else {
			c.logf("elastic: %v", err)
		}
	}

	// No client available
	c.mu.Lock()
	c.hasActive = false
	c.mu.Unlock()
}

// PerformRequest does a HTTP request to Elasticsearch while logging, tracing,
// marking dead connections, retrying, and reloading connections.
//
// It returns a response and an error on failure.
func (c *Client) PerformRequest(method, path string, params url.Values, body interface{}) (*Response, error) {
	start := time.Now().UTC()
	retries := c.maxRetries

	var err error
	var req *Request
	var resp *Response

	for {
		pathWithParams := path
		if len(params) > 0 {
			pathWithParams += "?" + params.Encode()
		}

		// Set up a new request
		req, err = c.NewRequest(method, pathWithParams)
		if err != nil {
			return nil, err
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
				return nil, err
			}
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
			continue // try again
		}

		// Tracing
		c.dumpResponse(res)

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
			float64(duration/time.Second))
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
