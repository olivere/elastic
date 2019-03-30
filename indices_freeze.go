// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/v7/uritemplates"
)

// IndicesFreezeService freezes an index.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/freeze-index-api.html
// and https://www.elastic.co/blog/creating-frozen-indices-with-the-elasticsearch-freeze-index-api
// for details.
type IndicesFreezeService struct {
	client              *Client
	pretty              bool
	index               string
	timeout             string
	masterTimeout       string
	ignoreUnavailable   *bool
	allowNoIndices      *bool
	expandWildcards     string
	waitForActiveShards string
}

// NewIndicesFreezeService creates a new IndicesFreezeService.
func NewIndicesFreezeService(client *Client) *IndicesFreezeService {
	return &IndicesFreezeService{
		client: client,
	}
}

// Index is the name of the index to freeze.
func (s *IndicesFreezeService) Index(index string) *IndicesFreezeService {
	s.index = index
	return s
}

// Timeout allows to specify an explicit timeout.
func (s *IndicesFreezeService) Timeout(timeout string) *IndicesFreezeService {
	s.timeout = timeout
	return s
}

// MasterTimeout allows to specify a timeout for connection to master.
func (s *IndicesFreezeService) MasterTimeout(masterTimeout string) *IndicesFreezeService {
	s.masterTimeout = masterTimeout
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *IndicesFreezeService) IgnoreUnavailable(ignoreUnavailable bool) *IndicesFreezeService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices expression
// resolves into no concrete indices. (This includes `_all` string or when
// no indices have been specified).
func (s *IndicesFreezeService) AllowNoIndices(allowNoIndices bool) *IndicesFreezeService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards specifies whether to expand wildcard expression to
// concrete indices that are open, closed or both..
func (s *IndicesFreezeService) ExpandWildcards(expandWildcards string) *IndicesFreezeService {
	s.expandWildcards = expandWildcards
	return s
}

// WaitForActiveShards sets the number of active shards to wait for
// before the operation returns.
func (s *IndicesFreezeService) WaitForActiveShards(numShards string) *IndicesFreezeService {
	s.waitForActiveShards = numShards
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *IndicesFreezeService) Pretty(pretty bool) *IndicesFreezeService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *IndicesFreezeService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/{index}/_freeze", map[string]string{
		"index": s.index,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if s.expandWildcards != "" {
		params.Set("expand_wildcards", s.expandWildcards)
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
	if s.waitForActiveShards != "" {
		params.Set("wait_for_active_shards", s.waitForActiveShards)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *IndicesFreezeService) Validate() error {
	var invalid []string
	if s.index == "" {
		invalid = append(invalid, "Index")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the service.
func (s *IndicesFreezeService) Do(ctx context.Context) (*IndicesFreezeResponse, error) {
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
	ret := new(IndicesFreezeResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// IndicesFreezeResponse is the outcome of freezing an index.
type IndicesFreezeResponse struct {
	Shards *ShardsInfo `json:"_shards"`
}
