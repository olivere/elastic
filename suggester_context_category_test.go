// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestSuggesterCategoryMapping(t *testing.T) {
	q := NewSuggesterCategoryMapping("color").DefaultValues("red")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"color":{"default":"red","type":"category"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterCategoryMappingWithTwoDefaultValues(t *testing.T) {
	q := NewSuggesterCategoryMapping("color").DefaultValues("red", "orange")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"color":{"default":["red","orange"],"type":"category"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterCategoryMappingWithFieldName(t *testing.T) {
	q := NewSuggesterCategoryMapping("color").
		DefaultValues("red", "orange").
		FieldName("color_field")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"color":{"default":["red","orange"],"path":"color_field","type":"category"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterCategoryQuery(t *testing.T) {
	q := NewSuggesterCategoryQuery("color", "red")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"color":[{"context":"red"}]}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterCategoryQueryWithTwoValues(t *testing.T) {
	q := NewSuggesterCategoryQuery("color", "red", "yellow")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expectedOption1 := `{"color":[{"context":"red"},{"context":"yellow"}]}`
	expectedOption2 := `{"color":[{"context":"yellow"},{"context":"red"}]}` // order is irrelevant to the results, and we model the query with a map which has no order guarantees
	if got != expectedOption1 && got != expectedOption2 {
		t.Errorf("expected either\n%s\nor\n%s\n,got:\n%s", expectedOption1, expectedOption2, got)
	}
}

func TestSuggesterCategoryQueryWithBoost(t *testing.T) {
	q := NewSuggesterCategoryQuery("color", "red")
	q.ValueWithBoost("yellow", 4)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expectedOption1 := `{"color":[{"context":"red"},{"boost":4,"context":"yellow"}]}`
	expectedOption2 := `{"color":[{"boost":4,"context":"yellow"},{"context":"red"}]}`
	if got != expectedOption1 && got != expectedOption2 {
		t.Errorf("expected either\n%s\nor\n%s\n,got:\n%s", expectedOption1, expectedOption2, got)
	}
}

func TestSuggesterCategoryQueryWithoutBoost(t *testing.T) {
	q := NewSuggesterCategoryQuery("color", "red")
	q.Value("yellow")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expectedOption1 := `{"color":[{"context":"red"},{"context":"yellow"}]}`
	expectedOption2 := `{"color":[{"context":"yellow"},{"context":"red"}]}`
	if got != expectedOption1 && got != expectedOption2 {
		t.Errorf("expected either\n%s\nor\n%s\n,got:\n%s", expectedOption1, expectedOption2, got)
	}
}
