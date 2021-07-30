// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CatMasterService shows information about the master node,
// including the ID, bound IP address, and name.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.12/cat-master.html
// for details.
type CatMasterService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	local         *bool
	masterTimeout string
	columns       []string
	sort          []string // list of columns for sort order
}

// NewCatMasterService creates a new CatMasterService
func NewCatMasterService(client *Client) *CatMasterService {
	return &CatMasterService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *CatMasterService) Pretty(pretty bool) *CatMasterService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *CatMasterService) Human(human bool) *CatMasterService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *CatMasterService) ErrorTrace(errorTrace bool) *CatMasterService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *CatMasterService) FilterPath(filterPath ...string) *CatMasterService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *CatMasterService) Header(name string, value string) *CatMasterService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *CatMasterService) Headers(headers http.Header) *CatMasterService {
	s.headers = headers
	return s
}

// Local indicates to return local information, i.e. do not retrieve
// the state from master node (default: false).
func (s *CatMasterService) Local(local bool) *CatMasterService {
	s.local = &local
	return s
}

// MasterTimeout is the explicit operation timeout for connection to master node.
func (s *CatMasterService) MasterTimeout(masterTimeout string) *CatMasterService {
	s.masterTimeout = masterTimeout
	return s
}

// Columns to return in the response.
// To get a list of all possible columns to return, run the following command
// in your terminal:
//
// Example:
//   curl 'http://localhost:9200/_cat/master?help'
//
// You can use Columns("*") to return all possible columns. That might take
// a little longer than the default set of columns.
func (s *CatMasterService) Columns(columns ...string) *CatMasterService {
	s.columns = columns
	return s
}

// Sort is a list of fields to sort by.
func (s *CatMasterService) Sort(fields ...string) *CatMasterService {
	s.sort = fields
	return s
}

// buildURL builds the URL for the operation.
func (s *CatMasterService) buildURL() (string, url.Values, error) {
	// Build URL
	path := "/_cat/master"

	// Add query string parameters
	params := url.Values{
		"format": []string{"json"}, // always returns as JSON
	}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	if v := s.local; v != nil {
		params.Set("local", fmt.Sprint(*v))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if len(s.sort) > 0 {
		params.Set("s", strings.Join(s.sort, ","))
	}
	if len(s.columns) > 0 {
		params.Set("h", strings.Join(s.columns, ","))
	}
	return path, params, nil
}

// Do executes the operation.
func (s *CatMasterService) Do(ctx context.Context) (CatMasterResponse, error) {
	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method:  "GET",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	var ret CatMasterResponse
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a get request.

// CatMasterResponse is the outcome of CatMasterService.Do.
type CatMasterResponse []CatMasterResponseRow

// CatMasterResponseRow is a single row in a CatMasterResponse.
// Notice that not all of these fields might be filled; that depends
// on the number of columns chose in the request (see CatMasterService.Columns).
type CatMasterResponseRow struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	IP   string `json:"ip"`
	Node string `json:"node"`
}
