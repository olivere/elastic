// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"io"
	_ "net/http"
	"testing"

	"golang.org/x/net/context"
)

func TestScroll(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
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

	// Should return all documents. Just don't call Do yet!
	svc := client.Scroll(testIndexName).Size(1)

	pages := 0
	numDocs := 0

	for {
		res, err := svc.Do()
		if err == EOS { // or err == io.EOF
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Errorf("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Errorf("expected results.Hits != nil; got nil")
		}
		if res.Hits.TotalHits != 3 {
			t.Errorf("expected results.Hits.TotalHits = %d; got %d", 3, res.Hits.TotalHits)
		}
		if len(res.Hits.Hits) != 1 {
			t.Errorf("expected len(results.Hits.Hits) = %d; got %d", 0, len(res.Hits.Hits))
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs++
		}

		if len(res.ScrollId) == 0 {
			t.Errorf("expeced scrollId in results; got %q", res.ScrollId)
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 3, numDocs)
	}

	err = svc.Clear(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Do()
	if err == nil {
		t.Fatal(err)
	}
}

func TestScrollWithQueryAndSort(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
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

	// Create a scroll service that returns tweets from user olivere
	// and returns them sorted by "message", in reverse order.
	//
	// Just don't call Do yet!
	svc := client.Scroll(testIndexName).
		Query(NewTermQuery("user", "olivere")).
		Sort("message", false).
		Size(1)

	numDocs := 0
	pages := 0
	for {
		res, err := svc.Do()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Errorf("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Errorf("expected results.Hits != nil; got nil")
		}
		if res.Hits.TotalHits != 2 {
			t.Errorf("expected results.Hits.TotalHits = %d; got %d", 2, res.Hits.TotalHits)
		}
		if len(res.Hits.Hits) != 1 {
			t.Errorf("expected len(results.Hits.Hits) = %d; got %d", 0, len(res.Hits.Hits))
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs++
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 2 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}
