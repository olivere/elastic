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

// IsmGetPolicyService get a ISM policy by its name.
// See https://opensearch.org/docs/latest/im-plugin/ism/api/#get-policy
type IsmGetPolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewIsmGetPolicyService creates a new IsmGetPolicyService.
func NewIsmGetPolicyService(client *Client) *IsmGetPolicyService {
	return &IsmGetPolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *IsmGetPolicyService) Pretty(pretty bool) *IsmGetPolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *IsmGetPolicyService) Human(human bool) *IsmGetPolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *IsmGetPolicyService) ErrorTrace(errorTrace bool) *IsmGetPolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *IsmGetPolicyService) FilterPath(filterPath ...string) *IsmGetPolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *IsmGetPolicyService) Header(name string, value string) *IsmGetPolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *IsmGetPolicyService) Headers(headers http.Header) *IsmGetPolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to get.
func (s *IsmGetPolicyService) Name(name string) *IsmGetPolicyService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *IsmGetPolicyService) buildURL() (string, url.Values, error) {
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
func (s *IsmGetPolicyService) Validate() error {
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
func (s *IsmGetPolicyService) Do(ctx context.Context) (*IsmGetPolicyResponse, error) {
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
	ret := new(IsmGetPolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// IsmGetPolicyResponse is the get index state management response object
// https://opensearch.org/docs/latest/im-plugin/ism/api/#get-policy
type IsmGetPolicyResponse struct {
	Id             string       `json:"_id"`
	Version        int64        `json:"_version"`
	SequenceNumber int64        `json:"_seq_no"`
	PrimaryTerm    int64        `json:"_primary_term"`
	Policy         IsmGetPolicy `json:"policy"`
}

// IsmPolicyBase is the base ISM policy
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/Policy.kt
type IsmPolicyBase struct {
	ID                *string               `json:"policy_id,omitempty"`
	Description       *string               `json:"description,omitempty"`
	ErrorNotification *IsmErrorNotification `json:"error_notification,omitempty"`
	DefaultState      *string               `json:"default_state,omitempty"`
	States            []IsmPolicyState      `json:"states,omitempty"`
	IsmTemplate       []IsmPolicyTemplate   `json:"ism_template,omitempty"`
}

// IsmGetPolicy is the ISM policy
type IsmGetPolicy struct {
	IsmPolicyBase   `json:",inline"`
	SchemaVersion   *int64 `json:"schema_version,omitempty"`
	LastUpdatedTime *int64 `json:"last_updated_time,omitempty"`
}

// IsmErrorNotification is the error notification object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/ErrorNotification.kt
type IsmErrorNotification struct {
	Destination     *IsmErrorNotificationDestination     `json:"destination,omitempty"`
	Channel         *IsmErrorNotificationChannel         `json:"channel,omitempty"`
	MessageTemplate *IsmErrorNotificationMessageTemplate `json:"message_template,omitempty"`
}

// IsmPolicyState is the policy state object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/State.kt
type IsmPolicyState struct {
	Name        string                     `json:"name"`
	Actions     []map[string]any           `json:"actions,omitempty"`
	Transitions []IsmPolicyStateTransition `json:"transitions,omitempty"`
}

// IsmErrorNotificationDestination is the error notification destination object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/destination/Destination.kt
type IsmErrorNotificationDestination struct {
	Type          string                                        `json:"type"`
	Chime         *IsmErrorNotificationDestinationChime         `json:"chime,omitempty"`
	Slack         *IsmErrorNotificationDestinationSlack         `json:"slack,omitempty"`
	CustomWebhook *IsmErrorNotificationDestinationCustomWebhook `json:"custom_webhook,omitempty"`
}

// IsmErrorNotificationDestinationChime is the chime notification object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/destination/Chime.kt
type IsmErrorNotificationDestinationChime struct {
	Url string `json:"url"`
}

// IsmErrorNotificationDestinationCustomWebhook is the webhook notification object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/destination/CustomWebhook.kt
type IsmErrorNotificationDestinationCustomWebhook struct {
	Url          *string           `json:"url,omitempty"`
	Scheme       *string           `json:"scheme,omitempty"`
	Host         *string           `json:"host,omitempty"`
	Port         *int64            `json:"port,omitempty"`
	Path         *string           `json:"path,omitempty"`
	QueryParams  map[string]string `json:"query_params,omitempty"`
	HeaderParams map[string]string `json:"header_params,omitempty"`
	Username     *string           `json:"username,omitempty"`
	Password     *string           `json:"password,omitempty"`
}

// IsmErrorNotificationDestinationSlack is the slack notification object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/destination/Slack.kt
type IsmErrorNotificationDestinationSlack struct {
	Url string `json:"url"`
}

// IsmErrorNotificationChannel is the channel object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/common/model/notification/Channel.kt
type IsmErrorNotificationChannel struct {
	ID string `json:"id"`
}

// IsmErrorNotificationMessageTemplate is the message template object
// Source: https://github.com/opensearch-project/OpenSearch/blob/main/server/src/main/java/org/opensearch/script/Script.java
type IsmErrorNotificationMessageTemplate struct {
	ScriptType string            `json:"type"`
	Lang       string            `json:"lang"`
	IdOrCode   string            `json:"idOrCode"`
	Options    map[string]string `json:"options,omitempty"`
	Params     map[string]string `json:"params,omitempty"`
}

// IsmPolicyStateTransition is the state transition object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/Transition.kt
type IsmPolicyStateTransition struct {
	StateName  string         `json:"state_name"`
	Conditions map[string]any `json:"conditions,omitempty"`
}

// IsmPolicyTemplate is the policy template object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/indexstatemanagement/model/ISMTemplate.kt
type IsmPolicyTemplate struct {
	IndexPatterns   []string `json:"index_patterns,omitempty"`
	Priority        *int64   `json:"priority,omitempty"`
	LastUpdatedTime *int64   `json:"last_updated_time,omitempty"`
}
