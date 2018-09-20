// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestTasksGetTaskBuildURL(t *testing.T) {
	client := setupTestClient(t)

	// Get specific task
	got, _, err := client.TasksGetTask().TaskId("node1:123").buildURL()
	if err != nil {
		t.Fatal(err)
	}
	want := "/_tasks/node1%3A123"
	if got != want {
		t.Errorf("want %q; got %q", want, got)
	}

	// Get specific task
	got, _, err = client.TasksGetTask().TaskIdFromNodeAndId("node2", 678).buildURL()
	if err != nil {
		t.Fatal(err)
	}
	want = "/_tasks/node2%3A678"
	if got != want {
		t.Errorf("want %q; got %q", want, got)
	}
}

func TestTasksGetTask(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Create a reindexing task
	var taskID string
	{
		res, err := client.Reindex().
			SourceIndex(testIndexName).
			DestinationIndex(testIndexName4).
			DoAsync(context.Background())
		if err != nil {
			t.Fatalf("unable to start reindexing task: %v", err)
		}
		taskID = res.TaskId
	}

	// Get the task by ID
	res, err := client.TasksGetTask().
		TaskId(taskID).
		Header("X-Opaque-Id", "987654").
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("response is nil")
	}
	if want, have := "987654", res.Header.Get("X-Opaque-Id"); want != have {
		t.Fatalf("expected HTTP header %#v; got: %#v", want, have)
	}
	if res.Task == nil {
		t.Fatal("task is nil")
	}
	// Elasticsearch <= 6.4.1 doesn't return the X-Opaque-Id in the body,
	// only in response header.
	/*
		have, found := res.Task.Headers["X-Opaque-Id"]
		if !found {
			t.Fatalf("expected to find headers[%q]", "X-Opaque-Id")
		}
		if want := "987654"; want != have {
			t.Fatalf("expected headers[%q]=%q; got: %q", "X-Opaque-Id", want, have)
		}
	*/
}
