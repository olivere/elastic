// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/url"
	"strings"
)

type AliasAction struct {
	// "add" or "remove"
	actionType string
	// Index name
	index string
	// Alias name
	alias string

	// Below only apply to "add" actions

	// Filter
	filter Filter
	// Routing value
	routing string
	// Search routing value
	searchRouting string
	// Index routing value
	indexRouting string
}

func NewAliasAddAction(index, alias string) AliasAction {
	return AliasAction{
		actionType: "add",
		index:      index,
		alias:      alias,
	}
}

func (a AliasAction) Filter(filter Filter) AliasAction {
	a.filter = filter
	return a
}

func (a AliasAction) Routing(routing string) AliasAction {
	a.routing = routing
	return a
}

func (a AliasAction) IndexRouting(routing string) AliasAction {
	a.indexRouting = routing
	return a
}

func (a AliasAction) SearchRouting(routings ...string) AliasAction {
	a.searchRouting = strings.Join(routings, ",")
	return a
}

func NewAliasRemoveAction(index, alias string) AliasAction {
	return AliasAction{
		actionType: "remove",
		index:      index,
		alias:      alias,
	}
}

type AliasService struct {
	client  *Client
	actions []AliasAction
	pretty  bool
}

func NewAliasService(client *Client) *AliasService {
	builder := &AliasService{
		client:  client,
		actions: make([]AliasAction, 0),
	}
	return builder
}

func (s *AliasService) Pretty(pretty bool) *AliasService {
	s.pretty = pretty
	return s
}

func (s *AliasService) Add(indexName string, aliasName string) *AliasService {
	action := AliasAction{actionType: "add", index: indexName, alias: aliasName}
	s.actions = append(s.actions, action)
	return s
}

func (s *AliasService) AddWithFilter(indexName string, aliasName string, filter *Filter) *AliasService {
	var f Filter
	if filter != nil {
		f = *filter
	}
	action := AliasAction{actionType: "add", index: indexName, alias: aliasName, filter: f}
	s.actions = append(s.actions, action)
	return s
}

func (s *AliasService) Remove(indexName string, aliasName string) *AliasService {
	action := AliasAction{actionType: "remove", index: indexName, alias: aliasName}
	s.actions = append(s.actions, action)
	return s
}

func (s *AliasService) Actions(actions ...AliasAction) *AliasService {
	s.actions = append(actions)
	return s
}

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
	actionsJson := make([]interface{}, 0)

	for _, action := range s.actions {
		actionJson := make(map[string]interface{})
		detailsJson := make(map[string]interface{})
		detailsJson["index"] = action.index
		detailsJson["alias"] = action.alias
		if action.filter != nil {
			detailsJson["filter"] = action.filter.Source()
		}
		if action.routing != "" {
			detailsJson["routing"] = action.routing
		}
		if action.indexRouting != "" {
			detailsJson["index_routing"] = action.indexRouting
		}
		if action.searchRouting != "" {
			detailsJson["search_routing"] = action.searchRouting
		}
		actionJson[action.actionType] = detailsJson
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

type AliasResult struct {
	Acknowledged bool `json:"acknowledged"`
}
