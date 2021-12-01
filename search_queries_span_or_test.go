// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestSpanOrQuery(t *testing.T) {
	q := NewSpanOrQuery(
		NewSpanTermQuery("field", "value1"),
		NewSpanTermQuery("field", "value2"),
		NewSpanTermQuery("field", "value3"),
	)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_or":{"clauses":[{"span_term":{"field":{"value":"value1"}}},{"span_term":{"field":{"value":"value2"}}},{"span_term":{"field":{"value":"value3"}}}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
