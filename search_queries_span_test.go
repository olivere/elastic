// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.
package elastic

import (
	"encoding/json"
	"testing"
)

func TestSpanNearQuery(t *testing.T) {

	q := NewSpanNearQuery(1, false, NewClause("field", "value"))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_near":{"clauses":[{"span_term":{"field":"value"}}],"in_order":false,"slop":1}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
