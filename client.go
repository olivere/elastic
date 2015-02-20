// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const (
	// Version is the current version of Elastic.
	Version = "1.3.1"

	// defaultUrl to be used as base for Elasticsearch requests.
	defaultUrl = "http://localhost:9200"

	// pingDuration is the time to periodically check the Elasticsearch URLs.
	pingDuration = 60 * time.Second
)

var (
	// ErrNoClient is raised when no active Elasticsearch client is available.
	ErrNoClient = errors.New("no active client")
)

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	urls []string // urls is a list of all clients for Elasticsearch queries

	c *http.Client // c is the net/http Client to use for requests

	log *log.Logger // output log

	mu        sync.RWMutex // mutex for the next two fields
	activeUrl string       // currently active connection url
	hasActive bool         // true if we have an active connection
}

// NewClient creates a new client to work with Elasticsearch.
func NewClient(client *http.Client, urls ...string) (*Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	c := &Client{c: client}
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
	c.pingUrls()
	go c.pinger() // start goroutine periodically ping all clients
	return c, nil
}

// SetLogger sets the logger for output from Elastic.
// If you don't set the logger, it will print to os.Stdout.
func (c *Client) SetLogger(log *log.Logger) {
	c.log = log
}

// printf is a helper to log output.
func (c *Client) printf(format string, args ...interface{}) {
	if c.log != nil {
		c.log.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// dumpRequest dumps the given HTTP request.
func (c *Client) dumpRequest(r *http.Request) {
	out, err := httputil.DumpRequestOut(r, true)
	if err == nil {
		c.printf("%s\n", string(out))
	}
}

// dumpResponse dumps the given HTTP response.
func (c *Client) dumpResponse(resp *http.Response) {
	out, err := httputil.DumpResponse(resp, true)
	if err == nil {
		c.printf("%s\n", string(out))
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
	return NewRequest(method, c.activeUrl+path)
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
				defer res.Body.Close()
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
				log.Printf("elastic: %v", err)
			}
		} else {
			log.Printf("elastic: %v", err)
		}
	}

	// No client available
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hasActive = false
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
