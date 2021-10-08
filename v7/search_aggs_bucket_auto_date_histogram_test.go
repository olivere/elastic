// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestAutoDateHistogramAggregation(t *testing.T) {
	agg := NewAutoDateHistogramAggregation().
		Field("date").
		Buckets(10)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"auto_date_histogram":{"buckets":10,"field":"date"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestAutoDateHistogramAggregationWithFormat(t *testing.T) {
	agg := NewAutoDateHistogramAggregation().Field("date").Format("yyyy-MM-dd").Buckets(5)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"auto_date_histogram":{"buckets":5,"field":"date","format":"yyyy-MM-dd"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestAutoDateHistogramAggregationWithMissing(t *testing.T) {
	agg := NewAutoDateHistogramAggregation().Field("date").Missing("1900")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"auto_date_histogram":{"field":"date","missing":"1900"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
