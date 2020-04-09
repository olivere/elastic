// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"errors"
	"testing"
	"time"
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

func TestTasksGetTaskWithError(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Create a reindexing task
	var taskID string
	{
		res, err := client.UpdateByQuery(testIndexName).
			WaitForCompletion(false).
			Conflicts("proceed").
			Script(NewScript("kaboom")).
			DoAsync(context.Background())
		if err != nil {
			t.Fatalf("unable to start update_by_query task: %v", err)
		}
		taskID = res.TaskId
	}

	var (
		response *TasksGetTaskResponse
		lastErr  error
	)
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		for {
			// Get the task by ID
			res, err := client.TasksGetTask().
				TaskId(taskID).
				Do(context.Background())
			if err != nil {
				lastErr = err
				return
			}
			if res == nil {
				lastErr = errors.New("response is nil")
				return
			}
			if res.Completed {
				lastErr = nil
				response = res
				return
			}
			time.Sleep(1 * time.Second) // retry
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("expected to finish the task after 5 seconds")
	}
	if lastErr != nil {
		t.Fatalf("expected no error, got %v", lastErr)
	}
	if response == nil {
		t.Fatal("expected a response, got nil")
	}
	if response.Error == nil {
		t.Fatal("expected a response with an error, got nil")
	}
	if want, have := "script_exception", response.Error.Type; want != have {
		t.Fatalf("expected an error type of %q, got %q", want, have)
	}
}
