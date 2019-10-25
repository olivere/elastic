// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestPercolatorQuery(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{
			"message": "Some message",
		})
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := jsoniter.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"document":{"message":"Some message"},"field":"query"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPercolatorQueryWithDetails(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{
			"message": "Some message",
		}).
		IndexedDocumentIndex("index").
		IndexedDocumentId("1").
		IndexedDocumentRouting("route").
		IndexedDocumentPreference("one").
		IndexedDocumentVersion(1)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := jsoniter.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"document":{"message":"Some message"},"field":"query","id":"1","index":"index","preference":"one","routing":"route","version":1}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPercolatorQueryWithMissingFields(t *testing.T) {
	q := NewPercolatorQuery() // no Field, Document, or Query
	_, err := q.Source()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
