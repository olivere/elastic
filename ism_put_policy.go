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

// IsmPutPolicyService update a ISM policy by its name.
// See https://opensearch.org/docs/latest/im-plugin/ism/api/#create-policy
type IsmPutPolicyService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	name           string
	body           interface{}
	sequenceNumber *int64
	primaryTerm    *int64
}

// NewIsmPutPolicyService creates a new IsmPutPolicyService.
func NewIsmPutPolicyService(client *Client) *IsmPutPolicyService {
	return &IsmPutPolicyService{
		client: client,
	}
}

// Pretty tells Opensearch whether to return a formatted JSON response.
func (s *IsmPutPolicyService) Pretty(pretty bool) *IsmPutPolicyService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *IsmPutPolicyService) Human(human bool) *IsmPutPolicyService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *IsmPutPolicyService) ErrorTrace(errorTrace bool) *IsmPutPolicyService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *IsmPutPolicyService) FilterPath(filterPath ...string) *IsmPutPolicyService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *IsmPutPolicyService) Header(name string, value string) *IsmPutPolicyService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *IsmPutPolicyService) Headers(headers http.Header) *IsmPutPolicyService {
	s.headers = headers
	return s
}

// Name is name of the policy to create.
func (s *IsmPutPolicyService) Name(name string) *IsmPutPolicyService {
	s.name = name
	return s
}

// Body specifies the policy. Use a string or a type that will get serialized as JSON.
func (s *IsmPutPolicyService) Body(body interface{}) *IsmPutPolicyService {
	s.body = body
	return s
}

// SequenceNumber specifies the sequence number to update.
func (s *IsmPutPolicyService) SequenceNumber(seqNum int64) *IsmPutPolicyService {
	s.sequenceNumber = ptr.To[int64](seqNum)
	return s
}

// PrimaryTerm specifies the primary term to update.
func (s *IsmPutPolicyService) PrimaryTerm(primaryTerm int64) *IsmPutPolicyService {
	s.primaryTerm = ptr.To[int64](primaryTerm)
	return s
}

// buildURL builds the URL for the operation.
func (s *IsmPutPolicyService) buildURL() (string, url.Values, error) {
	var (
		path string
		err  error
	)

	// Build URL
	if s.primaryTerm != nil && s.sequenceNumber != nil {
		path, err = uritemplates.Expand("/_plugins/_ism/policies/{name}?if_seq_no={seqNum}&if_primary_term={priTerm}", map[string]string{
			"name":    s.name,
			"seqNum":  strconv.FormatInt(*s.sequenceNumber, 10),
			"priTerm": strconv.FormatInt(*s.primaryTerm, 10),
		})
	} else {
		path, err = uritemplates.Expand("/_plugins/_ism/policies/{name}", map[string]string{
			"name": s.name,
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
func (s *IsmPutPolicyService) Validate() error {
	var invalid []string
	if s.name == "" {
		invalid = append(invalid, "Name")
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
func (s *IsmPutPolicyService) Do(ctx context.Context) (*IsmGetPolicyResponse, error) {
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
	ret := new(IsmGetPolicyResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type IsmPutPolicy struct {
	Policy IsmPolicy `json:"policy"`
}
