// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"io"
	_ "net/http"
	"testing"
)

func TestScroll(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Should return all documents. Just don't call Do yet!
	svc := client.Scroll(testIndexName).Size(1)

	pages := 0
	docs := 0

	for {
		res, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Fatal("expected results.Hits != nil; got nil")
		}
		if want, have := int64(3), res.TotalHits(); want != have {
			t.Fatalf("expected results.TotalHits() = %d; got %d", want, have)
		}
		if want, have := 1, len(res.Hits.Hits); want != have {
			t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs++
		}

		if len(res.ScrollId) == 0 {
			t.Fatalf("expected scrollId in results; got %q", res.ScrollId)
		}
	}

	if want, have := 3, pages; want != have {
		t.Fatalf("expected to retrieve %d pages; got %d", want, have)
	}
	if want, have := 3, docs; want != have {
		t.Fatalf("expected to retrieve %d hits; got %d", want, have)
	}

	err = svc.Clear(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Do(context.TODO())
	if err == nil {
		t.Fatal("expected to fail")
	}
}

func TestScrollWithQueryAndSort(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
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
		Size(1).
		TrackTotalHits(true)

	docs := 0
	pages := 0
	for {
		res, err := svc.Do(context.TODO())
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
			t.Fatal("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Fatal("expected results.Hits != nil; got nil")
		}
		if want, have := int64(2), res.TotalHits(); want != have {
			t.Fatalf("expected results.TotalHits() = %d; got %d", want, have)
		}
		if want, have := 1, len(res.Hits.Hits); want != have {
			t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs++
		}
	}

	if want, have := 2, pages; want != have {
		t.Fatalf("expected to retrieve %d pages; got %d", want, have)
	}
	if want, have := 2, docs; want != have {
		t.Fatalf("expected to retrieve %d hits; got %d", want, have)
	}
}

func TestScrollWithBody(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
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
	}

	for i, tt := range tests {
		// Should return all documents. Just don't call Do yet!
		svc := client.Scroll(testIndexName).Size(1).Body(tt.Body)

		pages := 0
		docs := 0

		for {
			res, err := svc.Do(context.TODO())
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if res == nil {
				t.Fatalf("#%d: expected results != nil; got nil", i)
			}
			if res.Hits == nil {
				t.Fatalf("#%d: expected results.Hits != nil; got nil", i)
			}
			if want, have := tt.ExpectedTotalHits, res.TotalHits(); want != have {
				t.Fatalf("#%d: expected results.TotalHits() = %d; got %d", i, want, have)
			}
			if want, have := 1, len(res.Hits.Hits); want != have {
				t.Fatalf("#%d: expected len(results.Hits.Hits) = %d; got %d", i, want, have)
			}

			pages++

			for _, hit := range res.Hits.Hits {
				if hit.Index != testIndexName {
					t.Fatalf("#%d: expected SearchResult.Hits.Hit.Index = %q; got %q", i, testIndexName, hit.Index)
				}
				item := make(map[string]interface{})
				err := json.Unmarshal(hit.Source, &item)
				if err != nil {
					t.Fatalf("#%d: %v", i, err)
				}
				docs++
			}

			if len(res.ScrollId) == 0 {
				t.Fatalf("#%d: expected scrollId in results; got %q", i, res.ScrollId)
			}
		}

		if want, have := tt.ExpectedPages, pages; want != have {
			t.Fatalf("#%d: expected to retrieve %d pages; got %d", i, want, have)
		}
		if want, have := tt.ExpectedDocs, docs; want != have {
			t.Fatalf("#%d: expected to retrieve %d hits; got %d", i, want, have)
		}

		err = svc.Clear(context.TODO())
		if err != nil {
			t.Fatalf("#%d: failed to clear scroll context: %v", i, err)
		}

		_, err = svc.Do(context.TODO())
		if err == nil {
			t.Fatalf("#%d: expected to fail", i)
		}
	}
}

func TestScrollWithSlice(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// Should return all documents. Just don't call Do yet!
	sliceQuery := NewSliceQuery().Id(0).Max(2)
	svc := client.Scroll(testIndexName).Slice(sliceQuery).Size(1)

	pages := 0
	docs := 0

	for {
		res, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Fatal("expected results.Hits != nil; got nil")
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs++
		}

		if len(res.ScrollId) == 0 {
			t.Fatalf("expected scrollId in results; got %q", res.ScrollId)
		}
	}

	if pages == 0 {
		t.Fatal("expected to retrieve some pages")
	}
	if docs == 0 {
		t.Fatal("expected to retrieve some hits")
	}

	if err := svc.Clear(context.TODO()); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.Do(context.TODO()); err == nil {
		t.Fatal("expected to fail")
	}
}

func TestScrollWithMaxResponseSize(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}
	tweet2 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Test response size error on first scroll request (first response is 418 bytes)
	svc := client.Scroll(testIndexName).Size(1).MaxResponseSize(400)
	_, err = svc.Do(context.TODO())
	if err != ErrResponseSize {
		t.Fatalf("expected response size error")
	}

	// Test response size error on second scroll request (first response is 418 bytes, second is 439 bytes)
	svc = client.Scroll(testIndexName).Size(1).MaxResponseSize(420)
	_, err = svc.Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Do(context.TODO())
	if err != ErrResponseSize {
		t.Fatalf("expected response size error")
	}
}

func TestScrollWithFilterPath(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Should return all documents. Just don't call Do yet!
	// Notice that we don't have to add "_scroll_id" to the FilterPath here:
	// It's been added automatically by the ScrollService.
	svc := client.Scroll(testIndexName).Size(1).
		FilterPath("hits.total", "hits.hits._index", "hits.hits._id", "hits.hits._source")

	pages := 0
	docs := 0

	for {
		res, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Fatal("expected results.Hits != nil; got nil")
		}
		if want, have := int64(3), res.TotalHits(); want != have {
			t.Fatalf("expected results.TotalHits() = %d; got %d", want, have)
		}
		if want, have := 1, len(res.Hits.Hits); want != have {
			t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			docs++
		}

		if len(res.ScrollId) == 0 {
			t.Fatalf("expected scrollId in results; got %q", res.ScrollId)
		}

		// Ensure we don't have "_scroll_id" in the query string more than once
		var scrollIdCount int
		for _, path := range svc.filterPath {
			if path == "_scroll_id" {
				scrollIdCount++
			}
		}
		if want, have := 1, scrollIdCount; want != have {
			t.Fatalf("expected _scroll_id to occur %d time(s), got %d", want, have)
		}
	}

	if want, have := 3, pages; want != have {
		t.Fatalf("expected to retrieve %d pages; got %d", want, have)
	}
	if want, have := 3, docs; want != have {
		t.Fatalf("expected to retrieve %d hits; got %d", want, have)
	}

	err = svc.Clear(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Do(context.TODO())
	if err == nil {
		t.Fatal("expected to fail")
	}
}

func TestScrollWithFilterPathKeepingContext(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Should return all documents. Just don't call Do yet!
	// Notice that we don't have to add "_scroll_id" to the FilterPath here:
	// It's been added automatically by the ScrollService.
	svc := client.Scroll(testIndexName).Size(1).
		FilterPath("hits.total", "hits.hits._index", "hits.hits._id")

	pages := 0
	docs := 0

	for {
		res, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected results != nil; got nil")
		}
		if res.Hits == nil {
			t.Fatal("expected results.Hits != nil; got nil")
		}
		if want, have := int64(3), res.TotalHits(); want != have {
			t.Fatalf("expected results.TotalHits() = %d; got %d", want, have)
		}
		if want, have := 1, len(res.Hits.Hits); want != have {
			t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", want, have)
		}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			if hit.Source != nil {
				t.Fatal("expected SearchResult.Hits.Hit.Source = nil")
			}
			docs++
		}

		if len(res.ScrollId) == 0 {
			t.Fatalf("expected scrollId in results; got %q", res.ScrollId)
		}
	}

	if want, have := 3, pages; want != have {
		t.Fatalf("expected to retrieve %d pages; got %d", want, have)
	}
	if want, have := 3, docs; want != have {
		t.Fatalf("expected to retrieve %d hits; got %d", want, have)
	}

	err = svc.Clear(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Do(context.TODO())
	if err == nil {
		t.Fatal("expected to fail")
	}
}
