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

// SecurityPutConfigService update a config by its name.
// See https://opensearch.org/docs/latest/security/access-control/api/#update-configuration
type SecurityPutConfigService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	body interface{}
}

// NewSecurityPutConfigService creates a new SecurityPutConfigService.
func NewSecurityPutConfigService(client *Client) *SecurityPutConfigService {
	return &SecurityPutConfigService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityPutConfigService) Pretty(pretty bool) *SecurityPutConfigService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityPutConfigService) Human(human bool) *SecurityPutConfigService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityPutConfigService) ErrorTrace(errorTrace bool) *SecurityPutConfigService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityPutConfigService) FilterPath(filterPath ...string) *SecurityPutConfigService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityPutConfigService) Header(name string, value string) *SecurityPutConfigService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityPutConfigService) Headers(headers http.Header) *SecurityPutConfigService {
	s.headers = headers
	return s
}

// Body specifies the config. Use a string or a type that will get serialized as JSON.
func (s *SecurityPutConfigService) Body(body interface{}) *SecurityPutConfigService {
	s.body = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityPutConfigService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/securityconfig/config", nil)
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
func (s *SecurityPutConfigService) Validate() error {
	var invalid []string
	if s.body == nil {
		invalid = append(invalid, "Body")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SecurityPutConfigService) Do(ctx context.Context) (*SecurityPutConfigResponse, error) {
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
	ret := new(SecurityPutConfigResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityPutConfigResponse is the response of SecurityPutConfigService.Do.
type SecurityPutConfigResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
