// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	_ "net/http"
	"testing"
)

func TestHighlightWithTermQuery(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun to do."}

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

	// Specify highlighter
	hl := NewHighlight()
	hl = hl.Fields(NewHighlighterField("message"))
	hl = hl.PreTags("<em>").PostTags("</em>")

	// Match all should return all documents
	query := NewPrefixQuery("message", "golang")
	searchResult, err := client.Search().
		Index(testIndexName).
		Highlight(hl).
		Query(&query).
		//Debug(true).Pretty(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Fatalf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.Hits.TotalHits != 1 {
		t.Fatalf("expected SearchResult.Hits.TotalHits = %d; got %d", 1, searchResult.Hits.TotalHits)
	}
	if len(searchResult.Hits.Hits) != 1 {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", 1, len(searchResult.Hits.Hits))
	}

	hit := searchResult.Hits.Hits[0]
	var tw tweet
	if err := json.Unmarshal(*hit.Source, &tw); err != nil {
		t.Fatal(err)
	}
	if hit.Highlight == nil || len(hit.Highlight) == 0 {
		t.Fatal("expected hit to have a highlight; got nil")
	}
	if hl, found := hit.Highlight["message"]; found {
		if len(hl) != 1 {
			t.Fatalf("expected to have one highlight for field \"message\"; got %d", len(hl))
		}
		expected := "Welcome to <em>Golang</em> and ElasticSearch."
		if hl[0] != expected {
			t.Errorf("expected to have highlight \"%s\"; got \"%s\"", expected, hl[0])
		}
	} else {
		t.Fatal("expected to have a highlight on field \"message\"; got none")
	}
}
