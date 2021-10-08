// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestTasksListBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		TaskId   []string
		Expected string
	}{
		{
			TaskId:   []string{},
			Expected: "/_tasks",
		},
		{
			TaskId:   []string{"node1:42"},
			Expected: "/_tasks/node1%3A42",
		},
		{
			TaskId:   []string{"node1:42", "node2:37"},
			Expected: "/_tasks/node1%3A42%2Cnode2%3A37",
		},
	}

	for i, tt := range tests {
		path, _, err := client.TasksList().
			TaskId(tt.TaskId...).
			buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != tt.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, tt.Expected, path)
		}
	}
}

func TestTasksList(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))
	res, err := client.TasksList().
		Pretty(true).
		Human(true).
		Header("X-Opaque-Id", "123456").
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("response is nil")
	}
	if len(res.Nodes) == 0 {
		t.Fatalf("expected at least 1 node; got: %d", len(res.Nodes))
	}
	if want, have := "123456", res.Header.Get("X-Opaque-Id"); want != have {
		t.Fatalf("expected HTTP header %#v; got: %#v", want, have)
	}
	for _, node := range res.Nodes {
		if len(node.Tasks) == 0 {
			t.Fatalf("expected at least 1 task; got: %d", len(node.Tasks))
		}
		for _, task := range node.Tasks {
			have, found := task.Headers["X-Opaque-Id"]
			if !found {
				t.Logf("shaky test: expected to find headers[%q]", "X-Opaque-Id")
			} else if want := "123456"; want != have {
				t.Fatalf("expected headers[%q]=%q; got: %q", "X-Opaque-Id", want, have)
			}
		}
	}
}
