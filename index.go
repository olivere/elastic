// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// The result of indexing a document in ElasticSearch.
type IndexResult struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Created bool   `json:"created"`
}

// The Index service adds documents to ElasticSearch.
type IndexService struct {
	client     *Client
	index      string
	_type      string
	id         string
	bodyString string
	bodyJson   interface{}
	pretty     bool
	debug      bool
}

func NewIndexService(client *Client) *IndexService {
	builder := &IndexService{
		client: client,
	}
	return builder
}

func (b *IndexService) Index(name string) *IndexService {
	b.index = name
	return b
}

func (b *IndexService) Type(_type string) *IndexService {
	b._type = _type
	return b
}

func (b *IndexService) Id(id string) *IndexService {
	b.id = id
	return b
}

func (b *IndexService) BodyString(body string) *IndexService {
	b.bodyString = body
	return b
}

func (b *IndexService) BodyJson(json interface{}) *IndexService {
	b.bodyJson = json
	return b
}

func (b *IndexService) Pretty(pretty bool) *IndexService {
	b.pretty = pretty
	return b
}

func (b *IndexService) Debug(debug bool) *IndexService {
	b.debug = debug
	return b
}

func (b *IndexService) Do() (*IndexResult, error) {
	// Build url
	urls := "/{index}/{type}/{id}"
	urls = strings.Replace(urls, "{index}", cleanPathString(b.index), 1)
	urls = strings.Replace(urls, "{type}", cleanPathString(b._type), 1)
	urls = strings.Replace(urls, "{id}", cleanPathString(b.id), 1)

	// Parameters
	params := make(url.Values)
	if b.pretty {
		params.Set("pretty", fmt.Sprintf("%v", b.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := b.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	if b.bodyJson != nil {
		req.SetBodyJson(b.bodyJson)
	} else {
		req.SetBodyString(b.bodyString)
	}

	if b.debug {
		out, _ := httputil.DumpRequestOut((*http.Request)(req), true)
		log.Printf("%s\n", string(out))
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
		out, _ := httputil.DumpResponse(res, true)
		log.Printf("%s\n", string(out))
	}

	ret := new(IndexResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
