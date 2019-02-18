// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestQueryStringQuery(t *testing.T) {
	q := NewQueryStringQuery(`this AND that OR thus`)
	q = q.DefaultField("content")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"query_string":{"default_field":"content","query":"this AND that OR thus"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestQueryStringQueryTimeZone(t *testing.T) {
	q := NewQueryStringQuery(`tweet_date:[2015-01-01 TO 2017-12-31]`)
	q = q.TimeZone("Europe/Berlin")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"query_string":{"query":"tweet_date:[2015-01-01 TO 2017-12-31]","time_zone":"Europe/Berlin"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestQueryStringQueryIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	q := NewQueryStringQuery("Golang")

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
	if got, want := searchResult.TotalHits(), int64(1); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 1; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}
}
