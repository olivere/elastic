// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestReverseNestedAggregation(t *testing.T) {
	agg := NewReverseNestedAggregation()
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"reverse_nested":{}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestReverseNestedAggregationWithPath(t *testing.T) {
	agg := NewReverseNestedAggregation().Path("comments")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"reverse_nested":{"path":"comments"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestReverseNestedAggregationWithSubAggregation(t *testing.T) {
	avgPriceAgg := NewAvgAggregation().Field("price")
	agg := NewReverseNestedAggregation().
		Path("a_path").
		SubAggregation("avg_price", avgPriceAgg)
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"aggregations":{"avg_price":{"avg":{"field":"price"}}},"reverse_nested":{"path":"a_path"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestReverseNestedAggregationWithMeta(t *testing.T) {
	agg := NewReverseNestedAggregation().
		Path("a_path").
		Meta(map[string]interface{}{"name": "Oliver"})
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"meta":{"name":"Oliver"},"reverse_nested":{"path":"a_path"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
