// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// ErrMissingIndex is returned e.g. from DeleteService if the index is missing.
	ErrMissingIndex = errors.New("elastic: index is missing")

	// ErrMissingType is returned e.g. from DeleteService if the type is missing.
	ErrMissingType = errors.New("elastic: type is missing")

	// ErrMissingId is returned e.g. from DeleteService if the document identifier is missing.
	ErrMissingId = errors.New("elastic: id is missing")
)

// checkResponse will return an error if the request/response indicates
// an error returned from Elasticsearch.
//
// HTTP status codes between in the range [200..299] are considered successful.
// All other errors are considered errors except they are specified in
// ignoreErrors. This is necessary because for some services, HTTP status 404
// is a valid response from Elasticsearch (e.g. the Exists service).
//
// The func tries to parse error details as returned from Elasticsearch
// and encapsulates them in type elastic.Error.
func checkResponse(req *http.Request, res *http.Response, ignoreErrors ...int) error {
	// 200-299 and 404 are valid status codes
	if (res.StatusCode >= 200 && res.StatusCode <= 299) || res.StatusCode == http.StatusNotFound {
		return nil
	}
	// Ignore certain errors?
	for _, code := range ignoreErrors {
		if code == res.StatusCode {
			return nil
		}
	}
	if res.Body == nil {
		return fmt.Errorf("elastic: Error %d (%s)", res.StatusCode, http.StatusText(res.StatusCode))
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("elastic: Error %d (%s) when reading body: %v", res.StatusCode, http.StatusText(res.StatusCode), err)
	}
	return createResponseError(res.StatusCode, slurp)
}

// createResponseError creates an Error structure from the HTTP response,
// its status code and the error information sent by Elasticsearch.
func createResponseError(statusCode int, data []byte) error {
	errReply := new(Error)
	err := json.Unmarshal(data, errReply)
	if err != nil {
		return fmt.Errorf("elastic: Error %d (%s)", statusCode, http.StatusText(statusCode))
	}
	if errReply != nil {
		if errReply.Status == 0 {
			errReply.Status = statusCode
		}
		return errReply
	}
	return fmt.Errorf("elastic: Error %d (%s)", statusCode, http.StatusText(statusCode))
}

// Error encapsulates error details as returned from Elasticsearch.
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("elastic: Error %d (%s): %s", e.Status, http.StatusText(e.Status), e.Message)
	} else {
		return fmt.Sprintf("elastic: Error %d (%s)", e.Status, http.StatusText(e.Status))
	}
}
