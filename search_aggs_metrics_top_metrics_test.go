// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestTopMetricsAggregation(t *testing.T) {
	agg := NewTopMetricsAggregation().
		Sort("f1", false).
		Field("a").
		Field("b").
		Size(3)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"top_metrics":{"metrics":[{"field":"a"},{"field":"b"}],"size":3,"sort":{"f1":{"order":"desc"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestTopMetricsAggregation_SortBy(t *testing.T) {
	agg := NewTopMetricsAggregation().
		SortBy(SortByDoc{}).
		Field("a")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"top_metrics":{"metrics":[{"field":"a"}],"sort":"_doc"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestTopMetricsAggregation_SortWithInfo(t *testing.T) {
	agg := NewTopMetricsAggregation().
		SortWithInfo(SortInfo{Field: "f2", Ascending: true, UnmappedType: "int"}).
		Field("b")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"top_metrics":{"metrics":[{"field":"b"}],"sort":{"f2":{"order":"asc","unmapped_type":"int"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestTopMetricsAggregation_FailNoSorter(t *testing.T) {
	agg := NewTopMetricsAggregation().
		Field("a").
		Field("b")
	_, err := agg.Source()
	if err == nil || err.Error() != "sorter is required for the top metrics aggregation" {
		t.Fatal(err)
	}
}

func TestTopMetricsAggregation_FailNoFields(t *testing.T) {
	agg := NewTopMetricsAggregation().
		Sort("f1", false)
	_, err := agg.Source()
	if err == nil || err.Error() != "field list is required for the top metrics aggregation" {
		t.Fatal(err)
	}
}
