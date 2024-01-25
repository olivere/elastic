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

// SmGetPolicyService get a ISM policy by its name.
// See https://opensearch.org/docs/latest/tuning-your-cluster/availability-and-recovery/snapshots/sm-api/#get-policies
type SmGetPolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewSmGetPolicyService creates a new SmGetPolicyService.
func NewSmGetPolicyService(client *Client) *SmGetPolicyService {
	return &SmGetPolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SmGetPolicyService) Pretty(pretty bool) *SmGetPolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SmGetPolicyService) Human(human bool) *SmGetPolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SmGetPolicyService) ErrorTrace(errorTrace bool) *SmGetPolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SmGetPolicyService) FilterPath(filterPath ...string) *SmGetPolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SmGetPolicyService) Header(name string, value string) *SmGetPolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SmGetPolicyService) Headers(headers http.Header) *SmGetPolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to get.
func (s *SmGetPolicyService) Name(name string) *SmGetPolicyService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *SmGetPolicyService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_sm/policies/{name}", map[string]string{
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
func (s *SmGetPolicyService) Validate() error {
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
func (s *SmGetPolicyService) Do(ctx context.Context) (*SmGetPolicyResponse, error) {
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
	ret := new(SmGetPolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SmGetPolicyResponse is the get policy response object
// https://opensearch.org/docs/latest/im-plugin/ism/api/#get-policy
type SmGetPolicyResponse struct {
	Id             string   `json:"_id"`
	Version        int64    `json:"_version"`
	SequenceNumber int64    `json:"_seq_no"`
	PrimaryTerm    int64    `json:"_primary_term"`
	Policy         SmPolicy `json:"sm_policy"`
}

// SmPolicy is the snapshot policy object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/snapshotmanagement/model/SMPolicy.kt
type SmPolicy struct {
	SmPutPolicy    `json:",inline"`
	Name           *string        `json:"policy_id,omitempty"`
	SchemaVersion  *int64         `json:"schema_version,omitempty"`
	LastUpdateTime *int64         `json:"last_updated_time,omitempty"`
	EnabledTime    *int64         `json:"enabled_time,omitempty"`
	Schedule       map[string]any `json:"schedule,omitempty"`
}
