// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type DeleteByQueryService struct {
	client *Client
}

func NewDeleteByQueryService(client *Client) *DeleteByQueryService {
	builder := &DeleteByQueryService{
		client: client,
	}
	return builder
}
