// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRangeFilter(t *testing.T) {
	f := NewRangeFilter("postDate").From("2010-03-01").To("2010-04-01")
	f = f.FilterName("MyFilterName")
	f = f.Execution("index")
	src, err := f.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"range":{"_name":"MyFilterName","execution":"index","postDate":{"from":"2010-03-01","include_lower":true,"include_upper":true,"to":"2010-04-01"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRangeFilterWithTimeZone(t *testing.T) {
	f := NewRangeFilter("born").
		Gte("2012-01-01").
		Lte("now").
		TimeZone("+1:00")
	src, err := f.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"range":{"born":{"from":"2012-01-01","include_lower":true,"include_upper":true,"time_zone":"+1:00","to":"now"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
