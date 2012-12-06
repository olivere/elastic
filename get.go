// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type GetService struct {
	client     *Client
	index      string
	_type      string
	id         string
	routing    string
	preference string
	fields     []string
	refresh    bool
	realtime   bool
}

func NewGetService(client *Client) *GetService {
	builder := &GetService{
		client:   client,
		_type:    "_all",
		realtime: true,
	}
	return builder
}

func (b *GetService) String() string {
	return fmt.Sprintf("[%v][%v][%v]: routing [%v]",
		b.index,
		b._type,
		b.id,
		b.routing)
}

func (b *GetService) Index(index string) *GetService {
	b.index = index
	return b
}

func (b *GetService) Type(_type string) *GetService {
	b._type = _type
	return b
}

func (b *GetService) Id(id string) *GetService {
	b.id = id
	return b
}

func (b *GetService) Parent(parent string) *GetService {
	if b.routing == "" {
		b.routing = parent
	}
	return b
}

func (b *GetService) Routing(routing string) *GetService {
	b.routing = routing
	return b
}

func (b *GetService) Preference(preference string) *GetService {
	b.preference = preference
	return b
}

func (b *GetService) Fields(fields ...string) *GetService {
	if b.fields == nil {
		b.fields = make([]string, 0)
	}
	b.fields = append(b.fields, fields...)
	return b
}

func (b *GetService) Refresh(refresh bool) *GetService {
	b.refresh = refresh
	return b
}

func (b *GetService) Realtime(realtime bool) *GetService {
	b.realtime = realtime
	return b
}

func (b *GetService) Do() (*GetResult, error) {
	// Build url
	urls := "/{index}/{type}/{id}"
	urls = strings.Replace(urls, "{index}", cleanPathString(b.index), 1)
	urls = strings.Replace(urls, "{type}", cleanPathString(b._type), 1)
	urls = strings.Replace(urls, "{id}", cleanPathString(b.id), 1)

	params := make(url.Values)
	urls += params.Encode()

	// Set up a new request
	req, err := b.client.NewRequest("GET", urls)
	if err != nil {
		return nil, err
	}

	// No body in get requests

	// Get response
	res, err := b.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	ret := new(GetResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a get request.

type GetResult struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Source json.RawMessage `json:"source,omitempty"`
}
