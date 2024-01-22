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

// SecurityGetRoleService retrieves a role by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#get-role
type SecurityGetRoleService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityGetRoleService creates a new SecurityGetRoleService.
func NewSecurityGetRoleService(client *Client) *SecurityGetRoleService {
	return &SecurityGetRoleService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetRoleService) Pretty(pretty bool) *SecurityGetRoleService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetRoleService) Human(human bool) *SecurityGetRoleService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetRoleService) ErrorTrace(errorTrace bool) *SecurityGetRoleService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetRoleService) FilterPath(filterPath ...string) *SecurityGetRoleService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetRoleService) Header(name string, value string) *SecurityGetRoleService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetRoleService) Headers(headers http.Header) *SecurityGetRoleService {
	s.headers = headers
	return s
}

// Name is name of the role to retrieve.
func (s *SecurityGetRoleService) Name(name string) *SecurityGetRoleService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetRoleService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/roles/{name}", map[string]string{
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
func (s *SecurityGetRoleService) Validate() error {
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
func (s *SecurityGetRoleService) Do(ctx context.Context) (*SecurityGetRoleResponse, error) {
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
	ret := new(SecurityGetRoleResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetRoleResponse is the response of SecurityGetRoleService.Do.
type SecurityGetRoleResponse map[string]SecurityRole

// SecurityRole is the role object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/RoleV7.java
type SecurityRole struct {
	SecurityPutRole SecurityPutRole `json:",inline"`
	Reserved        *bool           `json:"reserved,omitempty"`
	Hidden          *bool           `json:"hidden,omitempty"`
	Static          *bool           `json:"static,omitempty"`
}

// SecurityTenantPermissions is the tenant permission object
type SecurityTenantPermissions struct {
	TenantPatterns []string `json:"tenant_patterns"`
	AllowedAction  []string `json:"allowed_actions"`
}

// SecurityIndexPermissions is the index permission object
type SecurityIndexPermissions struct {
	IndexPatterns         []string `json:"index_patterns"`
	MaskedFields          []string `json:"masked_fields"`
	AllowedActions        []string `json:"allowed_actions"`
	DocumentLevelSecurity *string  `json:"dls,omitempty"`
	FieldLevelSecurity    []string `json:"fls,omitempty"`
}
