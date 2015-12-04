// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/olivere/elastic.v2/uritemplates"
)

// TermvectorService returns information and statistics on terms in the
// fields of a particular document. The document could be stored in the
// index or artificially provided by the user.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/1.7/docs-termvectors.html
// for documentation.
type TermvectorService struct {
	client           *Client
	pretty           bool
	index            string
	typ              string
	id               string
	doc              interface{}
	fieldStatistics  *bool
	fields           []string
	perFieldAnalyzer map[string]string
	offsets          *bool
	parent           string
	payloads         *bool
	positions        *bool
	preference       string
	realtime         *bool
	routing          string
	termStatistics   *bool
	bodyJson         interface{}
	bodyString       string
}

// NewTermvectorService creates a new TermvectorService.
func NewTermvectorService(client *Client) *TermvectorService {
	return &TermvectorService{
		client: client,
	}
}

// Index in which the document resides.
func (s *TermvectorService) Index(index string) *TermvectorService {
	s.index = index
	return s
}

// Type of the document.
func (s *TermvectorService) Type(typ string) *TermvectorService {
	s.typ = typ
	return s
}

// Id of the document.
func (s *TermvectorService) Id(id string) *TermvectorService {
	s.id = id
	return s
}

// Doc is the document to analyze.
func (s *TermvectorService) Doc(doc interface{}) *TermvectorService {
	s.doc = doc
	return s
}

// FieldStatistics specifies if document count, sum of document frequencies
// and sum of total term frequencies should be returned.
func (s *TermvectorService) FieldStatistics(fieldStatistics bool) *TermvectorService {
	s.fieldStatistics = &fieldStatistics
	return s
}

// Fields a list of fields to return.
func (s *TermvectorService) Fields(fields ...string) *TermvectorService {
	if s.fields == nil {
		s.fields = make([]string, 0)
	}
	s.fields = append(s.fields, fields...)
	return s
}

// PerFieldAnalyzer allows to specify a different analyzer than the one
// at the field.
func (s *TermvectorService) PerFieldAnalyzer(perFieldAnalyzer map[string]string) *TermvectorService {
	s.perFieldAnalyzer = perFieldAnalyzer
	return s
}

// Offsets specifies if term offsets should be returned.
func (s *TermvectorService) Offsets(offsets bool) *TermvectorService {
	s.offsets = &offsets
	return s
}

// Parent id of documents.
func (s *TermvectorService) Parent(parent string) *TermvectorService {
	s.parent = parent
	return s
}

// Payloads specifies if term payloads should be returned.
func (s *TermvectorService) Payloads(payloads bool) *TermvectorService {
	s.payloads = &payloads
	return s
}

// Positions specifies if term positions should be returned.
func (s *TermvectorService) Positions(positions bool) *TermvectorService {
	s.positions = &positions
	return s
}

// Preference specify the node or shard the operation
// should be performed on (default: random).
func (s *TermvectorService) Preference(preference string) *TermvectorService {
	s.preference = preference
	return s
}

// Realtime specifies if request is real-time as opposed to
// near-real-time (default: true).
func (s *TermvectorService) Realtime(realtime bool) *TermvectorService {
	s.realtime = &realtime
	return s
}

// Routing is a specific routing value.
func (s *TermvectorService) Routing(routing string) *TermvectorService {
	s.routing = routing
	return s
}

// TermStatistics specifies if total term frequency and document frequency
// should be returned.
func (s *TermvectorService) TermStatistics(termStatistics bool) *TermvectorService {
	s.termStatistics = &termStatistics
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *TermvectorService) Pretty(pretty bool) *TermvectorService {
	s.pretty = pretty
	return s
}

// BodyJson defines the body parameters. See documentation.
func (s *TermvectorService) BodyJson(body interface{}) *TermvectorService {
	s.bodyJson = body
	return s
}

// BodyString defines the body parameters as a string. See documentation.
func (s *TermvectorService) BodyString(body string) *TermvectorService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *TermvectorService) buildURL() (string, url.Values, error) {
	var pathParam = map[string]string{
		"index": s.index,
		"type":  s.typ,
	}
	var path string
	var err error

	// Build URL
	if s.id != "" {
		pathParam["id"] = s.id
		path, err = uritemplates.Expand("/{index}/{type}/{id}/_termvector", pathParam)
	} else {
		path, err = uritemplates.Expand("/{index}/{type}/_termvector", pathParam)
	}

	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.fieldStatistics != nil {
		params.Set("field_statistics", fmt.Sprintf("%v", *s.fieldStatistics))
	}
	if len(s.fields) > 0 {
		params.Set("fields", strings.Join(s.fields, ","))
	}
	if s.offsets != nil {
		params.Set("offsets", fmt.Sprintf("%v", *s.offsets))
	}
	if s.parent != "" {
		params.Set("parent", s.parent)
	}
	if s.payloads != nil {
		params.Set("payloads", fmt.Sprintf("%v", *s.payloads))
	}
	if s.positions != nil {
		params.Set("positions", fmt.Sprintf("%v", *s.positions))
	}
	if s.preference != "" {
		params.Set("preference", s.preference)
	}
	if s.realtime != nil {
		params.Set("realtime", fmt.Sprintf("%v", *s.realtime))
	}
	if s.routing != "" {
		params.Set("routing", s.routing)
	}
	if s.termStatistics != nil {
		params.Set("term_statistics", fmt.Sprintf("%v", *s.termStatistics))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *TermvectorService) Validate() error {
	var invalid []string
	if s.index == "" {
		invalid = append(invalid, "Index")
	}
	if s.typ == "" {
		invalid = append(invalid, "Type")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *TermvectorService) Do() (*TermvectorResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request body
	var body interface{}
	if s.bodyJson != nil {
		body = s.bodyJson
	} else if s.bodyString != "" {
		body = s.bodyString
	} else if s.doc != nil || s.perFieldAnalyzer != nil {
		data := make(map[string]interface{})
		if s.doc != nil {
			data["doc"] = s.doc
		}
		if len(s.perFieldAnalyzer) > 0 {
			data["per_field_analyzer"] = s.perFieldAnalyzer
		}
		body = data
	} else {
		body = ""
	}

	// Get HTTP response
	res, err := s.client.PerformRequest("GET", path, params, body)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(TermvectorResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type TokenInfo struct {
	StartOffset int64  `json:"start_offset"`
	EndOffset   int64  `json:"end_offset"`
	Position    int64  `json:"position"`
	Payload     string `json:"payload"`
}

type TermsInfo struct {
	DocFreq  int64       `json:"doc_freq"`
	TermFreq int64       `json:"term_freq"`
	Ttf      int64       `json:"ttf"`
	Tokens   []TokenInfo `json:"tokens"`
}

type FieldStatistics struct {
	DocCount   int64 `json:"doc_count"`
	SumDocFreq int64 `json:"sum_doc_freq"`
	SumTtf     int64 `json:"sum_ttf"`
}

type TermVectorsFieldInfo struct {
	FieldStatistics FieldStatistics      `json:"field_statistics"`
	Terms           map[string]TermsInfo `json:"terms"`
}

// TermvectorResponse is the response of TermvectorService.Do.
type TermvectorResponse struct {
	Index       string                          `json:"_index"`
	Type        string                          `json:"_type"`
	Id          string                          `json:"_id,omitempty"`
	Version     int                             `json:"_version"`
	Found       bool                            `json:"found"`
	Took        int64                           `json:"took"`
	TermVectors map[string]TermVectorsFieldInfo `json:"term_vectors"`
}
