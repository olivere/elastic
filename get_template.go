// Copyright 2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// GetTemplateService reads a search template.
// It is documented at http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-template.html.
type GetTemplateService struct {
	client      *Client
	debug       bool
	pretty      bool
	id          string
	version     interface{}
	versionType string
}

// NewGetTemplateService creates a new GetTemplateService.
func NewGetTemplateService(client *Client) *GetTemplateService {
	return &GetTemplateService{
		client: client,
	}
}

// Id is documented as: Template ID.
func (s *GetTemplateService) Id(id string) *GetTemplateService {
	s.id = id
	return s
}

// Version is documented as: Explicit version number for concurrency control.
func (s *GetTemplateService) Version(version interface{}) *GetTemplateService {
	s.version = version
	return s
}

// VersionType is documented as: Specific version type.
func (s *GetTemplateService) VersionType(versionType string) *GetTemplateService {
	s.versionType = versionType
	return s
}

// buildURL builds the URL for the operation.
func (s *GetTemplateService) buildURL() (string, error) {
	// Build URL
	urls, err := uritemplates.Expand("/_search/template/{id}", map[string]string{
		"id": s.id,
	})
	if err != nil {
		return "", err
	}

	// Add query string parameters
	params := url.Values{}
	if s.version != nil {
		params.Set("version", fmt.Sprintf("%v", s.version))
	}
	if s.versionType != "" {
		params.Set("version_type", s.versionType)
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	return urls, nil
}

// Validate checks if the operation is valid.
func (s *GetTemplateService) Validate() error {
	var invalid []string
	if s.id == "" {
		invalid = append(invalid, "Id")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation and returns the template.
func (s *GetTemplateService) Do() (*GetTemplateResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	urls, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request
	req, err := s.client.NewRequest("GET", urls)
	if err != nil {
		return nil, err
	}

	// Debug output?
	if s.debug {
		s.client.dumpRequest((*http.Request)(req))
	}

	// Get HTTP response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Debug output?
	if s.debug {
		s.client.dumpResponse(res)
	}

	// Decode response
	resp := new(GetTemplateResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type GetTemplateResponse struct {
	Template string `json:"template"`
}
