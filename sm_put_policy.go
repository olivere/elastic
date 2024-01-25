package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/disaster37/opensearch/v2/uritemplates"
)

// SmPutPolicyService update a SM policy by its name.
// See https://opensearch.org/docs/latest/tuning-your-cluster/availability-and-recovery/snapshots/sm-api/#create-or-update-a-policy
type SmPutPolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name           string
	body           interface{}
	sequenceNumber int64
	primaryTerm    int64
}

// NewSmPutPolicyService creates a new SmPutPolicyService.
func NewSmPutPolicyService(client *Client) *SmPutPolicyService {
	return &SmPutPolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *SmPutPolicyService) Pretty(pretty bool) *SmPutPolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SmPutPolicyService) Human(human bool) *SmPutPolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SmPutPolicyService) ErrorTrace(errorTrace bool) *SmPutPolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SmPutPolicyService) FilterPath(filterPath ...string) *SmPutPolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SmPutPolicyService) Header(name string, value string) *SmPutPolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SmPutPolicyService) Headers(headers http.Header) *SmPutPolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to create.
func (s *SmPutPolicyService) Name(name string) *SmPutPolicyService {
	s.name = name
	return s
}

// Body specifies the policy. Use a string or a type that will get serialized as JSON.
func (s *SmPutPolicyService) Body(body interface{}) *SmPutPolicyService {
	s.body = body
	return s
}

// SequenceNumber specifies the sequence number to update.
func (s *SmPutPolicyService) SequenceNumber(seqNum int64) *SmPutPolicyService {
	s.sequenceNumber = seqNum
	return s
}

// PrimaryTerm specifies the primary term to update.
func (s *SmPutPolicyService) PrimaryTerm(primaryTerm int64) *SmPutPolicyService {
	s.primaryTerm = primaryTerm
	return s
}

// buildURL builds the URL for the operation.
func (s *SmPutPolicyService) buildURL() (string, url.Values, error) {

	// Build URL

	path, err := uritemplates.Expand("/_plugins/_sm/policies/{name}?if_seq_no={seqNum}&if_primary_term={priTerm}", map[string]string{
		"name":    s.name,
		"seqNum":  strconv.FormatInt(s.sequenceNumber, 10),
		"priTerm": strconv.FormatInt(s.primaryTerm, 10),
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
func (s *SmPutPolicyService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
	}
	if s.body == nil {
		invalid = append(invalid, "Body")
	}
	if s.primaryTerm == 0 {
		invalid = append(invalid, "PrimaryTerm")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SmPutPolicyService) Do(ctx context.Context) (*SmGetPolicyResponse, error) {
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
	ret := new(SmGetPolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SmPutPolicy is object to put when you should to create or update snapshot policy
type SmPutPolicy struct {
	Description    *string                `json:"description,omitempty"`
	Enabled        *bool                  `json:"enabled,omitempty"`
	SnapshotConfig SmPolicySnapshotConfig `json:"snapshot_config"`
	Creation       SmPolicyCreation       `json:"creation"`
	Deletion       *SmPolicyDeletion      `json:"deletion,omitempty"`
	Notification   *SmPolicyNotification  `json:"notification,omitempty"`
}

// SmPolicySnapshotConfig is the snapshot config
type SmPolicySnapshotConfig struct {
	DateFormat         *string        `json:"date_format,omitempty"`
	Timezone           *string        `json:"timezone,omitempty"`
	Indices            *string        `json:"indices,omitempty"`
	Repository         string         `json:"repository,omitempty"`
	IgnoreUnavailable  *bool          `json:"ignore_unavailable,omitempty"`
	IncludeGlobalState *bool          `json:"include_global_state,omitempty"`
	Partial            *bool          `json:"partial,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

// SmPolicyCreation is the creation object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/snapshotmanagement/model/SMPolicy.kt#L252
type SmPolicyCreation struct {
	Schedule  map[string]any `json:"schedule"`
	TimeLimit *string        `json:"time_limit,omitempty"`
}

// SmPolicyDeletion is the deletion object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/snapshotmanagement/model/SMPolicy.kt#L300
type SmPolicyDeletion struct {
	Schedule  map[string]any           `json:"schedule,omitempty"`
	Condition *SmPolicyDeleteCondition `json:"condition,omitempty"`
	TimeLimit *string                  `json:"time_limit,omitempty"`
}

// SmPolicyDeleteCondition is the delete condition object
type SmPolicyDeleteCondition struct {
	MaxCount *int64  `json:"max_count,omitempty"`
	MaxAge   *string `json:"max_age,omitempty"`
	MinCount *int64  `json:"min_count,omitempty"`
}

// SmPolicyNotification is the notification object
// Source: https://github.com/opensearch-project/index-management/blob/main/src/main/kotlin/org/opensearch/indexmanagement/snapshotmanagement/model/NotificationConfig.kt#L28
type SmPolicyNotification struct {
	Channel    SmPolicyNotificationChannel    `json:"channel"`
	Conditions *SmPolicyNotificationCondition `json:"conditions,omitempty"`
}

// SmPolicyNotificationChannel is the channel object
type SmPolicyNotificationChannel struct {
	ID string `json:"id"`
}

// SmPolicyNotificationCondition is the notification condition object
type SmPolicyNotificationCondition struct {
	Creation          *bool `json:"creation,omitempty"`
	Deletion          *bool `json:"deletion,omitempty"`
	Failure           *bool `json:"failure,omitempty"`
	TimeLimitExceeded *bool `json:"time_limit_exceeded,omitempty"`
}
