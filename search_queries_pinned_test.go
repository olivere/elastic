// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestPinnedQueryTest(t *testing.T) {
	tests := []struct {
		Query    Query
		Expected string
	}{
		// #0
		{
			Query:    NewPinnedQuery(),
			Expected: `{"pinned":{}}`,
		},
		// #1
		{
			Query:    NewPinnedQuery().Ids("1", "2", "3"),
			Expected: `{"pinned":{"ids":["1","2","3"]}}`,
		},
		// #2
		{
			Query:    NewPinnedQuery().Organic(NewMatchAllQuery()),
			Expected: `{"pinned":{"organic":{"match_all":{}}}}`,
		},
		// #3
		{
			Query:    NewPinnedQuery().Ids("1", "2", "3").Organic(NewMatchAllQuery()),
			Expected: `{"pinned":{"ids":["1","2","3"],"organic":{"match_all":{}}}}`,
		},
	}

	for i, tt := range tests {
		src, err := tt.Query.Source()
		if err != nil {
			t.Fatalf("#%d: encoding Source failed: %v", i, err)
		}
		data, err := json.Marshal(src)
		if err != nil {
			t.Fatalf("#%d: marshaling to JSON failed: %v", i, err)
		}
		if want, got := tt.Expected, string(data); want != got {
			t.Fatalf("#%d: expected\n%s\ngot:\n%s", i, want, got)
		}
	}
}
