// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	_ "net/http"
	"testing"
)

func TestSuggestService(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Test _suggest endpoint
	termSuggesterName := "my-term-suggester"
	termSuggester := NewTermSuggester(termSuggesterName).Text("Goolang").Field("message")
	phraseSuggesterName := "my-phrase-suggester"
	phraseSuggester := NewPhraseSuggester(phraseSuggesterName).Text("Goolang").Field("message")

	result, err := client.Suggest().
		Index(testIndexName).
		Suggester(termSuggester).
		Suggester(phraseSuggester).
		// Debug(true).Pretty(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Errorf("expected result != nil; got nil")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 suggester results; got %d", len(result))
	}

	termSuggestions, found := result[termSuggesterName]
	if !found {
		t.Errorf("expected to find Suggest[%s]; got false", termSuggesterName)
	}
	if termSuggestions == nil {
		t.Errorf("expected Suggest[%s] != nil; got nil", termSuggesterName)
	}
	if len(termSuggestions) != 1 {
		t.Errorf("expected 1 suggestion; got %d", len(termSuggestions))
	}

	phraseSuggestions, found := result[phraseSuggesterName]
	if !found {
		t.Errorf("expected to find Suggest[%s]; got false", phraseSuggesterName)
	}
	if phraseSuggestions == nil {
		t.Errorf("expected Suggest[%s] != nil; got nil", phraseSuggesterName)
	}
	if len(phraseSuggestions) != 1 {
		t.Errorf("expected 1 suggestion; got %d", len(phraseSuggestions))
	}
}
