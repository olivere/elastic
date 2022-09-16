// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestVariableWidthHistogramAggregation(t *testing.T) {
	agg := NewVariableWidthHistogramAggregation().Field("price").Buckets(50)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"variable_width_histogram":{"buckets":50,"field":"price"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestVariableWidthHistogramAggregationWithMetaData(t *testing.T) {
	agg := NewVariableWidthHistogramAggregation().Field("price").InitialBuffer(100).ShardSize(500).Meta(map[string]interface{}{"name": "Oliver"})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"meta":{"name":"Oliver"},"variable_width_histogram":{"field":"price","initial_buffer":100,"shard_size":500}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
