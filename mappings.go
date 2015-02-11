package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// MappingsService is a convenient service for create and
// retrieve elasticsearch mappings.
type MappingsService struct {
	client  *Client
	indices []string
	types   []string
	debug   bool
	pretty  bool
}

// FieldProperty stores field properties
type FieldProperty struct {
	Type                       string            `json:"type"`
	Format                     string            `json:"format,omitempty"`
	FieldData                  map[string]string `json:"fielddata,omitempty"`
	LatLon                     bool              `json:"lat_lon,omitempty"`
	GeoHash                    bool              `json:"geohash,omitempty"`
	GeoHashPrefix              bool              `json:"geohash_prefix,omitempty"`
	GeoHashPrecision           int               `json:"geohash_precision,omitempty"`
	Tree                       string            `json:"tree,omitempty"`
	TreeLevels                 int               `json:"tree_levels,omitempty"`
	Analyzer                   string            `json:"analyzer,omitempty"`
	IndexAnalyzer              string            `json:"index_analyzer,omitempty"`
	SearchAnalyzer             string            `json:"search_analyzer,omitempty"`
	Payloads                   bool              `json:"payloads,omitempty"`
	PreserveSeparators         bool              `json:"preserve_separators,omitempty"`
	PreservePositionIncrements bool              `json:"preserve_position_increments,omitempty"`
	MaxInputLength             int               `json:"max_input_length,omitempty"`
}

// TypeMappings stores the type's mappings
type TypeMappings struct {
	Properties map[string]FieldProperty `json:"properties"`
}

// Mappings stores the /_mappings result section
type Mappings struct {
	Mappings map[string]TypeMappings `json:"mappings,omitempty"`
}

// NewMappingsService creates a new mappings service
func NewMappingsService(client *Client) *MappingsService {
	return &MappingsService{
		client: client,
	}
}

// Index set the desired index
func (s *MappingsService) Index(index string) *MappingsService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, index)
	return s
}

// Indices sets one or more desired indices
func (s *MappingsService) Indices(indices ...string) *MappingsService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, indices...)
	return s
}

// Type set the desired type
func (s *MappingsService) Type(typ string) *MappingsService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, typ)
	return s
}

// Types set one or more desired types
func (s *MappingsService) Types(types ...string) *MappingsService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, types...)
	return s
}

// Pretty enable/disable pretty print of JSON request/response debug
func (s *MappingsService) Pretty(pretty bool) *MappingsService {
	s.pretty = pretty
	return s
}

// Debug enable/disable debug
func (s *MappingsService) Debug(debug bool) *MappingsService {
	s.debug = debug
	return s
}

// Do make the request
func (s *MappingsService) Do() (map[string]Mappings, error) {
	var err error

	// Build url
	urls := "/"

	// Indices part
	indexPart := make([]string, 0)
	for _, index := range s.indices {
		index, err = uritemplates.Expand("{index}", map[string]string{
			"index": index,
		})
		if err != nil {
			return nil, err
		}
		indexPart = append(indexPart, index)
	}
	if len(indexPart) > 0 {
		urls += strings.Join(indexPart, ",")
	}

	// Types part
	typesPart := make([]string, 0)
	for _, typ := range s.types {
		typ, err = uritemplates.Expand("{type}", map[string]string{
			"type": typ,
		})
		if err != nil {
			return nil, err
		}
		typesPart = append(typesPart, typ)
	}
	if len(typesPart) > 0 {
		urls += "/" + strings.Join(typesPart, ",")
	}

	// Search
	urls += "/_mappings"

	// Parameters
	params := make(url.Values)
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("GET", urls)

	if err != nil {
		return nil, err
	}

	if s.debug {
		s.client.dumpRequest((*http.Request)(req))
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

	if s.debug {
		s.client.dumpResponse(res)
	}

	ret := map[string]Mappings{}

	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	if ret != nil {
		return ret, nil
	}

	return nil, errors.New("Failed to get mappings")
}
