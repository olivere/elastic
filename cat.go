// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatService executes cat commands against the cluster.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.2/cat.html
// for details.
type CatService struct {
	client *Client
	pretty bool
	model  string
}

const (
	aliases      string = "aliases"
	allocation   string = "allocation"
	count        string = "count"
	fielddata    string = "fielddata"
	health       string = "health"
	indices      string = "indices"
	master       string = "master"
	nodeattrs    string = "nodeattrs"
	nodes        string = "nodes"
	pendingtasks string = "pending_tasks"
	plugins      string = "plugins"
	recovery     string = "recovery"
	repositories string = "repositories"
	shards       string = "shards"
	segments     string = "segments"
	templates    string = "templates"
	threadpool   string = "thread_pool"
)

var catModels = map[string]struct{}{
	aliases:      struct{}{},
	allocation:   struct{}{},
	count:        struct{}{},
	fielddata:    struct{}{},
	health:       struct{}{},
	indices:      struct{}{},
	master:       struct{}{},
	nodeattrs:    struct{}{},
	nodes:        struct{}{},
	pendingtasks: struct{}{},
	plugins:      struct{}{},
	recovery:     struct{}{},
	repositories: struct{}{},
	shards:       struct{}{},
	segments:     struct{}{},
	templates:    struct{}{},
	threadpool:   struct{}{},
}

// NewCatService creates a new CatService.
func NewCatService(client *Client) *CatService {
	return &CatService{
		client: client,
		pretty: true,
	}
}

// Aliases sets the service to show information about currently configured aliases to indices including filter and routing infos.
func (s *CatService) Aliases() *CatService {
	s.model = aliases
	return s
}

// Allocation sets the service to show a snapshot of how many shards are allocated to each data node and how much disk space they are using.
func (s *CatService) Allocation() *CatService {
	s.model = allocation
	return s
}

// Count sets the service to show the document count of the entire cluster, or individual indices.
func (s *CatService) Count() *CatService {
	s.model = count
	return s
}

// FieldData sets the service to show how much heap memory is currently being used by fielddata on every data node in the cluster.
func (s *CatService) FieldData() *CatService {
	s.model = fielddata
	return s
}

// Health sets the service to show a terse, one-line representation of the same information from /_cluster/health
func (s *CatService) Health() *CatService {
	s.model = health
	return s
}

// Indices sets the service to show a cross-section of each index.
func (s *CatService) Indices() *CatService {
	s.model = indices
	return s
}

// Master sets the service to show the master’s node ID, bound IP address, and node name.
func (s *CatService) Master() *CatService {
	s.model = master
	return s
}

// NodeAttrs sets the service to show custom node attributes.
func (s *CatService) NodeAttrs() *CatService {
	s.model = nodeattrs
	return s
}

// Nodes sets the service to show the cluster topology.
func (s *CatService) Nodes() *CatService {
	s.model = nodes
	return s
}

// PendingTasks sets the service to show data about any pending updates to the cluster state.
func (s *CatService) PendingTasks() *CatService {
	s.model = pendingtasks
	return s
}

// Plugins sets the service to show a view per node of running plugins.
func (s *CatService) Plugins() *CatService {
	s.model = plugins
	return s
}

// Recovery sets the service to show a view of index shard recoveries, both on-going and previously completed.
func (s *CatService) Recovery() *CatService {
	s.model = recovery
	return s
}

// Repositories sets the service to show the snapshot repositories registered in the cluster.
func (s *CatService) Repositories() *CatService {
	s.model = repositories
	return s
}

// Segments sets the service to show low level information about the segments in the shards of an index.
func (s *CatService) Segments() *CatService {
	s.model = segments
	return s
}

// Shards sets the service to show the detailed view of what nodes contain which shards.
func (s *CatService) Shards() *CatService {
	s.model = shards
	return s
}

// Templates sets the service to show information about existing templates.
func (s *CatService) Templates() *CatService {
	s.model = templates
	return s
}

// ThreadPool sets the service to show cluster wide thread pool statistics per node.
func (s *CatService) ThreadPool() *CatService {
	s.model = threadpool
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatService) Pretty(pretty bool) *CatService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *CatService) buildURL() (string, url.Values, error) {
	path := fmt.Sprintf("/_cat/%s", s.model)

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *CatService) Validate() error {
	if s.model == "" {
		return fmt.Errorf("no model specified for cat command")
	}

	if _, ok := catModels[s.model]; !ok {
		return fmt.Errorf("unknown model '%s' specified for cat command", s.model)
	}

	return nil
}

// Do executes the operation.
func (s *CatService) Do(ctx context.Context) (string, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return "", err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return "", err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "GET",
		Path:   path,
		Params: params,
	})

	if err != nil {
		return "", err
	}

	return string(res.Body), nil
}

// CatResponse is the response of CatService.Do.
type CatResponse struct {
	Response string
}
