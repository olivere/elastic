// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestPointInTimeOpenAndClose(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Create a Point In Time
	openResp, err := client.OpenPointInTime(testIndexName).
		KeepAlive("1m").
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if openResp == nil {
		t.Fatal("expected non-nil Point In Time")
	}
	if openResp.Id == "" {
		t.Fatal("expected non-blank Point In Time ID")
	}

	// Close the Point in Time
	closeResp, err := client.ClosePointInTime(openResp.Id).Pretty(true).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if closeResp == nil {
		t.Fatal("expected non-nil Point In Time")
	}
	if want, have := true, closeResp.Succeeded; want != have {
		t.Fatalf("want Succeeded=%v, have %v", want, have)
	}
	if want, have := 1, closeResp.NumFreed; want != have {
		t.Fatalf("want NumFreed=%v, have %v", want, have)
	}
}

func TestPointInTimeLifecycle(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Create a Point In Time
	pitResp, err := client.OpenPointInTime().
		Index(testIndexName).
		KeepAlive("1m").
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pitResp == nil {
		t.Fatal("expected non-nil Point In Time")
	}
	if pitResp.Id == "" {
		t.Fatal("expected non-blank Point In Time ID")
	}

	// We remove the documents here, but will be able to still search with
	// the PIT previously created
	_, err = client.DeleteByQuery(testIndexName).
		Query(NewMatchAllQuery()).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Search with the Point in Time ID
	searchResult, err := client.Search().
		// Index(testIndexName). // <-- you may not use indices with PointInTime!
		Query(NewMatchAllQuery()).
		PointInTime(NewPointInTimeWithKeepAlive(pitResp.Id, "1m")).
		Size(100).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(3); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}

}
