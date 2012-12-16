// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"errors"
	"net/http"
)

// Current version.
const Version = "0.1"

// Default URL to be used as base for ElasticSearch requests.
const defaultUrl = "http://localhost:9200"

// An ElasticSearch client.
type Client struct {
	// Base URL for ElasticSearch queries
	Url string

	c *http.Client
}

// Creates a new client to work with ElasticSearch.
func NewClient(client *http.Client) (*Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	c := &Client{Url: defaultUrl, c: client}
	return c, nil
}

// Helper function that adds e.g. the base URL to the request.
func (c *Client) NewRequest(method, path string) (*Request, error) {
	return NewRequest(method, c.Url+path)
}

// Create a new index.
func (c *Client) CreateIndex(name string) *CreateIndexService {
	builder := NewCreateIndexService(c)
	builder.Index(name)
	return builder
}

// Delete an index.
func (c *Client) DeleteIndex(name string) *DeleteIndexService {
	builder := NewDeleteIndexService(c)
	builder.Index(name)
	return builder
}

// Check if an index exists.
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

// Delete a document.
func (c *Client) Delete() *DeleteService {
	builder := NewDeleteService(c)
	return builder
}

// Delete documents by query.
func (c *Client) DeleteByQuery() *DeleteByQueryService {
	builder := NewDeleteByQueryService(c)
	return builder
}

// Get a document.
func (c *Client) Get() *GetService {
	builder := NewGetService(c)
	return builder
}

// Check if a document exists.
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

// Start a search.
func (c *Client) Search(indices ...string) *SearchService {
	builder := NewSearchService(c)
	builder.Indices(indices...)
	return builder
}

// Flush.
func (c *Client) Flush() *FlushService {
	builder := NewFlushService(c)
	return builder
}

// Bulk.
func (c *Client) Bulk() *BulkService {
	builder := NewBulkService(c)
	return builder
}

// Aliases.
func (c *Client) Alias() *AliasService {
	builder := NewAliasService(c)
	return builder
}
