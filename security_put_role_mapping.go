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

// SecurityPutRoleMappingService retrieves a role by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#create-role-mapping
type SecurityPutRoleMappingService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
	body interface{}
}

// NewSecurityPutRoleMappingService creates a new SecurityPutRoleMappingService.
func NewSecurityPutRoleMappingService(client *Client) *SecurityPutRoleMappingService {
	return &SecurityPutRoleMappingService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityPutRoleMappingService) Pretty(pretty bool) *SecurityPutRoleMappingService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityPutRoleMappingService) Human(human bool) *SecurityPutRoleMappingService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityPutRoleMappingService) ErrorTrace(errorTrace bool) *SecurityPutRoleMappingService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityPutRoleMappingService) FilterPath(filterPath ...string) *SecurityPutRoleMappingService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityPutRoleMappingService) Header(name string, value string) *SecurityPutRoleMappingService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityPutRoleMappingService) Headers(headers http.Header) *SecurityPutRoleMappingService {
	s.headers = headers
	return s
}

// Name is name of the role to create.
func (s *SecurityPutRoleMappingService) Name(name string) *SecurityPutRoleMappingService {
	s.name = name
	return s
}

// Body specifies the role. Use a string or a type that will get serialized as JSON.
func (s *SecurityPutRoleMappingService) Body(body interface{}) *SecurityPutRoleMappingService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityPutRoleMappingService) buildURL() (string, url.Values, error) {
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
func (s *SecurityPutRoleMappingService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
	}
	if s.body == nil {
		invalid = append(invalid, "Body")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SecurityPutRoleMappingService) Do(ctx context.Context) (*SecurityPutRoleMappingResponse, error) {
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
		Method:  "PUT",
		Path:    path,
		Params:  params,
		Body:    s.body,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SecurityPutRoleMappingResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityPutRoleMappingResponse is the response of SecurityPutRoleMappingService.Do.
type SecurityPutRoleMappingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
