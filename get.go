// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

type GetService struct {
	client                     *Client
	index                      string
	_type                      string
	id                         string
	routing                    string
	preference                 string
	fields                     []string
	fsc                        *FetchSourceContext
	refresh                    *bool
	realtime                   *bool
	version                    *int64 // see org.elasticsearch.common.lucene.uid.Versions
	versionType                string // see org.elasticsearch.index.VersionType
	ignoreErrOnGeneratedFields *bool
}

func NewGetService(client *Client) *GetService {
	builder := &GetService{
		client: client,
		_type:  "_all",
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

func (s *GetService) FetchSource(fetchSource bool) *GetService {
	if s.fsc == nil {
		s.fsc = NewFetchSourceContext(fetchSource)
	} else {
		s.fsc.SetFetchSource(fetchSource)
	}
	return s
}

func (s *GetService) FetchSourceContext(fetchSourceContext *FetchSourceContext) *GetService {
	s.fsc = fetchSourceContext
	return s
}

func (b *GetService) Refresh(refresh bool) *GetService {
	b.refresh = &refresh
	return b
}

func (b *GetService) Realtime(realtime bool) *GetService {
	b.realtime = &realtime
	return b
}

// Version can be MatchAny (-3), MatchAnyPre120 (0), NotFound (-1),
// or NotSet (-2). These are specified in org.elasticsearch.common.lucene.uid.Versions.
// The default is MatchAny (-3).
func (b *GetService) Version(version int64) *GetService {
	b.version = &version
	return b
}

// VersionType can be "internal", "external", "external_gt", "external_gte",
// or "force". See org.elasticsearch.index.VersionType in Elasticsearch source.
// It is "internal" by default.
func (b *GetService) VersionType(versionType string) *GetService {
	b.versionType = versionType
	return b
}

func (b *GetService) Do() (*GetResult, error) {
	// Build url
	urls, err := uritemplates.Expand("/{index}/{type}/{id}", map[string]string{
		"index": b.index,
		"type":  b._type,
		"id":    b.id,
	})
	if err != nil {
		return nil, err
	}

	params := make(url.Values)
	if b.realtime != nil {
		params.Add("realtime", fmt.Sprintf("%v", *b.realtime))
	}
	if len(b.fields) > 0 {
		params.Add("fields", strings.Join(b.fields, ","))
	}
	if b.routing != "" {
		params.Add("routing", b.routing)
	}
	if b.preference != "" {
		params.Add("preference", b.preference)
	}
	if b.refresh != nil {
		params.Add("refresh", fmt.Sprintf("%v", *b.refresh))
	}
	if b.version != nil {
		params.Add("_version", fmt.Sprintf("%d", *b.version))
	}
	if b.versionType != "" {
		params.Add("_version_type", b.versionType)
	}
	if b.fsc != nil {
		for k, values := range b.fsc.Query() {
			params.Add(k, strings.Join(values, ","))
		}
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := b.client.NewRequest("GET", urls)
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
	ret := new(GetResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a get request.

type GetResult struct {
	Index   string           `json:"_index"`
	Type    string           `json:"_type"`
	Id      string           `json:"_id"`
	Version int64            `json:"_version,omitempty"`
	Source  *json.RawMessage `json:"_source,omitempty"`
	Found   bool             `json:"found,omitempty"`
	Fields  []string         `json:"fields,omitempty"`
	Error   string           `json:"error,omitempty"` // used only in MultiGet
}
