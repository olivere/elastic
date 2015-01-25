// Copyright 2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

type OptimizeService struct {
	client             *Client
	indices            []string
	maxNumSegments     *int
	onlyExpungeDeletes *bool
	flush              *bool
	waitForMerge       *bool
	force              *bool
	pretty             bool
	debug              bool
}

func NewOptimizeService(client *Client) *OptimizeService {
	builder := &OptimizeService{
		client:  client,
		indices: make([]string, 0),
	}
	return builder
}

func (s *OptimizeService) Index(index string) *OptimizeService {
	s.indices = append(s.indices, index)
	return s
}

func (s *OptimizeService) Indices(indices ...string) *OptimizeService {
	s.indices = append(s.indices, indices...)
	return s
}

func (s *OptimizeService) MaxNumSegments(maxNumSegments int) *OptimizeService {
	s.maxNumSegments = &maxNumSegments
	return s
}

func (s *OptimizeService) OnlyExpungeDeletes(onlyExpungeDeletes bool) *OptimizeService {
	s.onlyExpungeDeletes = &onlyExpungeDeletes
	return s
}

func (s *OptimizeService) Flush(flush bool) *OptimizeService {
	s.flush = &flush
	return s
}

func (s *OptimizeService) WaitForMerge(waitForMerge bool) *OptimizeService {
	s.waitForMerge = &waitForMerge
	return s
}

func (s *OptimizeService) Force(force bool) *OptimizeService {
	s.force = &force
	return s
}

func (s *OptimizeService) Pretty(pretty bool) *OptimizeService {
	s.pretty = pretty
	return s
}

func (s *OptimizeService) Debug(debug bool) *OptimizeService {
	s.debug = debug
	return s
}

func (s *OptimizeService) Do() (*OptimizeResult, error) {
	// Build url
	urls := "/"

	// Indices part
	indexPart := make([]string, 0)
	for _, index := range s.indices {
		index, err := uritemplates.Expand("{index}", map[string]string{
			"index": index,
		})
		if err != nil {
			return nil, err
		}
		indexPart = append(indexPart, index)
	}
	if len(indexPart) > 0 {
		urls += strings.Join(indexPart, ",")
	}

	urls += "/_optimize"

	// Parameters
	params := make(url.Values)
	if s.maxNumSegments != nil {
		params.Set("max_num_segments", fmt.Sprintf("%d", *s.maxNumSegments))
	}
	if s.onlyExpungeDeletes != nil {
		params.Set("only_expunge_deletes", fmt.Sprintf("%v", *s.onlyExpungeDeletes))
	}
	if s.flush != nil {
		params.Set("flush", fmt.Sprintf("%v", *s.flush))
	}
	if s.waitForMerge != nil {
		params.Set("wait_for_merge", fmt.Sprintf("%v", *s.waitForMerge))
	}
	if s.force != nil {
		params.Set("force", fmt.Sprintf("%v", *s.force))
	}
	if s.pretty {
		params.Set("pretty", fmt.Sprintf("%v", s.pretty))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	// Set up a new request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	if s.debug {
		s.client.dumpRequest((*http.Request)(req))
	}

	// Get response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if s.debug {
		s.client.dumpResponse(res)
	}

	ret := new(OptimizeResult)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// -- Result of an optimize request.

type OptimizeResult struct {
	Shards shardsInfo `json:"_shards,omitempty"`
}
