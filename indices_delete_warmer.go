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

// DeleteWarmerService allows to delete a warmer.
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-warmers.html.
type DeleteWarmerService struct {
	client        *Client
	pretty        bool
	index         []string
	name          []string
	masterTimeout string
}

// NewDeleteWarmerService creates a new DeleteWarmerService.
func NewDeleteWarmerService(client *Client) *DeleteWarmerService {
	return &DeleteWarmerService{
		client: client,
		index:  make([]string, 0),
		name:   make([]string, 0),
	}
}

// Index is a list of index names the mapping should be added to
// (supports wildcards); use `_all` or omit to add the mapping on all indices.
func (s *DeleteWarmerService) Index(indices ...string) *DeleteWarmerService {
	s.index = append(s.index, indices...)
	return s
}

// Name is a list of warmer names to delete (supports wildcards);
// use `_all` to delete all warmers in the specified indices.
func (s *DeleteWarmerService) Name(name ...string) *DeleteWarmerService {
	s.name = append(s.name, name...)
	return s
}

// MasterTimeout specifies the timeout for connection to master.
func (s *DeleteWarmerService) MasterTimeout(masterTimeout string) *DeleteWarmerService {
	s.masterTimeout = masterTimeout
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *DeleteWarmerService) Pretty(pretty bool) *DeleteWarmerService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *DeleteWarmerService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/{index}/_warmer/{name}", map[string]string{
		"index": strings.Join(s.index, ","),
		"name":  strings.Join(s.name, ","),
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if len(s.name) > 0 {
		params.Set("name", strings.Join(s.name, ","))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *DeleteWarmerService) Validate() error {
	var invalid []string
	if len(s.index) == 0 {
		invalid = append(invalid, "Index")
	}
	if len(s.name) == 0 {
		invalid = append(invalid, "Name")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *DeleteWarmerService) Do() (*DeleteWarmerResponse, error) {
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
	res, err := s.client.PerformRequest("DELETE", path, params, nil)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(DeleteWarmerResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteWarmerResponse is the response of DeleteWarmerService.Do.
type DeleteWarmerResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
