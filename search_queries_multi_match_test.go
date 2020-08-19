// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestMultiMatchQuery(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryWithNoFields(t *testing.T) {
	q := NewMultiMatchQuery("accident")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":[],"query":"accident"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryBestFields(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("best_fields").TieBreaker(0)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":0,"type":"best_fields"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryMostFields(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("most_fields").TieBreaker(1)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":1,"type":"most_fields"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryCrossFields(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("cross_fields").TieBreaker(0)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":0,"type":"cross_fields"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryPhrase(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("phrase").TieBreaker(0)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":0,"type":"phrase"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryPhrasePrefix(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("phrase_prefix").TieBreaker(0)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":0,"type":"phrase_prefix"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryBestFieldsWithCustomTieBreaker(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").
		Type("best_fields").
		TieBreaker(0.3)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":0.3,"type":"best_fields"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryBoolPrefix(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("bool_prefix").TieBreaker(1)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","tie_breaker":1,"type":"bool_prefix"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiMatchQueryOptions(t *testing.T) {
	tests := []struct {
		Query *MultiMatchQuery
		Want  string
	}{
		{
			Query: NewMultiMatchQuery("this is a test", "message"),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test"}}`,
		},
		{
			Query: NewMultiMatchQuery("this is a test", "message").Type("best_fields").TieBreaker(0),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test","tie_breaker":0,"type":"best_fields"}}`,
		},
		{
			Query: NewMultiMatchQuery("this is a test", "message").Type("best_fields").TieBreaker(0.5),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test","tie_breaker":0.5,"type":"best_fields"}}`,
		},
		{
			Query: NewMultiMatchQuery("this is a test", "message").TieBreaker(0.5).Type("best_fields"),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test","tie_breaker":0.5,"type":"best_fields"}}`,
		},
		{
			Query: NewMultiMatchQuery("this is a test", "message").Type("cross_fields").TieBreaker(0.5),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test","tie_breaker":0.5,"type":"cross_fields"}}`,
		},
		{
			Query: NewMultiMatchQuery("this is a test", "message").TieBreaker(0.5).Type("cross_fields"),
			Want:  `{"multi_match":{"fields":["message"],"query":"this is a test","tie_breaker":0.5,"type":"cross_fields"}}`,
		},
	}
	for i, tt := range tests {
		src, err := tt.Query.Source()
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		data, err := json.Marshal(src)
		if err != nil {
			t.Fatalf("#%d: marshaling to JSON failed: %v", i, err)
		}
		if want, have := tt.Want, string(data); want != have {
			t.Errorf("#%d: want\n%s\n, have:\n%s", i, want, have)
		}
	}
}

func TestMultiMatchQueryDefaultTieBreaker(t *testing.T) {
	// should not add tie_breaker field within the query DSL unless specified
	q := NewMultiMatchQuery("this is a test", "subject", "message").Type("bool_prefix")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test","type":"bool_prefix"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
