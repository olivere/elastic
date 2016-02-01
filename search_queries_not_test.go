// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestNotQuery(t *testing.T) {
	f := NewNotQuery(NewTermQuery("user", "olivere"))
	src, err := f.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"not":{"query":{"term":{"user":"olivere"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestNotQueryWithParams(t *testing.T) {
	postDateFilter := NewRangeQuery("postDate").Gte("2010-03-01").Lte("2010-04-01")
	f := NewNotQuery(postDateFilter)
	f = f.QueryName("MyQueryName")
	src, err := f.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"not":{"_name":"MyQueryName","query":{"range":{"postDate":{"gte":"2010-03-01","lte":"2010-04-01"}}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
