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

// IsmDeletePolicyService delete a ISM policy by its name.
// See https://opensearch.org/docs/latest/im-plugin/ism/api/#delete-policy
type IsmDeletePolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewIsmDeletePolicyService creates a new IsmDeletePolicyService.
func NewIsmDeletePolicyService(client *Client) *IsmDeletePolicyService {
	return &IsmDeletePolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *IsmDeletePolicyService) Pretty(pretty bool) *IsmDeletePolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *IsmDeletePolicyService) Human(human bool) *IsmDeletePolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *IsmDeletePolicyService) ErrorTrace(errorTrace bool) *IsmDeletePolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *IsmDeletePolicyService) FilterPath(filterPath ...string) *IsmDeletePolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *IsmDeletePolicyService) Header(name string, value string) *IsmDeletePolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *IsmDeletePolicyService) Headers(headers http.Header) *IsmDeletePolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to delete.
func (s *IsmDeletePolicyService) Name(name string) *IsmDeletePolicyService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *IsmDeletePolicyService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_ism/policies/{name}", map[string]string{
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
func (s *IsmDeletePolicyService) Validate() error {
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
func (s *IsmDeletePolicyService) Do(ctx context.Context) (*IsmDeletePolicyResponse, error) {
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
		Method:  "DELETE",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(IsmDeletePolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type IsmDeletePolicyResponse struct {
	Index          *string        `json:"_index,omitempty"`
	ID             *string        `json:"_id,omitempty"`
	Version        *int64         `json:"_version,omitempty"`
	Result         *string        `json:"result,omitempty"`
	ForcedRefresh  *bool          `json:"forced_refresh,omitempty"`
	Shards         map[string]any `json:"_shards,omitempty"`
	SequenceNumber *int64         `json:"_seq_no,omitempty"`
	PrimaryTerm    *int64         `json:"_primary_term,omitempty"`
}
