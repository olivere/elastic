// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// -- Bulk delete request --

// Bulk request to remove document from Elasticsearch.
type BulkDeleteRequest struct {
	BulkableRequest
	index       string
	typ         string
	id          string
	routing     string
	refresh     *bool
	version     int64  // default is MATCH_ANY
	versionType string // default is "internal"

	source []string
}

func NewBulkDeleteRequest() *BulkDeleteRequest {
	return &BulkDeleteRequest{}
}

func (r *BulkDeleteRequest) Index(index string) *BulkDeleteRequest {
	r.index = index
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) Type(typ string) *BulkDeleteRequest {
	r.typ = typ
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) Id(id string) *BulkDeleteRequest {
	r.id = id
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) Routing(routing string) *BulkDeleteRequest {
	r.routing = routing
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) Refresh(refresh bool) *BulkDeleteRequest {
	r.refresh = &refresh
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) Version(version int64) *BulkDeleteRequest {
	r.version = version
	r.source = nil
	return r
}

// VersionType can be "internal" (default), "external", "external_gte",
// "external_gt", or "force".
func (r *BulkDeleteRequest) VersionType(versionType string) *BulkDeleteRequest {
	r.versionType = versionType
	r.source = nil
	return r
}

func (r *BulkDeleteRequest) String() string {
	lines, err := r.Source()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return strings.Join(lines, "\n")
}

func (r *BulkDeleteRequest) Source() ([]string, error) {
	if r.source != nil {
		return r.source, nil
	}
	lines := make([]string, 1)

	source := make(map[string]interface{})
	deleteCommand := make(map[string]interface{})
	if r.index != "" {
		deleteCommand["_index"] = r.index
	}
	if r.typ != "" {
		deleteCommand["_type"] = r.typ
	}
	if r.id != "" {
		deleteCommand["_id"] = r.id
	}
	if r.routing != "" {
		deleteCommand["_routing"] = r.routing
	}
	if r.version > 0 {
		deleteCommand["_version"] = r.version
	}
	if r.versionType != "" {
		deleteCommand["_version_type"] = r.versionType
	}
	if r.refresh != nil {
		deleteCommand["refresh"] = *r.refresh
	}
	source["delete"] = deleteCommand

	body, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	lines[0] = string(body)
	r.source = lines

	return lines, nil
}
