// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestDeleteScript(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	scriptID := "example-delete-script-id"

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
	putRes, err := client.PutScript().
		Id(scriptID).
		BodyString(script).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if putRes == nil {
		t.Errorf("expected result to be != nil; got: %v", putRes)
	}
	if !putRes.Acknowledged {
		t.Errorf("expected ack for PutScript op; got %v", putRes.Acknowledged)
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

	// DeleteScript API
	res, err := client.DeleteScript().
		Id(scriptID).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Errorf("expected result to be != nil; got: %v", res)
	}
	if !res.Acknowledged {
		t.Errorf("expected ack for DeleteScript op; got %v", res.Acknowledged)
	}

	// Must not exist now
	_, err = client.PerformRequest(
		context.Background(),
		PerformRequestOptions{
			Method: "DELETE",
			Path:   "/_scripts/" + scriptID,
		})
	if err != nil && !IsNotFound(err) {
		t.Fatal(err)
	}
}
