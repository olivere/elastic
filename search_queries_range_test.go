// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRangeQuery(t *testing.T) {
	q := NewRangeQuery("postDate").Gte("2010-03-01").Lte("2010-04-01")
	q = q.QueryName("my_query")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"range":{"_name":"my_query","postDate":{"gte":"2010-03-01","lte":"2010-04-01"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithTimeZone(t *testing.T) {
	q := NewRangeQuery("born").
		Gte("2012-01-01").
		Lte("now").
		TimeZone("+1:00")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"range":{"born":{"gte":"2012-01-01","time_zone":"+1:00","lte":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeQueryWithFormat(t *testing.T) {
	q := NewRangeQuery("born").
		Gte("2012/01/01").
		Lte("now").
		Format("yyyy/MM/dd")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"range":{"born":{"format":"yyyy/MM/dd","gte":"2012/01/01","lte":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
