// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRareTermsAggregation(t *testing.T) {
	agg := NewRareTermsAggregation().Field("genre")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"rare_terms":{"field":"genre"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRareTermsAggregationWithArgs(t *testing.T) {
	agg := NewRareTermsAggregation().
		Field("genre").
		MaxDocCount(2).
		Precision(0.1).
		Missing("n/a")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"rare_terms":{"field":"genre","max_doc_count":2,"missing":"n/a","precision":0.1}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRareTermsAggregationWithIncludeExclude(t *testing.T) {
	agg := NewRareTermsAggregation().Field("genre").Include("swi*").Exclude("electro*")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"rare_terms":{"exclude":"electro*","field":"genre","include":"swi*"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRareTermsAggregationWithIncludeExcludeValues(t *testing.T) {
	agg := NewRareTermsAggregation().Field("genre").IncludeValues("swing", "rock").ExcludeValues("jazz")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"rare_terms":{"exclude":["jazz"],"field":"genre","include":["swing","rock"]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRareTermsAggregationSubAggregation(t *testing.T) {
	genres := NewRareTermsAggregation().Field("genre")
	agg := NewTermsAggregation().Field("force")
	agg = agg.SubAggregation("genres", genres)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"aggregations":{"genres":{"rare_terms":{"field":"genre"}}},"terms":{"field":"force"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestRareTermsAggregationWithMetaData(t *testing.T) {
	agg := NewRareTermsAggregation().Field("genre")
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
	expected := `{"meta":{"name":"Oliver"},"rare_terms":{"field":"genre"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
