// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestMatchBoolPrefixQuery(t *testing.T) {
	q := NewMatchBoolPrefixQuery("query_name", "this is a test").
		Analyzer("custom_analyzer").
		MinimumShouldMatch("75%").
		Operator("AND").
		Fuzziness("AUTO").
		PrefixLength(1).
		MaxExpansions(5).
		FuzzyTranspositions(false).
		FuzzyRewrite("constant_score").
		Boost(0.3)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"match_bool_prefix":{"query_name":{"analyzer":"custom_analyzer","boost":0.3,"fuzziness":"AUTO","fuzzy_rewrite":"constant_score","fuzzy_transpositions":false,"max_expansions":5,"minimum_should_match":"75%","operator":"AND","prefix_length":1,"query":"this is a test"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
