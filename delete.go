// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type DeleteService struct {
	client *Client
}

func NewDeleteService(client *Client) *DeleteService {
	builder := &DeleteService{
		client: client,
	}
	return builder
}
