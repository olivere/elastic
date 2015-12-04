// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

var (
	_ = fmt.Print
	_ = httputil.DumpRequest
	_ = log.Print
	_ = strings.Index
	_ = uritemplates.Expand
	_ = url.Parse
)

// GetWarmerService allows to get the definition of a warmer.
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-warmers.html.
type GetWarmerService struct {
	client            *Client
	pretty            bool
	index             []string
	name              []string
	typ               []string
	allowNoIndices    *bool
	expandWildcards   string
	ignoreUnavailable *bool
	local             *bool
}

// NewGetWarmerService creates a new GetWarmerService.
func NewGetWarmerService(client *Client) *GetWarmerService {
	return &GetWarmerService{
		client: client,
		typ:    make([]string, 0),
		index:  make([]string, 0),
		name:   make([]string, 0),
	}
}

// Index is a list of index names to restrict the operation; use `_all` to perform the operation on all indices.
func (s *GetWarmerService) Index(indices ...string) *GetWarmerService {
	s.index = append(s.index, indices...)
	return s
}

// Name is the name of the warmer (supports wildcards); leave empty to get all warmers.
func (s *GetWarmerService) Name(name ...string) *GetWarmerService {
	s.name = append(s.name, name...)
	return s
}

// Type is a list of type names the mapping should be added to
// (supports wildcards); use `_all` or omit to add the mapping on all types.
func (s *GetWarmerService) Type(typ ...string) *GetWarmerService {
	s.typ = append(s.typ, typ...)
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices.
// This includes `_all` string or when no indices have been specified.
func (s *GetWarmerService) AllowNoIndices(allowNoIndices bool) *GetWarmerService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *GetWarmerService) ExpandWildcards(expandWildcards string) *GetWarmerService {
	s.expandWildcards = expandWildcards
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *GetWarmerService) IgnoreUnavailable(ignoreUnavailable bool) *GetWarmerService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// Local indicates wether or not to return local information,
// do not retrieve the state from master node (default: false).
func (s *GetWarmerService) Local(local bool) *GetWarmerService {
	s.local = &local
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *GetWarmerService) Pretty(pretty bool) *GetWarmerService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *GetWarmerService) buildURL() (string, url.Values, error) {
	var index, typ, name []string

	if len(s.index) > 0 {
		index = s.index
	} else {
		index = []string{"_all"}
	}

	if len(s.typ) > 0 {
		typ = s.typ
	} else {
		typ = []string{"*"}
	}

	if len(s.name) > 0 {
		name = s.name
	} else {
		name = []string{"_all"}
	}

	// Build URL
	path, err := uritemplates.Expand("/{index}/{type}/_warmer/{name}", map[string]string{
		"index": strings.Join(index, ","),
		"type":  strings.Join(typ, ","),
		"name":  strings.Join(name, ","),
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.allowNoIndices != nil {
		params.Set("allow_no_indices", fmt.Sprintf("%v", *s.allowNoIndices))
	}
	if s.expandWildcards != "" {
		params.Set("expand_wildcards", s.expandWildcards)
	}
	if s.ignoreUnavailable != nil {
		params.Set("ignore_unavailable", fmt.Sprintf("%v", *s.ignoreUnavailable))
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *GetWarmerService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *GetWarmerService) Do() (map[string]interface{}, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest("GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	// Return operation response
	var ret map[string]interface{}
	if err := json.Unmarshal(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
