// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// PingService checks if an Elasticsearch server on a given URL is alive.
// When asked for, it can also return various information about the
// Elasticsearch server, e.g. the Elasticsearch version number.
//
// Ping simply starts a HTTP GET request to the URL of the server.
// If the server responds with HTTP Status code 200 OK, the server is alive.
type PingService struct {
	client       *Client
	url          string
	timeout      string
	httpHeadOnly bool
	pretty       bool
}

// PingResult is the result returned from querying the Elasticsearch server.
type PingResult struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Version     struct {
		Number                           string `json:"number"`                              // e.g. "7.0.0"
		BuildFlavor                      string `json:"build_flavor"`                        // e.g. "oss" or "default"
		BuildType                        string `json:"build_type"`                          // e.g. "docker"
		BuildHash                        string `json:"build_hash"`                          // e.g. "b7e28a7"
		BuildDate                        string `json:"build_date"`                          // e.g. "2019-04-05T22:55:32.697037Z"
		BuildSnapshot                    bool   `json:"build_snapshot"`                      // e.g. false
		LuceneVersion                    string `json:"lucene_version"`                      // e.g. "8.0.0"
		MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`  // e.g. "6.7.0"
		MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"` // e.g. "6.0.0-beta1"
	} `json:"version"`
	TagLine string `json:"tagline"`
}

func NewPingService(client *Client) *PingService {
	return &PingService{
		client:       client,
		url:          DefaultURL,
		httpHeadOnly: false,
		pretty:       false,
	}
}

func (s *PingService) URL(url string) *PingService {
	s.url = url
	return s
}

func (s *PingService) Timeout(timeout string) *PingService {
	s.timeout = timeout
	return s
}

// HeadOnly makes the service to only return the status code in Do;
// the PingResult will be nil.
func (s *PingService) HttpHeadOnly(httpHeadOnly bool) *PingService {
	s.httpHeadOnly = httpHeadOnly
	return s
}

func (s *PingService) Pretty(pretty bool) *PingService {
	s.pretty = pretty
	return s
}

// Do returns the PingResult, the HTTP status code of the Elasticsearch
// server, and an error.
func (s *PingService) Do(ctx context.Context) (*PingResult, int, error) {
	s.client.mu.RLock()
	basicAuth := s.client.basicAuth
	basicAuthUsername := s.client.basicAuthUsername
	basicAuthPassword := s.client.basicAuthPassword
	s.client.mu.RUnlock()

	url_ := s.url + "/"

	params := make(url.Values)
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if len(params) > 0 {
		url_ += "?" + params.Encode()
	}

	var method string
	if s.httpHeadOnly {
		method = "HEAD"
	} else {
		method = "GET"
	}

	// Notice: This service must NOT use PerformRequest!
	req, err := NewRequest(method, url_)
	if err != nil {
		return nil, 0, err
	}

	if basicAuth {
		req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
	}

	res, err := s.client.c.Do((*http.Request)(req).WithContext(ctx))
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var ret *PingResult
	if !s.httpHeadOnly {
		ret = new(PingResult)
		if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
			return nil, res.StatusCode, err
		}
	}

	return ret, res.StatusCode, nil
}
