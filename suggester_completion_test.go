// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestCompletionSuggesterSource(t *testing.T) {
	s := NewCompletionSuggester("song-suggest").
		Text("n").
		Field("suggest")
	src, err := s.Source(true)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"song-suggest":{"text":"n","completion":{"field":"suggest"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCompletionSuggesterSourceWithMultipleContexts(t *testing.T) {
	s := NewCompletionSuggester("song-suggest").
		Text("n").
		Field("suggest").
		ContextQueries(
			NewSuggesterCategoryQuery("artist", "Sting"),
			NewSuggesterCategoryQuery("label", "BMG"),
		)
	src, err := s.Source(true)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expectedOption1 := `{"song-suggest":{"text":"n","completion":{"contexts":{"artist":[{"context":"Sting"}],"label":[{"context":"BMG"}]},"field":"suggest"}}}`
	expectedOption2 := `{"song-suggest":{"text":"n","completion":{"contexts":{"label":[{"context":"BMG"}],"artist":[{"context":"Sting"}]},"field":"suggest"}}}`
	if got != expectedOption1 && got != expectedOption2 {
		t.Errorf("expected either\n%s\nor\n%s\n,got:\n%s", expectedOption1, expectedOption2, got)
	}
}
