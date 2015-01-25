// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"net/http"

	"github.com/olivere/elastic/uritemplates"
)

type CreateIndexService struct {
	client *Client
	index  string
	body   string
	pretty bool
	debug  bool
}

func NewCreateIndexService(client *Client) *CreateIndexService {
	builder := &CreateIndexService{
		client: client,
	}
	return builder
}

func (b *CreateIndexService) Index(index string) *CreateIndexService {
	b.index = index
	return b
}

func (b *CreateIndexService) Body(body string) *CreateIndexService {
	b.body = body
	return b
}

func (b *CreateIndexService) Pretty(pretty bool) *CreateIndexService {
	b.pretty = pretty
	return b
}

func (b *CreateIndexService) Debug(debug bool) *CreateIndexService {
	b.debug = debug
	return b
}

func (b *CreateIndexService) Do() (*CreateIndexResult, error) {
	// Build url
	urls, err := uritemplates.Expand("/{index}/", map[string]string{
		"index": b.index,
	})
	if err != nil {
		return nil, err
	}

	// Set up a new request
	req, err := b.client.NewRequest("PUT", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	req.SetBodyString(b.body)

	if b.debug {
		b.client.dumpRequest((*http.Request)(req))
	}

	// Get response
	res, err := b.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if b.debug {
		b.client.dumpResponse(res)
	}

	ret := new(CreateIndexResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a create index request.

type CreateIndexResult struct {
	Acknowledged bool `json:"acknowledged"`
}
