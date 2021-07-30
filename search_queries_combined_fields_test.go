// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestCombinedFieldsQuery(t *testing.T) {
	q := NewCombinedFieldsQuery("query text", "f1", "f2").
		Field("f3").
		FieldWithBoost("f4", 2.0).
		AutoGenerateSynonymsPhraseQuery(false).
		Operator("AND").
		MinimumShouldMatch("3").
		ZeroTermsQuery("all")

	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"combined_fields":{"auto_generate_synonyms_phrase_query":false,"fields":["f1","f2","f3","f4^2.000000"],"minimum_should_match":"3","operator":"AND","query":"query text","zero_terms_query":"all"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
