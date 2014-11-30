// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"errors"
	"net/http"
)

const (
	// Version is the current version of Elastic.
	Version = "1.3.1"

	// defaultUrl to be used as base for Elasticsearch requests.
	defaultUrl = "http://localhost:9200"
)

var (
	// ErrNoClient is raised when no active Elasticsearch client is available.
	ErrNoClient = errors.New("no active client")
)

// Client is an Elasticsearch client. Create one by calling NewClient.
type Client struct {
	c    *http.Client    // c is the net/http Client to use for requests
	pool *ConnectionPool // pool is a list of connections to Elasticsearch servers
}

// NewClient creates a new client to work with Elasticsearch.
func NewClient(client *http.Client, urls ...string) (*Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	pool := NewConnectionPool(client, urls...)
	c := &Client{c: client, pool: pool}
	return c, nil
}

// NewRequest creates a new request with the given method and prepends
// the base URL to the path. If no active connection to Elasticsearch
// is available, ErrNoClient is returned.
func (c *Client) NewRequest(method, path string) (*Request, error) {
	url, err := c.pool.GetNextRequestURL()
	if err != nil {
		return nil, ErrNoClient
	}
	return NewRequest(method, url+path)
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
