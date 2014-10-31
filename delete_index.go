// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"net/http"

	"github.com/olivere/elastic/uritemplates"
)

type DeleteIndexService struct {
	client *Client
	index  string
}

func NewDeleteIndexService(client *Client) *DeleteIndexService {
	builder := &DeleteIndexService{
		client: client,
	}
	return builder
}

func (b *DeleteIndexService) Index(index string) *DeleteIndexService {
	b.index = index
	return b
}

func (b *DeleteIndexService) Do() (*DeleteIndexResult, error) {
	// Build url
	urls, err := uritemplates.Expand("/{index}/", map[string]string{
		"index": b.index,
	})
	if err != nil {
		return nil, err
	}

	// Set up a new request
	req, err := b.client.NewRequest("DELETE", urls)
	if err != nil {
		return nil, err
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
	ret := new(DeleteIndexResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a delete index request.

type DeleteIndexResult struct {
	Acknowledged bool `json:"acknowledged"`
}
