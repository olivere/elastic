package elastic

import (
	"context"
	"fmt"
)

type service struct {
	client                 *Client
	waitForCompletion      *bool
	performInternalRequest func(context.Context) (*Response, error)
}

func newService(client *Client, request func(context.Context) (*Response, error)) *service {
	return &service{
		client:                 client,
		performInternalRequest: request,
	}
}

func (s *service) doAsync(ctx context.Context) (*StartTaskResult, error) {
	// DoAsync only makes sense with WaitForCompletion set to true
	if s.waitForCompletion != nil && *s.waitForCompletion {
		return nil, fmt.Errorf("cannot start a task with WaitForCompletion set to true")
	}
	f := false
	s.waitForCompletion = &f

	res, err := s.performInternalRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(StartTaskResult)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *service) Do(ctx context.Context) (*BulkIndexByScrollResponse, error) {
	// Check pre-conditions
	res, err := s.performInternalRequest(ctx)
	if err != nil {
		return nil, err
	}
	// Return operation response (BulkIndexByScrollResponse is defined in DeleteByQuery)
	ret := new(BulkIndexByScrollResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
