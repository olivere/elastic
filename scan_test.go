// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"io"
	_ "net/http"
	"testing"
)

func TestScan(t *testing.T) {
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

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if want, have := int64(3), cursor.Results.Hits.TotalHits; want != have {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", want, have)
	}
	if want, have := 0, len(cursor.Results.Hits.Hits); want != have {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
	}

	pages := 0
	docs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages++

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs++
		}
	}

	if pages != 4 {
		t.Fatalf("expected to retrieve %d pages; got %d", 4, pages)
	}
	if docs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 3, docs)
	}
}

func TestScanWithSort(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

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

	// We sort on a numerical field, because sorting on the 'message' string field would
	// raise the whole question of tokenizing and analyzing.
	cursor, err := client.Scan(testIndexName).Sort("retweets", true).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatal("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatal("expected results.Hits != nil; got nil")
	}
	if want, have := int64(3), cursor.Results.Hits.TotalHits; want != have {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", want, have)
	}
	if want, have := 1, len(cursor.Results.Hits.Hits); want != have {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
	}
	if want, have := "3", cursor.Results.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected hitID = %v; got %v", want, have)
	}

	docs := 1 // The cursor already gave us a result
	pages := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs += 1
		}
	}

	if pages != 3 {
		t.Fatalf("expected to retrieve %d pages; got %d", 3, pages)
	}
	if docs != 3 {
		t.Fatalf("expected to retrieve %d hits; got %d", 3, docs)
	}
}

func TestScanWithSortByDoc(t *testing.T) {
	// Sorting by doc is introduced in Elasticsearch 2.1,
	// and replaces the deprecated search_type=scan.
	// See https://www.elastic.co/guide/en/elasticsearch/reference/2.x/breaking_21_search_changes.html#_literal_search_type_scan_literal_deprecated
	client := setupTestClientAndCreateIndex(t)

	esversion, err := client.ElasticsearchVersion(DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	if esversion < "2.1" {
		t.Skipf(`Elasticsearch %s does not have {"sort":["_doc"]}`, esversion)
		return
	}

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	comment1 := comment{User: "nico", Comment: "You bet."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("comment").Id("1").Parent("1").BodyJson(&comment1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).Sort("_doc", true).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	docs := 0
	pages := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for range searchResult.Hits.Hits {
			docs += 1
		}
	}

	if pages != 3 {
		t.Fatalf("expected to retrieve %d pages; got %d", 3, pages)
	}
	if docs != 2 {
		t.Fatalf("expected to retrieve %d hits; got %d", 2, docs)
	}
}

func TestScanWithSearchSource(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

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

	src := NewSearchSource().
		Query(NewTermQuery("user", "olivere")).
		FetchSourceContext(NewFetchSourceContext(true).Include("retweets"))
	cursor, err := client.Scan(testIndexName).SearchSource(src).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if want, have := int64(2), cursor.Results.Hits.TotalHits; want != have {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", want, have)
	}

	docs := 0
	pages := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			if _, found := item["message"]; found {
				t.Fatalf("expected to not see field %q; got: %#v", "message", item)
			}
			docs += 1
		}
	}

	if pages != 3 {
		t.Fatalf("expected to retrieve %d pages; got %d", 3, pages)
	}
	if docs != 2 {
		t.Fatalf("expected to retrieve %d hits; got %d", 2, docs)
	}
}

func TestScanWithBody(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

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

	// Test with simple strings and a map
	var tests = []struct {
		Body              interface{}
		ExpectedTotalHits int64
		ExpectedDocs      int
		ExpectedPages     int
	}{
		{
			Body:              `{"query":{"match_all":{}}}`,
			ExpectedTotalHits: 3,
			ExpectedDocs:      3,
			ExpectedPages:     3,
		},
		/*
			{
				Body:              `{"query":{"term":{"user":"olivere"}},"sort":["_doc"]}`,
				ExpectedTotalHits: 2,
				ExpectedDocs:      2,
				ExpectedPages:     2,
			},
			{
				Body:              `{"query":{"term":{"user":"olivere"}},"sort":[{"retweets":"desc"}]}`,
				ExpectedTotalHits: 2,
				ExpectedDocs:      2,
				ExpectedPages:     2,
			},
			{
				Body: map[string]interface{}{
					"query": map[string]interface{}{
						"term": map[string]interface{}{
							"user": "olivere",
						},
					},
					"sort": []interface{}{"_doc"},
				},
				ExpectedTotalHits: 2,
				ExpectedDocs:      2,
				ExpectedPages:     2,
			},
		*/
	}

	for i, tt := range tests {
		cursor, err := client.Scan(testIndexName).Body(tt.Body).Size(1).Do()
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		if cursor.Results == nil {
			t.Fatalf("#%d: expected search results, got nil", i)
		}
		if want, have := tt.ExpectedTotalHits, cursor.Results.Hits.TotalHits; want != have {
			t.Fatalf("#%d: expected results.Hits.TotalHits = %d; got %d", i, want, have)
		}
		docs := len(cursor.Results.Hits.Hits)
		for {
			_, err = cursor.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("#%d: %v", i, err)
			}
			if cursor.Results == nil {
				t.Fatalf("#%d: expected search results, got nil", i)
			}
			docs += len(cursor.Results.Hits.Hits)
		}
		if want, have := tt.ExpectedDocs, docs; want != have {
			t.Fatalf("#%d: expected to retrieve %d documents; got %d", i, want, have)
		}
	}
}

func TestScanWithQuery(t *testing.T) {
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

	// Return tweets from olivere only
	termQuery := NewTermQuery("user", "olivere")
	cursor, err := client.Scan(testIndexName).
		Size(1).
		Query(termQuery).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatal("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatal("expected results.Hits != nil; got nil")
	}
	if want, have := int64(2), cursor.Results.Hits.TotalHits; want != have {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", want, have)
	}
	if want, have := 0, len(cursor.Results.Hits.Hits); want != have {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
	}

	pages := 0
	docs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs += 1
		}
	}

	if pages != 3 {
		t.Fatalf("expected to retrieve at %d pages; got %d", 3, pages)
	}
	if docs != 2 {
		t.Fatalf("expected to retrieve %d hits; got %d", 2, docs)
	}
}

func TestScanAndScrollWithMissingIndex(t *testing.T) {
	client := setupTestClient(t) // does not create testIndexName

	cursor, err := client.Scan(testIndexName).Scroll("30s").Do()
	if err == nil {
		t.Fatalf("expected error != nil; got: %v", err)
	}
	if cursor != nil {
		t.Fatalf("expected cursor == nil; got: %v", cursor)
	}
}

func TestScanAndScrollWithEmptyIndex(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	if isTravis() {
		t.Skip("test on Travis failes regularly with " +
			"Error 503 (Service Unavailable): SearchPhaseExecutionException[Failed to execute phase [init_scan], all shards failed]")
	}

	_, err := client.Flush().Index(testIndexName).WaitIfOngoing(true).Do()
	if err != nil {
		t.Fatal(err)
	}

	cursor, err := client.Scan(testIndexName).Scroll("30s").Do()
	if err != nil {
		t.Fatal(err)
	}
	if cursor == nil {
		t.Fatalf("expected cursor; got: %v", cursor)
	}

	// First request returns no error, but no hits
	res, err := cursor.Next()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected results != nil; got: nil")
	}
	if res.ScrollId == "" {
		t.Fatalf("expected scrollId in results; got: %q", res.ScrollId)
	}
	if want, have := int64(0), res.TotalHits(); want != have {
		t.Fatalf("expected TotalHits() = %d; got %d", want, have)
	}
	if res.Hits == nil {
		t.Fatal("expected results.Hits != nil; got: nil")
	}
	if want, have := int64(0), res.Hits.TotalHits; want != have {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", want, have)
	}
	if res.Hits.Hits == nil {
		t.Fatalf("expected results.Hits.Hits != nil; got: %v", res.Hits.Hits)
	}
	if want, have := 0, len(res.Hits.Hits); want != have {
		t.Fatalf("expected len(results.Hits.Hits) == %d; got: %d", want, have)
	}

	// Subsequent requests return EOS
	res, err = cursor.Next()
	if err != EOS {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected results == %v; got: %v", nil, res)
	}

	res, err = cursor.Next()
	if err != EOS {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected results == %v; got: %v", nil, res)
	}
}

func TestScanIssue119(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	comment1 := comment{User: "nico", Comment: "You bet."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}

	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("comment").Id("1").Parent("1").BodyJson(&comment1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).Fields("_source", "_parent").Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		for _, hit := range searchResult.Hits.Hits {
			if hit.Type == "tweet" {
				if _, ok := hit.Fields["_parent"].(string); ok {
					t.Errorf("Type `tweet` cannot have any parent...")

					toPrint, _ := json.MarshalIndent(hit, "", "    ")
					t.Fatal(string(toPrint))
				}
			}

			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
