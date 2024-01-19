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

// SecurityGetRoleMappingService retrieves a role by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#get-role-mapping
type SecurityGetRoleMappingService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityGetRoleMappingService creates a new SecurityGetRoleMappingService.
func NewSecurityGetRoleMappingService(client *Client) *SecurityGetRoleMappingService {
	return &SecurityGetRoleMappingService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetRoleMappingService) Pretty(pretty bool) *SecurityGetRoleMappingService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetRoleMappingService) Human(human bool) *SecurityGetRoleMappingService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetRoleMappingService) ErrorTrace(errorTrace bool) *SecurityGetRoleMappingService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetRoleMappingService) FilterPath(filterPath ...string) *SecurityGetRoleMappingService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetRoleMappingService) Header(name string, value string) *SecurityGetRoleMappingService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetRoleMappingService) Headers(headers http.Header) *SecurityGetRoleMappingService {
	s.headers = headers
	return s
}

// Name is name of the role to retrieve.
func (s *SecurityGetRoleMappingService) Name(name string) *SecurityGetRoleMappingService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetRoleMappingService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/rolesmapping/{name}", map[string]string{
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
func (s *SecurityGetRoleMappingService) Validate() error {
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
func (s *SecurityGetRoleMappingService) Do(ctx context.Context) (*SecurityGetRoleMappingResponse, error) {
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
	ret := new(SecurityGetRoleMappingResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetRoleMappingResponse is the response of SecurityGetRoleMappingService.Do.
type SecurityGetRoleMappingResponse map[string]SecurityRoleMapping

// SecurityRoleMapping is the role mapping object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/RoleMappingsV7.java
type SecurityRoleMapping struct {
	BackendRoles    []string `json:"backend_roles"`
	AndBackendRoles []string `json:"and_backend_roles"`
	Hosts           []string `json:"hosts"`
	Users           []string `json:"users"`
	Reserved        bool     `json:"reserved,omitempty"`
	Hidden          bool     `json:"hidden,omitempty"`
}
