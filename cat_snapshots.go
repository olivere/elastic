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

	"github.com/olivere/elastic/v7/uritemplates"
)

// CatSnapshotsService returns the list of snapshots.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.13/cat-snapshots.html
// for details.
type CatSnapshotsService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	repository    string // snapshot repository used to limit the request
	masterTimeout string
	columns       []string
	sort          []string // list of columns for sort order
}

// NewCatSnapshotsService creates a new NewCatSnapshotsService.
func NewCatSnapshotsService(client *Client) *CatSnapshotsService {
	return &CatSnapshotsService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *CatSnapshotsService) Pretty(pretty bool) *CatSnapshotsService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *CatSnapshotsService) Human(human bool) *CatSnapshotsService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *CatSnapshotsService) ErrorTrace(errorTrace bool) *CatSnapshotsService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *CatSnapshotsService) FilterPath(filterPath ...string) *CatSnapshotsService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *CatSnapshotsService) Header(name string, value string) *CatSnapshotsService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *CatSnapshotsService) Headers(headers http.Header) *CatSnapshotsService {
	s.headers = headers
	return s
}

// Repository specifies the napshot repository used to limit the request.
func (s *CatSnapshotsService) Repository(repository string) *CatSnapshotsService {
	s.repository = repository
	return s
}

// MasterTimeout is the explicit operation timeout for connection to master node.
func (s *CatSnapshotsService) MasterTimeout(masterTimeout string) *CatSnapshotsService {
	s.masterTimeout = masterTimeout
	return s
}

// Columns to return in the response.
// To get a list of all possible columns to return, run the following command
// in your terminal:
//
// Example:
//   curl 'http://localhost:9200/_cat/snapshots/<repository>?help'
//
// You can use Columns("*") to return all possible columns. That might take
// a little longer than the default set of columns.
func (s *CatSnapshotsService) Columns(columns ...string) *CatSnapshotsService {
	s.columns = columns
	return s
}

// Sort is a list of fields to sort by.
func (s *CatSnapshotsService) Sort(fields ...string) *CatSnapshotsService {
	s.sort = fields
	return s
}

// buildURL builds the URL for the operation.
func (s *CatSnapshotsService) buildURL() (string, url.Values, error) {
	// Build URL
	var (
		path string
		err  error
	)

	if s.repository != "" {
		path, err = uritemplates.Expand("/_cat/snapshots/{repository}", map[string]string{
			"repository": s.repository,
		})
	} else {
		path = "/_cat/snapshots"
	}
	if err != nil {
		return "", url.Values{}, err
	}

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
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	// TODO
	if len(s.columns) > 0 {
		// loop through all columns and apply alias if needed
		for i, column := range s.columns {
			if fullValueRaw, isAliased := catSnapshotsResponseRowAliasesMap[column]; isAliased {
				// alias can be translated to multiple fields,
				// so if translated value contains a comma, than replace the first value
				// and append the others
				if strings.Contains(fullValueRaw, ",") {
					fullValues := strings.Split(fullValueRaw, ",")
					s.columns[i] = fullValues[0]
					s.columns = append(s.columns, fullValues[1:]...)
				} else {
					s.columns[i] = fullValueRaw
				}
			}
		}

		params.Set("h", strings.Join(s.columns, ","))
	}
	if len(s.sort) > 0 {
		params.Set("s", strings.Join(s.sort, ","))
	}
	return path, params, nil
}

// Do executes the operation.
func (s *CatSnapshotsService) Do(ctx context.Context) (CatSnapshotsResponse, error) {
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
	var ret CatSnapshotsResponse
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of a get request.

// CatSnapshotsResponse is the outcome of CatSnapshotsService.Do.
type CatSnapshotsResponse []CatSnapshotsResponseRow

// CatSnapshotssResponseRow specifies the data returned for one index
// of a CatSnapshotsResponse. Notice that not all of these fields might
// be filled; that depends on the number of columns chose in the
// request (see CatSnapshotsService.Columns).
type CatSnapshotsResponseRow struct {
	ID               string `json:"id"`                // ID of the snapshot, such as "snap1".
	Repository       string `json:"repository"`        // Name of the repository, such as "repo1".
	Status           string `json:"status"`            // One of "FAILED", "INCOMPATIBLE", "IN_PROGRESS", "PARTIAL" or "SUCCESS".
	StartEpoch       string `json:"start_epoch"`       // Unix epoch time at which the snapshot process started.
	StartTime        string `json:"start_time"`        // HH:MM:SS time at which the snapshot process started.
	EndEpoch         string `json:"end_epoch"`         // Unix epoch time at which the snapshot process ended.
	EndTime          string `json:"end_time"`          // HH:MM:SS time at which the snapshot process ended.
	Duration         string `json:"duration"`          // Time it took the snapshot process to complete in time units.
	Indices          string `json:"indices"`           // Number of indices in the snapshot.
	SuccessfulShards string `json:"successful_shards"` // Number of successful shards in the snapshot.
	FailedShards     string `json:"failed_shards"`     // Number of failed shards in the snapshot.
	TotalShards      string `json:"total_shards"`      // Total number of shards in the snapshot.
	Reason           string `json:"reason"`            // Reason for any snapshot failures.
}

// catSnapshotsResponseRowAliasesMap holds the global map for columns aliases
// the map is used by CatSnapshotsService.buildURL.
// For backwards compatibility some fields are able to have the same aliases
// that means that one alias can be translated to different columns (from different elastic versions)
// example for understanding: rto -> RefreshTotal, RefreshExternalTotal
var catSnapshotsResponseRowAliasesMap = map[string]string{
	"snapshot": "id",
	"re":       "repository",
	"s":        "status",
	"ste":      "start_epoch",
	"sti":      "start_time",
	"ete":      "end_epoch",
	"eti":      "end_time",
	"dur":      "duration",
	"i":        "indices",
	"ss":       "successful_shards",
	"fs":       "failed_shards",
	"ts":       "total_shards",
	"`r":       "reason",
}
