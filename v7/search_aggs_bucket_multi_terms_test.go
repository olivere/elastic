// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestMultiTermsAggregation(t *testing.T) {
	agg := NewMultiTermsAggregation().Terms("genre", "product")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_terms":{"terms":[{"field":"genre"},{"field":"product"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiTermsAggregationWithMultiTerms(t *testing.T) {
	agg := NewMultiTermsAggregation().MultiTerms(
		MultiTerm{Field: "genre", Missing: "n/a"},
		MultiTerm{Field: "product", Missing: "n/a"},
	)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_terms":{"terms":[{"field":"genre","missing":"n/a"},{"field":"product","missing":"n/a"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiTermsAggregationWithSubAggregation(t *testing.T) {
	subAgg := NewAvgAggregation().Field("height")
	agg := NewMultiTermsAggregation().Terms("genre", "product").Size(10).
		OrderByAggregation("avg_height", false)
	agg = agg.SubAggregation("avg_height", subAgg)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"aggregations":{"avg_height":{"avg":{"field":"height"}}},"multi_terms":{"order":[{"avg_height":"desc"}],"size":10,"terms":[{"field":"genre"},{"field":"product"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiTermsAggregationWithMultipleSubAggregation(t *testing.T) {
	subAgg1 := NewAvgAggregation().Field("height")
	subAgg2 := NewAvgAggregation().Field("width")
	agg := NewMultiTermsAggregation().Terms("genre", "product").Size(10).
		OrderByAggregation("avg_height", false)
	agg = agg.SubAggregation("avg_height", subAgg1)
	agg = agg.SubAggregation("avg_width", subAgg2)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"aggregations":{"avg_height":{"avg":{"field":"height"}},"avg_width":{"avg":{"field":"width"}}},"multi_terms":{"order":[{"avg_height":"desc"}],"size":10,"terms":[{"field":"genre"},{"field":"product"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiTermsAggregationWithMetaData(t *testing.T) {
	agg := NewMultiTermsAggregation().Terms("genre", "product").Size(10).OrderByKeyDesc()
	agg = agg.Meta(map[string]interface{}{"name": "Oliver"})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"meta":{"name":"Oliver"},"multi_terms":{"order":[{"_key":"desc"}],"size":10,"terms":[{"field":"genre"},{"field":"product"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMultiTermsAggregationWithMissing(t *testing.T) {
	agg := NewMultiTermsAggregation().MultiTerms(
		MultiTerm{Field: "genre"},
		MultiTerm{Field: "product", Missing: "n/a"},
	).Size(10)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_terms":{"size":10,"terms":[{"field":"genre"},{"field":"product","missing":"n/a"}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
