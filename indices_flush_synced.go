// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// IndicesSyncedFlushService performs a normal flush, then adds a generated
// unique marked (sync_id) to all shards.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.4/indices-synced-flush.html
// for details.
type IndicesSyncedFlushService struct {
	client            *Client
	pretty            bool
	index             []string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
}

// NewIndicesSyncedFlushService creates a new IndicesSyncedFlushService.
func NewIndicesSyncedFlushService(client *Client) *IndicesSyncedFlushService {
	return &IndicesSyncedFlushService{
		client: client,
	}
}

// Index is a list of index names; use `_all` or empty string for all indices.
func (s *IndicesSyncedFlushService) Index(indices ...string) *IndicesSyncedFlushService {
	s.index = append(s.index, indices...)
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *IndicesSyncedFlushService) IgnoreUnavailable(ignoreUnavailable bool) *IndicesSyncedFlushService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices expression
// resolves into no concrete indices. (This includes `_all` string or when
// no indices have been specified).
func (s *IndicesSyncedFlushService) AllowNoIndices(allowNoIndices bool) *IndicesSyncedFlushService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards specifies whether to expand wildcard expression to
// concrete indices that are open, closed or both..
func (s *IndicesSyncedFlushService) ExpandWildcards(expandWildcards string) *IndicesSyncedFlushService {
	s.expandWildcards = expandWildcards
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *IndicesSyncedFlushService) Pretty(pretty bool) *IndicesSyncedFlushService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *IndicesSyncedFlushService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string

	if len(s.index) > 0 {
		path, err = uritemplates.Expand("/{index}/_flush/synced", map[string]string{
			"index": strings.Join(s.index, ","),
		})
	} else {
		path = "/_flush/synced"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
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
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *IndicesSyncedFlushService) Validate() error {
	return nil
}

// Do executes the service.
func (s *IndicesSyncedFlushService) Do(ctx context.Context) (*IndicesSynedFlushResponse, error) {
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
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "POST",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(IndicesSynedFlushResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a flush request.

// IndicesSynedFlushResponse is the outcome of a synched flush call.
type IndicesSynedFlushResponse struct {
	Shards shardsInfo `json:"_shards"`

	// TODO Add information about the indices here from the root level
	// It looks like this:
	// {
	// 	"_shards" : {
	// 	  "total" : 4,
	// 	  "successful" : 4,
	// 	  "failed" : 0
	// 	},
	// 	"elastic-test" : {
	// 	  "total" : 1,
	// 	  "successful" : 1,
	// 	  "failed" : 0
	// 	},
	// 	"elastic-test2" : {
	// 	  "total" : 1,
	// 	  "successful" : 1,
	// 	  "failed" : 0
	// 	},
	// 	"elastic-orders" : {
	// 	  "total" : 1,
	// 	  "successful" : 1,
	// 	  "failed" : 0
	// 	},
	// 	"elastic-nosource-test" : {
	// 	  "total" : 1,
	// 	  "successful" : 1,
	// 	  "failed" : 0
	// 	}
	// }
}
