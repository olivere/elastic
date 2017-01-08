// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
)

// IndicesAnalyzeService analyze an index
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-analyze.html
// for detail

// IndicesAnalyzeService is documented at http://www.elastic.co/guide/en/elasticsearch/reference/master/indices-analyze.html.
type IndicesAnalyzeService struct {
	client      *Client
	pretty      bool
	index       string
	analyzer    string
	attributes  []string
	charFilters []string
	explain     bool
	field       string
	filters     []string
	format      string
	preferLocal bool
	text        []string
	tokenizer   string
	bodyJson    interface{}
	bodyString  string
}

// NewIndicesAnalyzeService creates a new IndicesAnalyzeService.
func NewIndicesAnalyzeService(client *Client) *IndicesAnalyzeService {
	return &IndicesAnalyzeService{
		client:      client,
		attributes:  make([]string, 0),
		charFilters: make([]string, 0),
		filters:     make([]string, 0),
		text:        make([]string, 0),
	}
}

// Index is documented as: The name of the index to scope the operation.
func (s *IndicesAnalyzeService) Index(index string) *IndicesAnalyzeService {
	s.index = index
	return s
}

// Analyzer is documented as: The name of the analyzer to use.
func (s *IndicesAnalyzeService) Analyzer(analyzer string) *IndicesAnalyzeService {
	s.analyzer = analyzer
	return s
}

// Attributes is documented as: A comma-separated list of token attributes to output, this parameter works only with `explain=true`.
func (s *IndicesAnalyzeService) Attributes(attributes []string) *IndicesAnalyzeService {
	s.attributes = attributes
	return s
}

// CharFilters is documented as: A comma-separated list of character filters to use for the analysis.
func (s *IndicesAnalyzeService) CharFilters(charFilters []string) *IndicesAnalyzeService {
	s.charFilters = charFilters
	return s
}

// Explain is documented as: With `true`, outputs more advanced details. (default: false).
func (s *IndicesAnalyzeService) Explain(explain bool) *IndicesAnalyzeService {
	s.explain = explain
	return s
}

// Field is documented as: Use the analyzer configured for this field (instead of passing the analyzer name).
func (s *IndicesAnalyzeService) Field(field string) *IndicesAnalyzeService {
	s.field = field
	return s
}

// Filters is documented as: A comma-separated list of filters to use for the analysis.
func (s *IndicesAnalyzeService) Filters(filters []string) *IndicesAnalyzeService {
	s.filters = filters
	return s
}

// Format of the output.
func (s *IndicesAnalyzeService) Format(format string) *IndicesAnalyzeService {
	s.format = format
	return s
}

// PreferLocal is documented as: With `true`, specify that a local shard should be used if available, with `false`, use a random shard (default: true).
func (s *IndicesAnalyzeService) PreferLocal(preferLocal bool) *IndicesAnalyzeService {
	s.preferLocal = preferLocal
	return s
}

// Text is documented as: The text on which the analysis should be performed (when request body is not used).
func (s *IndicesAnalyzeService) Text(text ...string) *IndicesAnalyzeService {
	s.text = text
	return s
}

// Tokenizer is documented as: The name of the tokenizer to use for the analysis.
func (s *IndicesAnalyzeService) Tokenizer(tokenizer string) *IndicesAnalyzeService {
	s.tokenizer = tokenizer
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *IndicesAnalyzeService) Pretty(pretty bool) *IndicesAnalyzeService {
	s.pretty = pretty
	return s
}

// BodyJson is documented as: The text on which the analysis should be performed.
func (s *IndicesAnalyzeService) BodyJson(body interface{}) *IndicesAnalyzeService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: The text on which the analysis should be performed.
func (s *IndicesAnalyzeService) BodyString(body string) *IndicesAnalyzeService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *IndicesAnalyzeService) buildURL() (string, string, error) {
	// Build URL
	var path string
	if s.index == "" {
		path = "/_analyze"
	} else {
		path = fmt.Sprintf("/%s/_analyze", s.index)
	}

	// Add query string parameters
	params := make(map[string]interface{})
	if s.pretty {
		params["pretty"] = "1"
	}
	if s.analyzer != "" {
		params["analyzer"] = s.analyzer
	}
	if len(s.attributes) > 0 {
		params["attributes"] = s.attributes
	}
	if len(s.charFilters) > 0 {
		params["char_filters"] = s.charFilters
	}
	if s.explain {
		params["explain"] = fmt.Sprintf("%v", s.explain)
	}
	if s.field != "" {
		params["field"] = s.field
	}
	if len(s.filters) > 0 {
		params["filters"] = s.filters
	}
	if s.format != "" {
		params["format"] = s.format
	}
	if s.preferLocal {
		params["prefer_local"] = fmt.Sprintf("%v", s.preferLocal)
	}
	if len(s.text) > 0 {
		params["text"] = s.text
	}
	if s.tokenizer != "" {
		params["tokenizer"] = s.tokenizer
	}

	body, err := json.Marshal(params)
	if err != nil {
		return "", "", err
	}

	return path, string(body), nil
}

func (s *IndicesAnalyzeService) Do() (*IndicesAnalyzeResponse, error) {
	if notValidated := s.Validate(); notValidated != nil {
		return nil, notValidated
	}

	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformRequest("GET", path, nil, params)
	if err != nil {
		return nil, err
	}

	ret := new(IndicesAnalyzeResponse)
	if err = s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *IndicesAnalyzeService) Validate() error {
	var invalid []string
	if len(s.text) == 0 {
		invalid = append(invalid, "Text")
	}

	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

type IndicesAnalyzeResponse struct {
	Tokens []AnalyzeToken `json:"tokens"` //json part for normal message
	Detail AnalyzeDetail  `json:"detail"` //json part for verbose message of explain request
}

type AnalyzeToken struct {
	Token       string `json:"token"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
	Type        string `json:"type"`
	Position    int    `json:"position"`
}

type AnalyzeDetail struct {
	CustomAnalyzer bool          `json:"custom_analyzer"`
	Charfilters    []interface{} `json:"charfilters"`
	Analyzer       struct {
		Name   string `json:"name"`
		Tokens []struct {
			Token          string `json:"token"`
			StartOffset    int    `json:"start_offset"`
			EndOffset      int    `json:"end_offset"`
			Type           string `json:"type"`
			Position       int    `json:"position"`
			Bytes          string `json:"bytes"`
			PositionLength int    `json:"positionLength"`
		} `json:"tokens"`
	} `json:"analyzer"`
	Tokenizer struct {
		Name   string `json:"name"`
		Tokens []struct {
			Token       string `json:"token"`
			StartOffset int    `json:"start_offset"`
			EndOffset   int    `json:"end_offset"`
			Type        string `json:"type"`
			Position    int    `json:"position"`
		} `json:"tokens"`
	} `json:"tokenizer"`
	Tokenfilters []struct {
		Name   string `json:"name"`
		Tokens []struct {
			Token       string `json:"token"`
			StartOffset int    `json:"start_offset"`
			EndOffset   int    `json:"end_offset"`
			Type        string `json:"type"`
			Position    int    `json:"position"`
			Keyword     bool   `json:"keyword"`
		} `json:"tokens"`
	} `json:"tokenfilters"`
}
