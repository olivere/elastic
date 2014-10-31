// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkResponse(res *http.Response) error {
	// 200-299 and 404 are valid status codes
	if (res.StatusCode >= 200 && res.StatusCode <= 299) ||
		res.StatusCode == http.StatusNotFound {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("elastic: got HTTP response code %d and error reading body: %v", res.StatusCode, err)
	}
	errReply := new(Error)
	err = json.Unmarshal(slurp, errReply)
	if err == nil && errReply != nil {
		return errReply.Error()
	}
	return nil
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

func (e *Error) Error() error {
	return fmt.Errorf("elastic: Error %d: %s", e.Status, e.Message)
}
