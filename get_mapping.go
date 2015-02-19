package elastic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

type GetMappingService struct {
	client *Client
	index  string
	pretty bool
	debug  bool
}

func NewGetMappingService(client *Client) *GetMappingService {
	builder := &GetMappingService{
		client: client,
	}
	return builder
}

func (b *GetMappingService) Get(index string) *GetMappingService {
	b.index = index
	return b
}

func (b *GetMappingService) Pretty(pretty bool) *GetMappingService {
	b.pretty = pretty
	return b
}

func (b *GetMappingService) Debug(debug bool) *GetMappingService {
	b.debug = debug
	return b
}

func (b *GetMappingService) Do() (string, error) {
	// Build url
	urls, err := uritemplates.Expand("/{index}/_mapping", map[string]string{
		"index": b.index,
	})
	if err != nil {
		return "", err
	}

	// Parameters
	params := make(url.Values)
	if b.pretty {
		params.Set("pretty", fmt.Sprintf("%v", b.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := b.client.NewRequest("GET", urls)
	if err != nil {
		return "", err
	}

	if b.debug {
		b.client.dumpRequest((*http.Request)(req))
	}

	// Get response
	res, err := b.client.c.Do((*http.Request)(req))
	if err != nil {
		return "", err
	}
	if err := checkResponse(res); err != nil {
		return "", err
	}
	defer res.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
