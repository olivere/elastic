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

// AlertingSearchMonitorService get a monitor by its name.
// See https://opensearch.org/docs/latest/observing-your-data/alerting/api/#search-monitors
type AlertingSearchMonitorService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	search interface{}
}

// NewAlertingSearchMonitorService creates a new AlertingSearchMonitorService.
func NewAlertingSearchMonitorService(client *Client) *AlertingSearchMonitorService {
	return &AlertingSearchMonitorService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *AlertingSearchMonitorService) Pretty(pretty bool) *AlertingSearchMonitorService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *AlertingSearchMonitorService) Human(human bool) *AlertingSearchMonitorService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *AlertingSearchMonitorService) ErrorTrace(errorTrace bool) *AlertingSearchMonitorService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *AlertingSearchMonitorService) FilterPath(filterPath ...string) *AlertingSearchMonitorService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *AlertingSearchMonitorService) Header(name string, value string) *AlertingSearchMonitorService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *AlertingSearchMonitorService) Headers(headers http.Header) *AlertingSearchMonitorService {
	s.headers = headers
	return s
}

// Name is name of the monitor to get.
func (s *AlertingSearchMonitorService) SearchByName(name string) *AlertingSearchMonitorService {
	s.search = map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"monitor.name": name,
			},
		},
	}
	return s
}

// Body specifies the search. Use a string or a type that will get serialized as JSON.
func (s *AlertingSearchMonitorService) Search(search interface{}) *AlertingSearchMonitorService {
	s.search = search
	return s
}

// buildURL builds the URL for the operation.
func (s *AlertingSearchMonitorService) buildURL() (string, url.Values, error) {
	var (
		path string
		err  error
	)

	// Build URL
	path, err = uritemplates.Expand("/_plugins/_alerting/monitors/_search", map[string]string{})
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
func (s *AlertingSearchMonitorService) Validate() error {
	var invalid []string
	if s.search == nil {
		invalid = append(invalid, "Search")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *AlertingSearchMonitorService) Do(ctx context.Context) ([]AlertingSearchMonitor, error) {
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
		Body:    s.search,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(AlertingSearchMonitorResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}

	return ret.Hits.Hits, nil
}

type AlertingSearchMonitorResponse struct {
	Hits AlertingSearchMonitorHit `json:"hits"`
}

type AlertingSearchMonitorHit struct {
	Hits []AlertingSearchMonitor `json:"hits"`
}

type AlertingSearchMonitor struct {
	Id      string             `json:"_id"`
	Monitor AlertingGetMonitor `json:"_source"`
}
