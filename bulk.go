// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type BulkService struct {
	client *Client

	index    string
	_type    string
	requests []BulkableRequest
	//replicationType string
	//consistencyLevel string
	refresh  *bool
	debug    bool
}

func NewBulkService(client *Client) *BulkService {
	builder := &BulkService{
		client: client,
		requests: make([]BulkableRequest, 0),
		debug: false,
	}
	return builder
}

func (s *BulkService) Index(index string) (*BulkService) {
	s.index = index
	return s
}

func (s *BulkService) Type(_type string) (*BulkService) {
	s._type = _type
	return s
}

func (s *BulkService) Debug(debug bool) (*BulkService) {
	s.debug = debug
	return s
}

func (s *BulkService) Add(r BulkableRequest) (*BulkService) {
	s.requests = append(s.requests, r)
	return s
}

func (s *BulkService) NumberOfActions() int {
	return len(s.requests)
}

func (s *BulkService) bodyAsString() (string, error) {
	buf := bytes.NewBufferString("")

	for _, req := range s.requests {
		source, err := req.Source()
		if err != nil {
			return "", err
		}
		for _, line := range source {
			_, err := buf.WriteString(fmt.Sprintf("%s\n", line))
			if err != nil {
				return "", nil
			}
		}
	}

	return buf.String(), nil
}

func (s *BulkService) Do() (*BulkResponse, error) {
	// Get body
	body, err := s.bodyAsString()
	if err != nil {
		return nil, err
	}

	// Build url
	urls := "/"
	if s.index != "" {
		urls += cleanPathString(s.index) + "/"
	}
	if s._type != "" {
		urls += cleanPathString(s._type) + "/"
	}
	urls += "_bulk"

	// TODO params

	// Set up a new request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	// Set body
	req.SetBodyString(body)

	// Debug
	if s.debug {
		out, _ := httputil.DumpRequestOut((*http.Request)(req), true)
		fmt.Printf("%s\n", string(out))
	}

	// Get response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Debug	
	if s.debug {
		out, _ := httputil.DumpResponse(res, true)
		fmt.Printf("%s\n", string(out))
	}

	ret := new(BulkResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}


// Response to bulk execution.
type BulkResponse struct {
	Took int `json:"took"`
}



// Generic interface to bulkable requests.
type BulkableRequest interface {
	Source() ([]string, error)
}



// Bulk request to add document to ElasticSearch.
type BulkIndexRequest struct {
	BulkableRequest
	Index      string
	Type       string
	Id         string
	Data       interface{}
}

func NewBulkIndexRequest(index,_type,id string, data interface{}) *BulkIndexRequest {
	return &BulkIndexRequest{
		Index: index,
		Type: _type,
		Id: id,
		Data: data,
	}
}

func (r BulkIndexRequest) Source() ([]string, error) {
	// { "index" : { "_index" : "test", "_type" : "type1", "_id" : "1" } }
	// { "field1" : "value1" }

	lines := make([]string, 2)

	// "index" ...
	command := make(map[string]interface{})
	indexCommand := make(map[string]interface{})
	command["index"] = indexCommand
	if r.Index != "" {
		indexCommand["_index"] = r.Index
	}
	if r.Type != "" {
		indexCommand["_type"] = r.Type
	}
	if r.Id != "" {
		indexCommand["_id"] = r.Id
	}
	// TODO _version
	// TODO _version_type
	// TODO _routing
	// TODO _percolate
	// TODO _parent
	// TODO _timestamp
	// TODO _ttl
	line, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	lines[0] = string(line)

	// "field1" ...
	if r.Data != nil {
		body, err := json.Marshal(r.Data)
		if err != nil {
			return nil, err
		}
		lines[1] = string(body)
	}

	return lines, nil
}



// Bulk request to remove document from ElasticSearch.
type BulkDeleteRequest struct {
	BulkableRequest
	Index  string
	Type   string
	Id     string
}

func NewBulkDeleteRequest(index,_type,id string) *BulkDeleteRequest {
	return &BulkDeleteRequest{
		Index: index,
		Type: _type,
		Id: id,
	}
}

func (r BulkDeleteRequest) Source() ([]string, error) {
	lines := make([]string, 1)

	source := make(map[string]interface{})
	data := make(map[string]interface{})
	source["delete"] = data

	if r.Index != "" {
		data["_index"] = r.Index
	}
	if r.Type != "" {
		data["_type"] = r.Type
	}
	if r.Id != "" {
		data["_id"] = r.Id
	}
	// TODO _version
	// TODO _version_type
	// TODO _routing
	// TODO _percolate
	// TODO _parent
	// TODO _timestamp
	// TODO _ttl

	body, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	lines[0] = string(body)
	
	return lines, nil
}