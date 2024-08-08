// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestBoostingQuery(t *testing.T) {
	q := NewBoostingQuery()
	q = q.Positive(NewTermQuery("tag", "wow"))
	q = q.Negative(NewRangeQuery("age").From(10).To(20))
	q = q.NegativeBoost(0.2)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"boosting":{"negative":{"range":{"age":{"gte":10,"lte":20}}},"negative_boost":0.2,"positive":{"term":{"tag":"wow"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
