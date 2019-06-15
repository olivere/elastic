// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// IndicesCreateService creates a new index.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.7/indices-create-index.html
// for details.
type IndicesCreateService struct {
	client          *Client
	pretty          bool
	index           string
	timeout         string
	masterTimeout   string
	includeTypeName *bool
	bodyJson        interface{}
	bodyString      string
}

// NewIndicesCreateService returns a new IndicesCreateService.
func NewIndicesCreateService(client *Client) *IndicesCreateService {
	return &IndicesCreateService{client: client}
}

// Index is the name of the index to create.
func (b *IndicesCreateService) Index(index string) *IndicesCreateService {
	b.index = index
	return b
}

// Timeout the explicit operation timeout, e.g. "5s".
func (s *IndicesCreateService) Timeout(timeout string) *IndicesCreateService {
	s.timeout = timeout
	return s
}

// MasterTimeout specifies the timeout for connection to master.
func (s *IndicesCreateService) MasterTimeout(masterTimeout string) *IndicesCreateService {
	s.masterTimeout = masterTimeout
	return s
}

// Body specifies the configuration of the index as a string.
// It is an alias for BodyString.
func (b *IndicesCreateService) Body(body string) *IndicesCreateService {
	b.bodyString = body
	return b
}

// BodyString specifies the configuration of the index as a string.
func (b *IndicesCreateService) BodyString(body string) *IndicesCreateService {
	b.bodyString = body
	return b
}

// BodyJson specifies the configuration of the index. The interface{} will
// be serializes as a JSON document, so use a map[string]interface{}.
func (b *IndicesCreateService) BodyJson(body interface{}) *IndicesCreateService {
	b.bodyJson = body
	return b
}

// Pretty indicates that the JSON response be indented and human readable.
func (b *IndicesCreateService) Pretty(pretty bool) *IndicesCreateService {
	b.pretty = pretty
	return b
}

// IncludeTypeName specifies whether requests and responses should include a type name.
func (s *IndicesCreateService) IncludeTypeName(includeTypeName bool) *IndicesCreateService {
	s.includeTypeName = &includeTypeName
	return s
}

// Do executes the operation.
func (b *IndicesCreateService) Do(ctx context.Context) (*IndicesCreateResult, error) {
	if b.index == "" {
		return nil, errors.New("missing index name")
	}

	// Build url
	path, err := uritemplates.Expand("/{index}", map[string]string{
		"index": b.index,
	})
	if err != nil {
		return nil, err
	}

	// Setup HTTP request body
	var body interface{}
	if b.bodyJson != nil {
		body = b.bodyJson
	} else {
		body = b.bodyString
	}

	// Get response
	res, err := b.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "PUT",
		Path:   path,
		Params: b.buildParams(),
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	ret := new(IndicesCreateResult)
	if err := b.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (b *IndicesCreateService) buildParams() url.Values {
	params := make(url.Values)
	if b.pretty {
		params.Set("pretty", "true")
	}
	if b.masterTimeout != "" {
		params.Set("master_timeout", b.masterTimeout)
	}
	if b.timeout != "" {
		params.Set("timeout", b.timeout)
	}
	if b.includeTypeName != nil {
		params.Set("include_type_name", fmt.Sprintf("%v", *b.includeTypeName))
	}
	return params
}

// -- Result of a create index request.

// IndicesCreateResult is the outcome of creating a new index.
type IndicesCreateResult struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index,omitempty"`
}
