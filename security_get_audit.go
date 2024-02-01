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

// SecurityGetAuditService retrieves a audit.
// See https://opensearch.org/docs/latest/security/access-control/api/#audit-logs
type SecurityGetAuditService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers
}

// NewSecurityGetAuditService creates a new SecurityGetAuditService.
func NewSecurityGetAuditService(client *Client) *SecurityGetAuditService {
	return &SecurityGetAuditService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetAuditService) Pretty(pretty bool) *SecurityGetAuditService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetAuditService) Human(human bool) *SecurityGetAuditService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetAuditService) ErrorTrace(errorTrace bool) *SecurityGetAuditService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetAuditService) FilterPath(filterPath ...string) *SecurityGetAuditService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetAuditService) Header(name string, value string) *SecurityGetAuditService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetAuditService) Headers(headers http.Header) *SecurityGetAuditService {
	s.headers = headers
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetAuditService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/audit", nil)
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
func (s *SecurityGetAuditService) Validate() error {
	var invalid []string
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SecurityGetAuditService) Do(ctx context.Context) (*SecurityGetAuditResponse, error) {
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
	ret := new(SecurityGetAuditResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetAuditResponse is the response of SecurityGetAuditService.Do.
type SecurityGetAuditResponse struct {
	Config SecurityAudit `json:"config"`
}

// SecurityAudit is the audit object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/AuditV7.java
type SecurityAudit struct {
	Enabled    *bool                   `json:"enabled,omitempty"`
	Compliance SecurityAuditCompliance `json:"compliance"`
	Audit      SecurityAuditSpec       `json:"audit"`
}

// SecurityAuditSpec is the audit spec
type SecurityAuditSpec struct {
	IgnoreUsers                 []string `json:"ignore_users,omitempty"`
	IgnoreRequests              []string `json:"ignore_requests,omitempty"`
	DisabledRestCategories      []string `json:"disabled_rest_categories,omitempty"`
	DisabledTransportCategories []string `json:"disabled_transport_categories,omitempty"`
	LogRequestBody              *bool    `json:"log_request_body,omitempty"`
	ResolveIndices              *bool    `json:"resolve_indices,omitempty"`
	ResolveBulkRequests         *bool    `json:"resolve_bulk_requests,omitempty"`
	ExcludeSensitiveHeaders     *bool    `json:"exclude_sensitive_headers,omitempty"`
	EnableTransport             *bool    `json:"enable_transport,omitempty"`
	EnableRest                  *bool    `json:"enable_rest,omitempty"`
}

// SecurityAuditCompliance is the compliance spec
type SecurityAuditCompliance struct {
	Enabled             *bool               `json:"enabled,omitempty"`
	WriteLogDiffs       *bool               `json:"write_log_diffs,omitempty"`
	ReadWatchedFields   map[string][]string `json:"read_watched_fields,omitempty"`
	ReadIgnoreUsers     []string            `json:"read_ignore_users,omitempty"`
	WriteWatchedIndices []string            `json:"write_watched_indices,omitempty"`
	WriteIgnoreUsers    []string            `json:"write_ignore_users,omitempty"`
	ReadMetadataOnly    *bool               `json:"read_metadata_only,omitempty"`
	WriteMetadataOnly   *bool               `json:"write_metadata_only,omitempty"`
	ExternalConfig      *bool               `json:"external_config,omitempty"`
	InternalConfig      *bool               `json:"internal_config,omitempty"`
}
