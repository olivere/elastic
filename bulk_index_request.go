// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Bulk request to add document to Elasticsearch.
type BulkIndexRequest struct {
	BulkableRequest
	index       string
	typ         string
	id          string
	opType      string
	routing     string
	parent      string
	timestamp   string
	ttl         int64
	refresh     *bool
	version     int64  // default is MATCH_ANY
	versionType string // default is "internal"
	doc         interface{}

	source []string
}

func NewBulkIndexRequest() *BulkIndexRequest {
	return &BulkIndexRequest{
		opType: "index",
	}
}

func (r *BulkIndexRequest) Index(index string) *BulkIndexRequest {
	r.index = index
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Type(typ string) *BulkIndexRequest {
	r.typ = typ
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Id(id string) *BulkIndexRequest {
	r.id = id
	r.source = nil
	return r
}

func (r *BulkIndexRequest) OpType(opType string) *BulkIndexRequest {
	r.opType = opType
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Routing(routing string) *BulkIndexRequest {
	r.routing = routing
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Parent(parent string) *BulkIndexRequest {
	r.parent = parent
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Timestamp(timestamp string) *BulkIndexRequest {
	r.timestamp = timestamp
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Ttl(ttl int64) *BulkIndexRequest {
	r.ttl = ttl
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Refresh(refresh bool) *BulkIndexRequest {
	r.refresh = &refresh
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Version(version int64) *BulkIndexRequest {
	r.version = version
	r.source = nil
	return r
}

func (r *BulkIndexRequest) VersionType(versionType string) *BulkIndexRequest {
	r.versionType = versionType
	r.source = nil
	return r
}

func (r *BulkIndexRequest) Doc(doc interface{}) *BulkIndexRequest {
	r.doc = doc
	r.source = nil
	return r
}

func (r *BulkIndexRequest) String() string {
	lines, err := r.Source()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return strings.Join(lines, "\n")
}

func (r *BulkIndexRequest) Source() ([]string, error) {
	// { "index" : { "_index" : "test", "_type" : "type1", "_id" : "1" } }
	// { "field1" : "value1" }

	if r.source != nil {
		return r.source, nil
	}

	lines := make([]string, 2)

	// "index" ...
	command := make(map[string]interface{})
	indexCommand := make(map[string]interface{})
	if r.index != "" {
		indexCommand["_index"] = r.index
	}
	if r.typ != "" {
		indexCommand["_type"] = r.typ
	}
	if r.id != "" {
		indexCommand["_id"] = r.id
	}
	if r.routing != "" {
		indexCommand["_routing"] = r.routing
	}
	if r.parent != "" {
		indexCommand["_parent"] = r.parent
	}
	if r.timestamp != "" {
		indexCommand["_timestamp"] = r.timestamp
	}
	if r.ttl > 0 {
		indexCommand["_ttl"] = r.ttl
	}
	if r.version > 0 {
		indexCommand["_version"] = r.version
	}
	if r.versionType != "" {
		indexCommand["_version_type"] = r.versionType
	}
	if r.refresh != nil {
		indexCommand["refresh"] = *r.refresh
	}
	command[r.opType] = indexCommand
	line, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	lines[0] = string(line)

	// "field1" ...
	if r.doc != nil {
		switch t := r.doc.(type) {
		default:
			body, err := json.Marshal(r.doc)
			if err != nil {
				return nil, err
			}
			lines[1] = string(body)
		case json.RawMessage:
			lines[1] = string(t)
		case *json.RawMessage:
			lines[1] = string(*t)
		case string:
			lines[1] = t
		case *string:
			lines[1] = *t
		}
	} else {
		lines[1] = "{}"
	}

	r.source = lines
	return lines, nil
}
