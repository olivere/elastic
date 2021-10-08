// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestPercolatorQuery(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{
			"message": "A new bonsai tree in the office",
		})
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"document":{"message":"A new bonsai tree in the office"},"field":"query"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPercolatorQueryWithMultipleDocuments(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		Document(
			map[string]interface{}{
				"message": "bonsai tree",
			}, map[string]interface{}{
				"message": "new tree",
			},
		)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"documents":[{"message":"bonsai tree"},{"message":"new tree"}],"field":"query"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPercolatorQueryWithExistingDocument(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		IndexedDocumentIndex("my-index").
		IndexedDocumentType("_doc").
		IndexedDocumentId("2").
		IndexedDocumentVersion(1)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"field":"query","id":"2","index":"my-index","type":"_doc","version":1}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPercolatorQueryWithDetails(t *testing.T) {
	q := NewPercolatorQuery().
		Field("query").
		Document(map[string]interface{}{
			"message": "A new bonsai tree in the office",
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
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"percolate":{"document":{"message":"A new bonsai tree in the office"},"field":"query","id":"1","index":"index","preference":"one","routing":"route","version":1}}`
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
