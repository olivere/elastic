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

// SecurityGetActionGroupService retrieves a action group by its name.
// https://opensearch.org/docs/latest/security/access-control/api/#get-action-groups
type SecurityGetActionGroupService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSecurityGetActionGroupService creates a new SecurityGetActionGroupService.
func NewSecurityGetActionGroupService(client *Client) *SecurityGetActionGroupService {
	return &SecurityGetActionGroupService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetActionGroupService) Pretty(pretty bool) *SecurityGetActionGroupService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetActionGroupService) Human(human bool) *SecurityGetActionGroupService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetActionGroupService) ErrorTrace(errorTrace bool) *SecurityGetActionGroupService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetActionGroupService) FilterPath(filterPath ...string) *SecurityGetActionGroupService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetActionGroupService) Header(name string, value string) *SecurityGetActionGroupService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetActionGroupService) Headers(headers http.Header) *SecurityGetActionGroupService {
	s.headers = headers
	return s
}

// Name is name of the action group to retrieve.
func (s *SecurityGetActionGroupService) Name(name string) *SecurityGetActionGroupService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetActionGroupService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/actiongroups/{name}", map[string]string{
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
func (s *SecurityGetActionGroupService) Validate() error {
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
func (s *SecurityGetActionGroupService) Do(ctx context.Context) (*SecurityGetActionGroupResponse, error) {
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
	ret := new(SecurityGetActionGroupResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetActionGroupResponse is the response of SecurityGetActionGroupService.Do.
type SecurityGetActionGroupResponse map[string]SecurityActionGroup

// SecurityActionGroup is the action group object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/ActionGroupsV7.java
type SecurityActionGroup struct {
	Reserved       bool     `json:"reserved,omitempty"`
	Hidden         bool     `json:"hidden,omitempty"`
	Static         bool     `json:"static,omitempty"`
	Description    string   `json:"description,omitempty"`
	Type           string   `json:"type,omitempty"`
	AllowedActions []string `json:"allowed_actions"`
}
