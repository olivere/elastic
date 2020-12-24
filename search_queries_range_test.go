// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRangeQuery(t *testing.T) {
	q := NewRangeQuery("postDate").
		Gte("2010-03-01")

	got := asJsonString(t, q)
	expected := `{"range":{"postDate":{"gte":"2010-03-01"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithBoost(t *testing.T) {
	q := NewRangeQuery("postDate").
		Gte("2010-03-01").
		Boost(3)

	got := asJsonString(t, q)
	expected := `{"range":{"postDate":{"boost":3,"gte":"2010-03-01"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithRelation(t *testing.T) {
	q := NewRangeQuery("postDate").
		Gte("2010-03-01").
		Relation(RelationWithin)

	got := asJsonString(t, q)
	expected := `{"range":{"postDate":{"gte":"2010-03-01","relation":"WITHIN"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithQueryName(t *testing.T) {
	q := NewRangeQuery("postDate").
		Gte("2010-03-01").
		QueryName("queryName")

	got := asJsonString(t, q)
	expected := `{"range":{"_name":"queryName","postDate":{"gte":"2010-03-01"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithTimeZone(t *testing.T) {
	q := NewRangeQuery("born").
		Gte("2012-01-01").
		Lte("now").
		TimeZone("+1:00")

	got := asJsonString(t, q)
	expected := `{"range":{"born":{"gte":"2012-01-01","lte":"now","time_zone":"+1:00"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithFormat(t *testing.T) {
	q := NewRangeQuery("born").
		Gte("2012/01/01").
		Lte("now").
		Format("yyyy/MM/dd")

	got := asJsonString(t, q)
	expected := `{"range":{"born":{"format":"yyyy/MM/dd","gte":"2012/01/01","lte":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDeprecatedRangeQueryTransformedIntoGtLt(t *testing.T) {
	q := NewRangeQuery("born").
		From("2012/01/01").
		To("now").
		IncludeLower(false).
		IncludeUpper(false)

	got := asJsonString(t, q)
	expected := `{"range":{"born":{"gt":"2012/01/01","lt":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDeprecatedRangeQueryTransformedIntoGteLte(t *testing.T) {
	q := NewRangeQuery("postDate").
		From("2012/01/01").
		To("now").
		IncludeLower(true).
		IncludeUpper(true)

	got := asJsonString(t, q)
	expected := `{"range":{"postDate":{"gte":"2012/01/01","lte":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func asJsonString(t *testing.T, q *RangeQuery) string {
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	return string(data)
}
