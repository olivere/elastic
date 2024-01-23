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

// SecurityGetUserService retrieves a role by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#get-user
type SecurityGetUserService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityGetUserService creates a new SecurityGetUserService.
func NewSecurityGetUserService(client *Client) *SecurityGetUserService {
	return &SecurityGetUserService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetUserService) Pretty(pretty bool) *SecurityGetUserService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetUserService) Human(human bool) *SecurityGetUserService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetUserService) ErrorTrace(errorTrace bool) *SecurityGetUserService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetUserService) FilterPath(filterPath ...string) *SecurityGetUserService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetUserService) Header(name string, value string) *SecurityGetUserService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetUserService) Headers(headers http.Header) *SecurityGetUserService {
	s.headers = headers
	return s
}

// Name is name of the role to retrieve.
func (s *SecurityGetUserService) Name(name string) *SecurityGetUserService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetUserService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/internalusers/{name}", map[string]string{
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
func (s *SecurityGetUserService) Validate() error {
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
func (s *SecurityGetUserService) Do(ctx context.Context) (*SecurityGetUserResponse, error) {
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
	ret := new(SecurityGetUserResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetUserResponse is the response of SecurityGetUserService.Do.
type SecurityGetUserResponse map[string]SecurityUser

// SecurityUser is the user object.
// source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/InternalUserV7.java
type SecurityUser struct {
	SecurityUserBase `json:",inline"`
	Reserved         *bool `json:"reserved,omitempty"`
	Hidden           *bool `json:"hidden,omitempty"`
	Static           *bool `json:"static,omitempty"`
}

type SecurityUserBase struct {
	Hash          *string           `json:"hash,omitempty"`
	BackendRoles  []string          `json:"backend_roles,omitempty"`
	SecurityRoles []string          `json:"opendistro_security_roles,omitempty"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	Description   *string           `json:"description,omitempty"`
	Enabled       *bool             `json:"enabled,omitempty"`
	Service       *bool             `json:"service,omitempty"`
}
