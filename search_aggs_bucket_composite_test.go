// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestCompositeAggregation(t *testing.T) {
	agg := NewCompositeAggregation().
		AddSource(NewCompositeAggregationSourceTerms("my_terms", "a_term").Missing("N/A")).
		AddSource(NewCompositeAggregationSourceHistogram("my_histogram", "price", 5)).
		AddSource(NewCompositeAggregationSourceDateHistogram("my_date_histogram", "purchase_date", "1d")).
		Size(10).
		After(map[string]interface{}{
			"my_terms":          "1",
			"my_histogram":      2,
			"my_date_histogram": "3",
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
	expected := `{"composite":{"after":{"my_date_histogram":"3","my_histogram":2,"my_terms":"1"},"size":10,"sources":[{"my_terms":{"terms":{"field":"a_term","missing":"N/A"}}},{"my_histogram":{"histogram":{"field":"price","interval":5}}},{"my_date_histogram":{"date_histogram":{"field":"purchase_date","interval":"1d"}}}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
