// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestXPackAsyncSearchLifecycle(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndAddDocs(t, SetURL("http://elastic:elastic@localhost:9210"), SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetURL("http://elastic:elastic@localhost:9210"))

	// Match all should return all documents
	resp, err := client.XPackAsyncSearchSubmit().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Size(100).
		Pretty(true).
		WaitForCompletionTimeout("10s"). // should be ready by then
		KeepOnCompletion(true).          // keep even after completion
		KeepAlive("2m").                 // keep for at least 2 minutes
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if want, have := false, resp.IsRunning; want != have {
		t.Errorf("expected IsRunning=%v; got %v", want, have)
	}
	if want, have := false, resp.IsPartial; want != have {
		t.Errorf("expected IsPartial=%v; got %v", want, have)
	}
	if resp.ID == "" {
		t.Error(`expected ID!=""`)
	}
	if resp.Response == nil {
		t.Fatal("expected Response; got nil")
	}
	if want, have := int64(3), resp.Response.TotalHits(); want != have {
		t.Errorf("expected TotalHits=%v; got %v", want, have)
	}
	for _, hit := range resp.Response.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get the search results with the given ID
	get, err := client.XPackAsyncSearchGet().ID(resp.ID).Pretty(true).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if get == nil {
		t.Fatal("expected response, got nil")
	}
	if want, have := false, get.IsRunning; want != have {
		t.Errorf("expected IsRunning=%v; got %v", want, have)
	}
	if want, have := false, get.IsPartial; want != have {
		t.Errorf("expected IsPartial=%v; got %v", want, have)
	}
	if want, have := get.ID, get.ID; want != have {
		t.Errorf("expected ID!=%q; got %q", want, have)
	}
	if get.Response == nil {
		t.Fatal("expected Response; got nil")
	}
	if want, have := int64(3), get.Response.TotalHits(); want != have {
		t.Errorf("expected TotalHits=%v; got %v", want, have)
	}

	// Delete the search results with the given ID
	del, err := client.XPackAsyncSearchDelete().ID(resp.ID).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if del == nil {
		t.Fatal("expected response, got nil")
	}
	if want, have := true, del.Acknowledged; want != have {
		t.Errorf("expected Acknowledged=%v; got %v", want, have)
	}
}
