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

// AlertingGetMonitorService get a monitor by its name.
// See https://opensearch.org/docs/latest/observing-your-data/alerting/api/#get-monitor
type AlertingGetMonitorService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name string
}

// NewAlertingGetMonitorService creates a new AlertingGetMonitorService.
func NewAlertingGetMonitorService(client *Client) *AlertingGetMonitorService {
	return &AlertingGetMonitorService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *AlertingGetMonitorService) Pretty(pretty bool) *AlertingGetMonitorService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *AlertingGetMonitorService) Human(human bool) *AlertingGetMonitorService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *AlertingGetMonitorService) ErrorTrace(errorTrace bool) *AlertingGetMonitorService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *AlertingGetMonitorService) FilterPath(filterPath ...string) *AlertingGetMonitorService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *AlertingGetMonitorService) Header(name string, value string) *AlertingGetMonitorService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *AlertingGetMonitorService) Headers(headers http.Header) *AlertingGetMonitorService {
	s.headers = headers
	return s
}

// Name is name of the monitor to get.
func (s *AlertingGetMonitorService) Name(name string) *AlertingGetMonitorService {
	s.name = name
	return s
}

// buildURL builds the URL for the operation.
func (s *AlertingGetMonitorService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_plugins/_alerting/monitors/{name}", map[string]string{
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
func (s *AlertingGetMonitorService) Validate() error {
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
func (s *AlertingGetMonitorService) Do(ctx context.Context) (*AlertingGetMonitorResponse, error) {
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
	ret := new(AlertingGetMonitorResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// AlertingGetMonitorResponse is the get monitor response object
type AlertingGetMonitorResponse struct {
	Id             string             `json:"_id"`
	Version        int64              `json:"_version"`
	SequenceNumber int64              `json:"_seq_no"`
	PrimaryTerm    int64              `json:"_primary_term"`
	Monitor        AlertingGetMonitor `json:"monitor"`
}

// AlertingMonitorBase is the base monitor object
type AlertingMonitor struct {
	Type        string           `json:"type"`
	Name        string           `json:"name"`
	MonitorType string           `json:"monitor_type"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Schedule    map[string]any   `json:"schedule"`
	Inputs      []map[string]any `json:"inputs"`
	Triggers    []map[string]any `json:"triggers"`
}

// AlertingGetMonitor is the ISM policy
type AlertingGetMonitor struct {
	AlertingMonitor `json:",inline"`
	EnabledTime     *int64 `json:"enabled_time,omitempty"`
	LastUpdatedTime *int64 `json:"last_updated_time,omitempty"`
}
