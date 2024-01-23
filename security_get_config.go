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

// SecurityGetConfigService retrieves a config.
// See https://opensearch.org/docs/latest/security/access-control/api/#get-configuration
type SecurityGetConfigService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers
}

// NewSecurityGetConfigService creates a new SecurityGetConfigService.
func NewSecurityGetConfigService(client *Client) *SecurityGetConfigService {
	return &SecurityGetConfigService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SecurityGetConfigService) Pretty(pretty bool) *SecurityGetConfigService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SecurityGetConfigService) Human(human bool) *SecurityGetConfigService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SecurityGetConfigService) ErrorTrace(errorTrace bool) *SecurityGetConfigService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SecurityGetConfigService) FilterPath(filterPath ...string) *SecurityGetConfigService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SecurityGetConfigService) Header(name string, value string) *SecurityGetConfigService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SecurityGetConfigService) Headers(headers http.Header) *SecurityGetConfigService {
	s.headers = headers
	return s
}

// buildURL builds the URL for the operation.
func (s *SecurityGetConfigService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_security/api/securityconfig", nil)
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
func (s *SecurityGetConfigService) Validate() error {
	var invalid []string
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SecurityGetConfigService) Do(ctx context.Context) (*SecurityGetConfigResponse, error) {
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
	ret := new(SecurityGetConfigResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SecurityGetConfigResponse is the response of SecurityGetConfigService.Do.
type SecurityGetConfigResponse struct {
	Config SecurityConfig `json:"config"`
}

// SecurityConfig is the config object.
// Source code: https://github.com/opensearch-project/security/blob/main/src/main/java/org/opensearch/security/securityconf/impl/v7/ConfigV7.java
type SecurityConfig struct {
	Dynamic SecurityConfigDynamic `json:"dynamic"`
}

type SecurityConfigDynamic struct {
	FilteredAliasMode            *string                                       `json:"filtered_alias_mode,omitempty"`
	DisableRestAuth              *bool                                         `json:"disable_rest_auth,omitempty"`
	DisableIntertransportAuth    *bool                                         `json:"disable_intertransport_auth,omitempty"`
	RespectRequestIndicesOptions *bool                                         `json:"respect_request_indices_options,omitempty"`
	License                      *string                                       `json:"license,omitempty"`
	Kibana                       *SecurityConfigKibana                         `json:"kibana,omitempty"`
	Http                         *SecurityConfigHttp                           `json:"http,omitempty"`
	Authc                        map[string]SecurityConfigAuthc                `json:"authc,omitempty"`
	Authz                        map[string]SecurityConfigAuthz                `json:"authz,omitempty"`
	AuthFailureListeners         map[string]SecurityConfigAuthFailureListeners `json:"auth_failure_listeners,omitempty"`
	DoNotFailOnForbidden         *bool                                         `json:"do_not_fail_on_forbidden,omitempty"`
	MultiRolespanEnabled         *bool                                         `json:"multi_rolespan_enabled,omitempty"`
	HostsResolverMode            *string                                       `json:"hosts_resolver_mode,omitempty"`
	TransportUserrnameAttribute  *string                                       `json:"transport_userrname_attribute,omitempty"`
	DoNotFailOnForbiddenEmpty    *bool                                         `json:"do_not_fail_on_forbidden_empty,omitempty"`
	OnBehalfOfSettings           *SecurityConfigOnBehalfOfSettings             `json:"on_behalf_of,omitempty"`
}

type SecurityConfigKibana struct {
	MultitenancyEnabled  *bool   `json:"multitenancy_enabled,omitempty"`
	PrivateTenantEnabled *bool   `json:"private_tenant_enabled,omitempty"`
	DefaultTenant        *string `json:"default_tenant,omitempty"`
	ServerUsername       *string `json:"server_username,omitempty"`
	OpendistroRole       *string `json:"opendistro_role,omitempty"`
	Index                *string `json:"index,omitempty"`
}

type SecurityConfigHttp struct {
	AnonymousAuthEnabled *bool              `json:"anonymous_auth_enabled,omitempty"`
	Xff                  *SecurityConfigXff `json:"xff,omitempty"`
}

type SecurityConfigAuthc struct {
	HttpEnabled            *bool                                `json:"http_enabled,omitempty"`
	TransportEnabled       *bool                                `json:"transport_enabled,omitempty"`
	Order                  *int64                               `json:"order,omitempty"`
	HttpAuthenticator      *SecurityConfigHttpAuthenticator     `json:"http_authenticator,omitempty"`
	AuthenticationBackendd *SecurityConfigAuthenticationBackend `json:"authentication_backend,omitempty"`
	Description            *string                              `json:"description,omitempty"`
}

type SecurityConfigAuthz struct {
	HttpEnabled          *bool                               `json:"http_enabled,omitempty"`
	TransportEnabled     *bool                               `json:"transport_enabled,omitempty"`
	AuthorizationBackend *SecurityConfigAuthorizationBackend `json:"authorization_backend,omitempty"`
	Description          *string                             `json:"description,omitempty"`
}

type SecurityConfigAuthFailureListeners struct {
	Type                  *string `json:"type,omitempty"`
	AuthenticationBackend *string `json:"authentication_backend,omitempty"`
	AllowedTries          *int64  `json:"allowed_tries,omitempty"`
	TimeWindowSeconds     *int64  `json:"time_window_seconds,omitempty"`
	BlockExpirySeconds    *int64  `json:"block_expiry_seconds,omitempty"`
	MaxBlockedClients     *int64  `json:"max_blocked_clients,omitempty"`
	MaxTrackedClients     *int64  `json:"max_tracked_clients,omitempty"`
}

type SecurityConfigOnBehalfOfSettings struct {
	OboEnabled    *bool   `json:"oboEnabled,omitempty"`
	SigningKey    *string `json:"signingKey,omitempty"`
	EncryptionKey *string `json:"encryptionKey,omitempty"`
}

type SecurityConfigXff struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	InternalProxies *string `json:"internalProxies,omitempty"`
	RemoteIpHeader  *string `json:"remoteIpHeader,omitempty"`
}

type SecurityConfigHttpAuthenticator struct {
	Challenge *bool          `json:"challenge,omitempty"`
	Type      *string        `json:"type,omitempty"`
	Config    map[string]any `json:"config,omitempty"`
}

type SecurityConfigAuthenticationBackend struct {
	Type   *string        `json:"type,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

type SecurityConfigAuthorizationBackend struct {
	Type   *string        `json:"type,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}
