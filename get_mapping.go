// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/olivere/elastic.v2/uritemplates"
)

// GetMappingService retrieves the mapping definitions for an index or
// index/type.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/1.7/indices-get-mapping.html
// for details.
type GetMappingService struct {
	client            *Client
	pretty            bool
	index             []string
	typ               []string
	local             *bool
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
}

// NewGetMappingService creates a new GetMappingService.
func NewGetMappingService(client *Client) *GetMappingService {
	return &GetMappingService{
		client: client,
		index:  make([]string, 0),
		typ:    make([]string, 0),
	}
}

// Index is a list of index names.
func (s *GetMappingService) Index(index ...string) *GetMappingService {
	s.index = append(s.index, index...)
	return s
}

// Type is a list of document types.
func (s *GetMappingService) Type(typ ...string) *GetMappingService {
	s.typ = append(s.typ, typ...)
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices.
// This includes `_all` string or when no indices have been specified.
func (s *GetMappingService) AllowNoIndices(allowNoIndices bool) *GetMappingService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both..
func (s *GetMappingService) ExpandWildcards(expandWildcards string) *GetMappingService {
	s.expandWildcards = expandWildcards
	return s
}

// Local indicates whether to return local information, do not retrieve
// the state from master node (default: false).
func (s *GetMappingService) Local(local bool) *GetMappingService {
	s.local = &local
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *GetMappingService) IgnoreUnavailable(ignoreUnavailable bool) *GetMappingService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *GetMappingService) Pretty(pretty bool) *GetMappingService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *GetMappingService) buildURL() (string, url.Values, error) {
	var index, typ []string

	if len(s.index) > 0 {
		index = s.index
	} else {
		index = []string{"_all"}
	}

	if len(s.typ) > 0 {
		typ = s.typ
	} else {
		typ = []string{"_all"}
	}

	// Build URL
	path, err := uritemplates.Expand("/{index}/_mapping/{type}", map[string]string{
		"index": strings.Join(index, ","),
		"type":  strings.Join(typ, ","),
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.ignoreUnavailable != nil {
		params.Set("ignore_unavailable", fmt.Sprintf("%v", *s.ignoreUnavailable))
	}
	if s.allowNoIndices != nil {
		params.Set("allow_no_indices", fmt.Sprintf("%v", *s.allowNoIndices))
	}
	if s.expandWildcards != "" {
		params.Set("expand_wildcards", s.expandWildcards)
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *GetMappingService) Validate() error {
	return nil
}

// Do runs DoC() with default context.
func (s *GetMappingService) Do() (map[string]interface{}, error) {
	return s.DoC(nil)
}

// DoC executes the operation.
func (s *GetMappingService) DoC(ctx context.Context) (map[string]interface{}, error) {
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
	res, err := s.client.PerformRequestC(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	// Return operation response
	var ret map[string]interface{}
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
