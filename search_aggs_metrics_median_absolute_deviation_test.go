// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestMedianAbsoluteDeviationAggregation(t *testing.T) {
	agg := NewMedianAbsoluteDeviationAggregation().Field("rating")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"median_absolute_deviation":{"field":"rating"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMedianAbsoluteDeviationAggregationWithOptions(t *testing.T) {
	agg := NewMedianAbsoluteDeviationAggregation().
		Field("rating").
		Compression(100).
		Missing(5)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"median_absolute_deviation":{"compression":100,"field":"rating","missing":5}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMedianAbsoluteDeviationAggregationWithScript(t *testing.T) {
	agg := NewMedianAbsoluteDeviationAggregation().
		Script(
			NewScript(`doc['rating'].value * params.scaleFactor`).
				Lang("painless").
				Param("scaleFactor", 2.0),
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
	expected := `{"median_absolute_deviation":{"script":{"lang":"painless","params":{"scaleFactor":2},"source":"doc['rating'].value * params.scaleFactor"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestMedianAbsoluteDeviationAggregationWithMetaData(t *testing.T) {
	agg := NewMedianAbsoluteDeviationAggregation().Field("rating").Meta(map[string]interface{}{"name": "Oliver"})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"median_absolute_deviation":{"field":"rating"},"meta":{"name":"Oliver"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
