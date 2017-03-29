// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"gopkg.in/olivere/elastic.v3/uritemplates"
)

// PutService put a document in Elasticsearch.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/docs-update.html
// for details.
type PutService struct {
	client           *Client
	index            string
	typ              string
	id               string
	routing          string
	parent           string
	fields           []string
	version          *int64
	versionType      string
	retryOnConflict  *int
	refresh          *bool
	replicationType  string
	consistencyLevel string
	doc              interface{}
	timeout          string
	pretty           bool
}

// NewPutService creates the service to replace documents in Elasticsearch.
func NewPutService(client *Client) *PutService {
	builder := &PutService{
		client: client,
		fields: make([]string, 0),
	}
	return builder
}

// Index is the name of the Elasticsearch index (required).
func (b *PutService) Index(name string) *PutService {
	b.index = name
	return b
}

// Type is the type of the document (required).
func (b *PutService) Type(typ string) *PutService {
	b.typ = typ
	return b
}

// Id is the identifier of the document to update (required).
func (b *PutService) Id(id string) *PutService {
	b.id = id
	return b
}

// Routing specifies a specific routing value.
func (b *PutService) Routing(routing string) *PutService {
	b.routing = routing
	return b
}

// Parent sets the id of the parent document.
func (b *PutService) Parent(parent string) *PutService {
	b.parent = parent
	return b
}

// RetryOnConflict specifies how many times the operation should be retried
// when a conflict occurs (default: 0).
func (b *PutService) RetryOnConflict(retryOnConflict int) *PutService {
	b.retryOnConflict = &retryOnConflict
	return b
}

// Fields is a list of fields to return in the response.
func (b *PutService) Fields(fields ...string) *PutService {
	b.fields = make([]string, 0, len(fields))
	b.fields = append(b.fields, fields...)
	return b
}

// Version defines the explicit version number for concurrency control.
func (b *PutService) Version(version int64) *PutService {
	b.version = &version
	return b
}

// VersionType is one of "internal" or "force".
func (b *PutService) VersionType(versionType string) *PutService {
	b.versionType = versionType
	return b
}

// Refresh the index after performing the update.
func (b *PutService) Refresh(refresh bool) *PutService {
	b.refresh = &refresh
	return b
}

// ReplicationType is one of "sync" or "async".
func (b *PutService) ReplicationType(replicationType string) *PutService {
	b.replicationType = replicationType
	return b
}

// ConsistencyLevel is one of "one", "quorum", or "all".
// It sets the write consistency setting for the update operation.
func (b *PutService) ConsistencyLevel(consistencyLevel string) *PutService {
	b.consistencyLevel = consistencyLevel
	return b
}

// Doc allows for updating a partial document.
func (b *PutService) Doc(doc interface{}) *PutService {
	b.doc = doc
	return b
}

// Timeout is an explicit timeout for the operation, e.g. "1000", "1s" or "500ms".
func (b *PutService) Timeout(timeout string) *PutService {
	b.timeout = timeout
	return b
}

// Pretty instructs to return human readable, prettified JSON.
func (b *PutService) Pretty(pretty bool) *PutService {
	b.pretty = pretty
	return b
}

// url returns the URL part of the document request.
func (b *PutService) url() (string, url.Values, error) {
	// Build url
	path := "/{index}/{type}/{id}"
	path, err := uritemplates.Expand(path, map[string]string{
		"index": b.index,
		"type":  b.typ,
		"id":    b.id,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Parameters
	params := make(url.Values)
	if b.pretty {
		params.Set("pretty", "true")
	}
	if b.routing != "" {
		params.Set("routing", b.routing)
	}
	if b.parent != "" {
		params.Set("parent", b.parent)
	}
	if b.timeout != "" {
		params.Set("timeout", b.timeout)
	}
	if b.refresh != nil {
		params.Set("refresh", fmt.Sprintf("%v", *b.refresh))
	}
	if b.replicationType != "" {
		params.Set("replication", b.replicationType)
	}
	if b.consistencyLevel != "" {
		params.Set("consistency", b.consistencyLevel)
	}
	if len(b.fields) > 0 {
		params.Set("fields", strings.Join(b.fields, ","))
	}
	if b.version != nil {
		params.Set("version", fmt.Sprintf("%d", *b.version))
	}
	if b.versionType != "" {
		params.Set("version_type", b.versionType)
	}
	if b.retryOnConflict != nil {
		params.Set("retry_on_conflict", fmt.Sprintf("%v", *b.retryOnConflict))
	}

	return path, params, nil
}

// body returns the body part of the document request.
func (b *PutService) body() (interface{}, error) {
	return b.doc, nil
}

// Do executes the update operation.
func (b *PutService) Do() (*PutResponse, error) {
	return b.DoC(nil)
}

// DoC executes the update operation.
func (b *PutService) DoC(ctx context.Context) (*PutResponse, error) {
	path, params, err := b.url()
	if err != nil {
		return nil, err
	}

	// Get body of the request
	body, err := b.body()
	if err != nil {
		return nil, err
	}

	// Get response
	res, err := b.client.PerformRequestC(ctx, "PUT", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return result
	ret := new(PutResponse)
	if err := b.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// PutResponse is the result of updating a document in Elasticsearch.
type PutResponse struct {
	Index     string     `json:"_index"`
	Type      string     `json:"_type"`
	Id        string     `json:"_id"`
	Version   int        `json:"_version"`
	Created   bool       `json:"created"`
	GetResult *GetResult `json:"get"`
}
