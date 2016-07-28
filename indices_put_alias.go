// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/url"
)

// AliasService manages index aliases.
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/indices-aliases.html.
type AliasService struct {
	client  *Client
	actions []aliasAction
	pretty  bool
}

// aliasAction is a single action applied to an alias: add or remove.
type aliasAction struct {
	// "add" or "remove"
	Type string
	// Index name
	Index string
	// Alias name
	Alias string
	// Filter
	Filter Query
}

// NewAliasService creates a new instance of AliasService.
func NewAliasService(client *Client) *AliasService {
	builder := &AliasService{
		client: client,
	}
	return builder
}

// Pretty asks Elasticsearch to return indented JSON.
func (s *AliasService) Pretty(pretty bool) *AliasService {
	s.pretty = pretty
	return s
}

// Add an alias.
func (s *AliasService) Add(indexName string, aliasName string) *AliasService {
	action := aliasAction{Type: "add", Index: indexName, Alias: aliasName}
	s.actions = append(s.actions, action)
	return s
}

// AddWithFilter adds an alias with a filter.
func (s *AliasService) AddWithFilter(indexName string, aliasName string, filter Query) *AliasService {
	action := aliasAction{Type: "add", Index: indexName, Alias: aliasName, Filter: filter}
	s.actions = append(s.actions, action)
	return s
}

// Remove removes an alias.
func (s *AliasService) Remove(indexName string, aliasName string) *AliasService {
	action := aliasAction{Type: "remove", Index: indexName, Alias: aliasName}
	s.actions = append(s.actions, action)
	return s
}

// Do executes the request.
func (s *AliasService) Do() (*AliasResult, error) {
	// Build url
	path := "/_aliases"

	// Parameters
	params := make(url.Values)
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}

	// Actions
	body := make(map[string]interface{})
	var actionsJson []interface{}

	for _, action := range s.actions {
		actionJson := make(map[string]interface{})
		detailsJson := make(map[string]interface{})
		detailsJson["index"] = action.Index
		detailsJson["alias"] = action.Alias
		if action.Filter != nil {
			src, err := action.Filter.Source()
			if err != nil {
				return nil, err
			}
			detailsJson["filter"] = src
		}
		actionJson[action.Type] = detailsJson
		actionsJson = append(actionsJson, actionJson)
	}

	body["actions"] = actionsJson

	// Get response
	res, err := s.client.PerformRequest("POST", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return results
	ret := new(AliasResult)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of an alias request.

// AliasResult is the outcome of AliasService.Do.
type AliasResult struct {
	Acknowledged bool `json:"acknowledged"`
}
