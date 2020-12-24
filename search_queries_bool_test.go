// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestBoolQuery(t *testing.T) {
	q := NewBoolQuery()
	q = q.Must(NewTermQuery("tag", "wow"))
	q = q.MustNot(NewRangeQuery("age").From(10).To(20))
	q = q.Filter(NewTermQuery("account", "1"))
	q = q.Should(NewTermQuery("tag", "sometag"), NewTermQuery("tag", "sometagtag"))
	q = q.Boost(10)
	q = q.QueryName("Test")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"bool":{"_name":"Test","boost":10,"filter":{"term":{"account":"1"}},"must":{"term":{"tag":"wow"}},"must_not":{"range":{"age":{"gte":10,"lte":20}}},"should":[{"term":{"tag":"sometag"}},{"term":{"tag":"sometagtag"}}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
