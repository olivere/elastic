// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestWrapperQueryIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	tq := NewTermQuery("user", "olivere")
	src, err := tq.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	base64string := base64.StdEncoding.EncodeToString(data)

	q := NewWrapperQuery(base64string)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(q).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 2; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}
}
