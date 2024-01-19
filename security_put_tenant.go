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

// SecurityPutTenantService update a tenant by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#create-tenant
type SecurityPutTenantService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
	body interface{}
}

// NewSecurityPutTenantService creates a new SecurityPutTenantService.
func NewSecurityPutTenantService(client *Client) *SecurityPutTenantService {
	return &SecurityPutTenantService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityPutTenantService) Pretty(pretty bool) *SecurityPutTenantService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityPutTenantService) Human(human bool) *SecurityPutTenantService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityPutTenantService) ErrorTrace(errorTrace bool) *SecurityPutTenantService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityPutTenantService) FilterPath(filterPath ...string) *SecurityPutTenantService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityPutTenantService) Header(name string, value string) *SecurityPutTenantService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityPutTenantService) Headers(headers http.Header) *SecurityPutTenantService {
	s.headers = headers
	return s
}

// Name is name of the tenant to create.
func (s *SecurityPutTenantService) Name(name string) *SecurityPutTenantService {
	s.name = name
	return s
}

// Body specifies the tenant. Use a string or a type that will get serialized as JSON.
func (s *SecurityPutTenantService) Body(body interface{}) *SecurityPutTenantService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityPutTenantService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/tenants/{name}", map[string]string{
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
func (s *SecurityPutTenantService) Validate() error {
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
func (s *SecurityPutTenantService) Do(ctx context.Context) (*SecurityPutTenantResponse, error) {
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
	ret := new(SecurityPutTenantResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityPutTenantResponse is the response of SecurityPutTenantService.Do.
type SecurityPutTenantResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SecurityPutTenant struct {
	Description string `json:"description"`
}
