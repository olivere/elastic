package elastic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Response represents a response from Elasticsearch.
type Response struct {
	// StatusCode is the HTTP status code, e.g. 200.
	StatusCode int
	// Header is the HTTP header from the HTTP response.
	// Keys in the map are canonicalized (see http.CanonicalHeaderKey).
	Header http.Header
	// Body is the deserialized response body.
	Body json.RawMessage
}

// newResponse creates a new response from the HTTP response.
func (c *Client) newResponse(res *http.Response) (*Response, error) {
	r := &Response{
		StatusCode: res.StatusCode,
		Header:     res.Header,
	}
	if res.Body != nil {
		slurp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		// HEAD requests return a body but no content
		if len(slurp) > 0 {
			if err := c.decoder.Decode(slurp, &r.Body); err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}
