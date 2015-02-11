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
type CreateMappingsService struct {
	client     *Client
	indices    []string
	types      []string
	body       map[string]TypeMappings
	bodyString string
	debug      bool
	pretty     bool
}

// NewMappingsService creates a new mappings service
func NewCreateMappingsService(client *Client) *CreateMappingsService {
	return &CreateMappingsService{
		client: client,
	}
}

// Index set the desired index
func (s *CreateMappingsService) Index(index string) *CreateMappingsService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, index)
	return s
}

// Indices sets one or more desired indices
func (s *CreateMappingsService) Indices(indices ...string) *CreateMappingsService {
	if s.indices == nil {
		s.indices = make([]string, 0)
	}
	s.indices = append(s.indices, indices...)
	return s
}

// Type set the desired type
func (s *CreateMappingsService) Type(typ string) *CreateMappingsService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, typ)
	return s
}

// Types set one or more desired types
func (s *CreateMappingsService) Types(types ...string) *CreateMappingsService {
	if s.types == nil {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, types...)
	return s
}

// Pretty enable/disable pretty print of JSON request/response debug
func (s *CreateMappingsService) Pretty(pretty bool) *CreateMappingsService {
	s.pretty = pretty
	return s
}

// Debug enable/disable debug
func (s *CreateMappingsService) Debug(debug bool) *CreateMappingsService {
	s.debug = debug
	return s
}

func (s *CreateMappingsService) Body(mappings map[string]TypeMappings) *CreateMappingsService {
	s.body = mappings
	return s
}

func (s *CreateMappingsService) BodyString(body string) *CreateMappingsService {
	s.bodyString = body
	return s
}

// Do make the request
func (s *CreateMappingsService) Do() (map[string]bool, error) {
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

	// mappings
	urls += "/_mappings"

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

	// Parameters
	params := make(url.Values)

	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("PUT", urls)

	if err != nil {
		return nil, err
	}

	if s.body != nil {
		req.SetBodyJson(&s.body)
	} else if s.bodyString != "" {
		req.SetBodyString(s.bodyString)
	} else {
		return nil, errors.New("Create mappings requires Body or BodyJson.")
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

	ret := map[string]bool{}

	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	if ret != nil {
		return ret, nil
	}

	return nil, errors.New("Failed to put mappings")
}
