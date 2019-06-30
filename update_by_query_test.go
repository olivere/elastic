// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestUpdateByQueryBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Indices   []string
		Types     []string
		Expected  string
		ExpectErr bool
	}{
		{
			[]string{},
			[]string{},
			"",
			true,
		},
		{
			[]string{"index1"},
			[]string{},
			"/index1/_update_by_query",
			false,
		},
		{
			[]string{"index1", "index2"},
			[]string{},
			"/index1%2Cindex2/_update_by_query",
			false,
		},
		{
			[]string{},
			[]string{"type1"},
			"",
			true,
		},
		{
			[]string{"index1"},
			[]string{"type1"},
			"/index1/type1/_update_by_query",
			false,
		},
		{
			[]string{"index1", "index2"},
			[]string{"type1", "type2"},
			"/index1%2Cindex2/type1%2Ctype2/_update_by_query",
			false,
		},
	}

	for i, test := range tests {
		builder := client.UpdateByQuery().Index(test.Indices...).Type(test.Types...)
		err := builder.Validate()
		if err != nil {
			if !test.ExpectErr {
				t.Errorf("case #%d: %v", i+1, err)
				continue
			}
		} else {
			// err == nil
			if test.ExpectErr {
				t.Errorf("case #%d: expected error", i+1)
				continue
			}
			path, _, _ := builder.buildURL()
			if path != test.Expected {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
			}
		}
	}
}

func TestUpdateByQueryBodyWithQuery(t *testing.T) {
	client := setupTestClient(t)
	out, err := client.UpdateByQuery().Query(NewTermQuery("user", "olivere")).getBody()
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"query":{"term":{"user":"olivere"}}}`
	if got != want {
		t.Fatalf("\ngot  %s\nwant %s", got, want)
	}
}

func TestUpdateByQueryBodyWithQueryAndScript(t *testing.T) {
	client := setupTestClient(t)
	out, err := client.UpdateByQuery().
		Query(NewTermQuery("user", "olivere")).
		Script(NewScriptInline("ctx._source.likes++")).
		getBody()
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"query":{"term":{"user":"olivere"}},"script":{"source":"ctx._source.likes++"}}`
	if got != want {
		t.Fatalf("\ngot  %s\nwant %s", got, want)
	}
}

func TestUpdateByQuery(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))
	esversion, err := client.ElasticsearchVersion(DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	if esversion < "2.3.0" {
		t.Skipf("Elasticsearch %v does not support update-by-query yet", esversion)
	}

	sourceCount, err := client.Count(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	res, err := client.UpdateByQuery(testIndexName).ProceedOnVersionConflict().Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("response is nil")
	}
	if res.Updated != sourceCount {
		t.Fatalf("expected %d; got: %d", sourceCount, res.Updated)
	}
}

func TestUpdateByQueryAsync(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))
	esversion, err := client.ElasticsearchVersion(DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	if esversion < "2.3.0" {
		t.Skipf("Elasticsearch %v does not support update-by-query yet", esversion)
	}

	sourceCount, err := client.Count(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	res, err := client.UpdateByQuery(testIndexName).
		ProceedOnVersionConflict().
		Slices("auto").
		DoAsync(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected result != nil")
	}
	if res.TaskId == "" {
		t.Errorf("expected a task id, got %+v", res)
	}

	tasksGetTask := client.TasksGetTask()
	taskStatus, err := tasksGetTask.TaskId(res.TaskId).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if taskStatus == nil {
		t.Fatal("expected task status result != nil")
	}
}

func TestUpdateByQueryConflict(t *testing.T) {
	fail := func(r *http.Request) (*http.Response, error) {
		body := `{
			"took": 3,
			"timed_out": false,
			"total": 1,
			"updated": 0,
			"deleted": 0,
			"batches": 1,
			"version_conflicts": 1,
			"noops": 0,
			"retries": {
			  "bulk": 0,
			  "search": 0
			},
			"throttled_millis": 0,
			"requests_per_second": -1,
			"throttled_until_millis": 0,
			"failures": [
			  {
				"index": "a",
				"type": "_doc",
				"id": "yjsmdGsBm363wfQmSbhj",
				"cause": {
				  "type": "version_conflict_engine_exception",
				  "reason": "[_doc][yjsmdGsBm363wfQmSbhj]: version conflict, current version [4] is different than the one provided [3]",
				  "index_uuid": "1rmL3mt8TimwshF-M1DxdQ",
				  "shard": "0",
				  "index": "a"
				},
				"status": 409
			  }
			]
		   }`
		return &http.Response{
			StatusCode:    http.StatusConflict,
			Body:          ioutil.NopCloser(strings.NewReader(body)),
			ContentLength: int64(len(body)),
		}, nil
	}

	// Run against a failing endpoint and see if PerformRequest
	// retries correctly.
	tr := &failingTransport{path: "/example/_update_by_query", fail: fail}
	httpClient := &http.Client{Transport: tr}
	client, err := NewClient(SetHttpClient(httpClient), SetHealthcheck(false))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.UpdateByQuery("example").ProceedOnVersionConflict().Do(context.TODO())
	if err != nil {
		t.Fatalf("mock should not be failed %+v", err)
	}
	if res.Took != 3 {
		t.Errorf("took should be 3, got %d", res.Took)
	}
	if res.Total != 1 {
		t.Errorf("total should be 1, got %d", res.Total)
	}
	if res.VersionConflicts != 1 {
		t.Errorf("total should be 1, got %d", res.VersionConflicts)
	}
	if len(res.Failures) != 1 {
		t.Errorf("failures length should be 1, got %d", len(res.Failures))
	}
	expected := bulkIndexByScrollResponseFailure{Index: "a", Type: "_doc", Id: "yjsmdGsBm363wfQmSbhj", Status: 409}
	if res.Failures[0] != expected {
		t.Errorf("failures should be %+v, got %+v", expected, res.Failures[0])
	}
}
