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
	"k8s.io/utils/ptr"
)

// AlertingPutMonitorService update a monitor by its name.
// See https://opensearch.org/docs/latest/observing-your-data/alerting/api/#update-monitor
type AlertingPutMonitorService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	id             string
	body           interface{}
	sequenceNumber *int64
	primaryTerm    *int64
}

// NewAlertingPutMonitorService creates a new AlertingPutMonitorService.
func NewAlertingPutMonitorService(client *Client) *AlertingPutMonitorService {
	return &AlertingPutMonitorService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *AlertingPutMonitorService) Pretty(pretty bool) *AlertingPutMonitorService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *AlertingPutMonitorService) Human(human bool) *AlertingPutMonitorService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *AlertingPutMonitorService) ErrorTrace(errorTrace bool) *AlertingPutMonitorService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *AlertingPutMonitorService) FilterPath(filterPath ...string) *AlertingPutMonitorService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *AlertingPutMonitorService) Header(name string, value string) *AlertingPutMonitorService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *AlertingPutMonitorService) Headers(headers http.Header) *AlertingPutMonitorService {
	s.headers = headers
	return s
}

// Id is id of the monitor to update.
func (s *AlertingPutMonitorService) Id(id string) *AlertingPutMonitorService {
	s.id = id
	return s
}

// Body specifies the policy. Use a string or a type that will get serialized as JSON.
func (s *AlertingPutMonitorService) Body(body interface{}) *AlertingPutMonitorService {
	s.body = body
	return s
}

// SequenceNumber specifies the sequence number to update.
func (s *AlertingPutMonitorService) SequenceNumber(seqNum int64) *AlertingPutMonitorService {
	s.sequenceNumber = ptr.To[int64](seqNum)
	return s
}

// PrimaryTerm specifies the primary term to update.
func (s *AlertingPutMonitorService) PrimaryTerm(primaryTerm int64) *AlertingPutMonitorService {
	s.primaryTerm = ptr.To[int64](primaryTerm)
	return s
}

// buildURL builds the URL for the operation.
func (s *AlertingPutMonitorService) buildURL() (string, url.Values, error) {
	var (
		path string
		err  error
	)

	// Build URL
	if s.primaryTerm != nil && s.sequenceNumber != nil {
		path, err = uritemplates.Expand("/_plugins/_alerting/monitors/{id}?if_seq_no={seqNum}&if_primary_term={priTerm}", map[string]string{
			"id":      s.id,
			"seqNum":  strconv.FormatInt(*s.sequenceNumber, 10),
			"priTerm": strconv.FormatInt(*s.primaryTerm, 10),
		})
	} else {
		path, err = uritemplates.Expand("/_plugins/_alerting/monitors/{id}", map[string]string{
			"id": s.id,
		})
	}

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
func (s *AlertingPutMonitorService) Validate() error {
	var invalid []string
	if s.id == "" {
		invalid = append(invalid, "Id")
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
func (s *AlertingPutMonitorService) Do(ctx context.Context) (*AlertingGetMonitorResponse, error) {
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
	ret := new(AlertingGetMonitorResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
