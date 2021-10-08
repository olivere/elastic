// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestWeightedAvgAggregation(t *testing.T) {
	agg := NewWeightedAvgAggregation().
		Value(&MultiValuesSourceFieldConfig{
			FieldName: "grade",
		}).
		Weight(&MultiValuesSourceFieldConfig{
			FieldName: "weight",
			Missing:   3,
		})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"weighted_avg":{"value":{"field":"grade"},"weight":{"field":"weight","missing":3}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
