// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type CountService struct {
	client  *Client
	indices []string
}

func NewCountService(client *Client) *CountService {
	builder := &CountService{
		client: client,
	}
	return builder
}

func (b *CountService) Index(index string) *CountService {
	if b.indices == nil {
		b.indices = make([]string, 0)
	}
	b.indices = append(b.indices, index)
	return b
}

func (b *CountService) Indices(indices ...string) *CountService {
	if b.indices == nil {
		b.indices = make([]string, 0)
	}
	b.indices = append(b.indices, indices...)
	return b
}
