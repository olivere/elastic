// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestPutScript(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	scriptID := "example-put-script-id"

	// Ensure the script does not exist
	_, err := client.PerformRequest(
		context.Background(),
		PerformRequestOptions{
			Method: "DELETE",
			Path:   "/_scripts/" + scriptID,
		})
	if err != nil && !IsNotFound(err) {
		t.Fatal(err)
	}

	// PutScript API
	script := `{
		"script": {
			"lang": "painless",
			"source": "ctx._source.message = params.new_message"
		}
	}`
	res, err := client.PutScript().
		Id(scriptID).
		BodyString(script).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Errorf("expected result to be != nil; got: %v", res)
	}
	if !res.Acknowledged {
		t.Errorf("expected ack for PutScript op; got %v", res.Acknowledged)
	}

	// Must exist now
	_, err = client.PerformRequest(
		context.Background(),
		PerformRequestOptions{
			Method: "GET",
			Path:   "/_scripts/" + scriptID,
		})
	if err != nil {
		t.Fatal(err)
	}
	// Cleanup
	client.PerformRequest(
		context.Background(),
		PerformRequestOptions{
			Method: "DELETE",
			Path:   "/_scripts/" + scriptID,
		})
}
