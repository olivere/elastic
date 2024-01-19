package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/disaster37/opensearch/v2/uritemplates"
)

// SecurityGetDistinguishedNameService retrieves a distinguished name by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#get-distinguished-names
type SecurityGetDistinguishedNameService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityGetDistinguishedNameService creates a new SecurityGetDistinguishedNameService.
func NewSecurityGetDistinguishedNameService(client *Client) *SecurityGetDistinguishedNameService {
	return &SecurityGetDistinguishedNameService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetDistinguishedNameService) Pretty(pretty bool) *SecurityGetDistinguishedNameService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetDistinguishedNameService) Human(human bool) *SecurityGetDistinguishedNameService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetDistinguishedNameService) ErrorTrace(errorTrace bool) *SecurityGetDistinguishedNameService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetDistinguishedNameService) FilterPath(filterPath ...string) *SecurityGetDistinguishedNameService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetDistinguishedNameService) Header(name string, value string) *SecurityGetDistinguishedNameService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetDistinguishedNameService) Headers(headers http.Header) *SecurityGetDistinguishedNameService {
	s.headers = headers
	return s
}

// Name is name of the distinguished name to retrieve.
func (s *SecurityGetDistinguishedNameService) Name(name string) *SecurityGetDistinguishedNameService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetDistinguishedNameService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/nodesdn/{name}", map[string]string{
		"name": s.name,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *SecurityGetDistinguishedNameService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SecurityGetDistinguishedNameService) Do(ctx context.Context) (*SecurityGetDistinguishedNameResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method:  "GET",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SecurityGetDistinguishedNameResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetDistinguishedNameResponse is the response of SecurityGetDistinguishedNameService.Do.
type SecurityGetDistinguishedNameResponse map[string]SecurityDistinguishedName

// SecurityDistinguishedName is the dn object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/NodesDn.java
type SecurityDistinguishedName struct {
	NodesDN []string                    `json:"nodes_dn"`
}
