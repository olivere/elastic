// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestIntervalQuery(t *testing.T) {
	q := NewIntervalQuery(
		"my_text",
		NewIntervalQueryRuleAllOf(
			NewIntervalQueryRuleMatch("my favorite food").
				MaxGaps(0).
				Ordered(true).
				Filter(
					NewIntervalQueryFilter().
						NotContaining(
							NewIntervalQueryRuleMatch("salty"),
						),
				),
			NewIntervalQueryRuleAnyOf(
				NewIntervalQueryRuleMatch("hot water"),
				NewIntervalQueryRuleMatch("cold porridge"),
			),
		).Ordered(true),
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
	expected := `{"intervals":{"my_text":{"all_of":{"intervals":[{"match":{"filter":{"not_containing":{"match":{"query":"salty"}}},"max_gaps":0,"ordered":true,"query":"my favorite food"}},{"any_of":{"intervals":[{"match":{"query":"hot water"}},{"match":{"query":"cold porridge"}}]}}],"ordered":true}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
