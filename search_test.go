// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestSearchMatchAll(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Size(100).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(3); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 3; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchWithCustomHTTPHeaders(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Match all should return all documents
	res, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Size(100).
		Pretty(true).
		Headers(http.Header{
			"X-ID":      []string{"A", "B"},
			"Custom-ID": []string{"olivere"},
		}).
		Header("X-ID", "12345").
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.TotalHits(), int64(3); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := res.Header.Get("Content-Type"), "application/json; charset=UTF-8"; got != want {
		t.Errorf("expected SearchResult.Header(%q) = %q; got %q", "Content-Type", want, got)
	}
}

func TestSearchMatchAllWithRequestCacheDisabled(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Match all should return all documents, with request cache disabled
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Size(100).
		Pretty(true).
		RequestCache(false).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(3); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 3; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}
}

func TestSearchTotalHits(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	count, err := client.Count(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatalf("expected more than %d documents", count)
	}

	// RestTotalHitsAsInt(false) (default)
	{
		res, err := client.Search().Index(testIndexName).Query(NewMatchAllQuery()).RestTotalHitsAsInt(false).Pretty(true).Do(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected SearchResult != nil; got nil")
		}
		if want, have := count, res.TotalHits(); want != have {
			t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, have)
		}
		if res.Hits == nil || res.Hits.TotalHits == nil {
			t.Fatal("expected SearchResult.Hits._ != nil; got nil")
		}
		if want, have := count, res.Hits.TotalHits.Value; want != have {
			t.Errorf("expected SearchResult.TotalHits.Value = %d; got %d", want, have)
		}
		if want, have := "eq", res.Hits.TotalHits.Relation; want != have {
			t.Errorf("expected SearchResult.TotalHits.Relation = %q; got %q", want, have)
		}
	}

	// RestTotalHitsAsInt(true)
	{
		res, err := client.Search().Index(testIndexName).Query(NewMatchAllQuery()).RestTotalHitsAsInt(true).Pretty(true).Do(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected SearchResult != nil; got nil")
		}
		if want, have := count, res.TotalHits(); want != have {
			t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, have)
		}
		if res.Hits == nil || res.Hits.TotalHits == nil {
			t.Fatal("expected SearchResult.Hits._ != nil; got nil")
		}
		if want, have := count, res.Hits.TotalHits.Value; want != have {
			t.Errorf("expected SearchResult.TotalHits.Value = %d; got %d", want, have)
		}
		if want, have := "eq", res.Hits.TotalHits.Relation; want != have {
			t.Errorf("expected SearchResult.TotalHits.Relation = %q; got %q", want, have)
		}
	}
}

func BenchmarkSearchMatchAll(b *testing.B) {
	client := setupTestClientAndCreateIndexAndAddDocs(b)

	for n := 0; n < b.N; n++ {
		// Match all should return all documents
		all := NewMatchAllQuery()
		searchResult, err := client.Search().Index(testIndexName).Query(all).Do(context.TODO())
		if err != nil {
			b.Fatal(err)
		}
		if searchResult.Hits == nil {
			b.Errorf("expected SearchResult.Hits != nil; got nil")
		}
		if searchResult.TotalHits() == 0 {
			b.Errorf("expected SearchResult.TotalHits() > %d; got %d", 0, searchResult.TotalHits())
		}
	}
}

func TestSearchResultTotalHits(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	count, err := client.Count(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	all := NewMatchAllQuery()
	searchResult, err := client.Search().Index(testIndexName).Query(all).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	got := searchResult.TotalHits()
	if got != count {
		t.Fatalf("expected %d hits; got: %d", count, got)
	}

	// No hits
	searchResult = &SearchResult{}
	got = searchResult.TotalHits()
	if got != 0 {
		t.Errorf("expected %d hits; got: %d", 0, got)
	}
}

func TestSearchResultWithProfiling(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	all := NewMatchAllQuery()
	searchResult, err := client.Search().Index(testIndexName).Query(all).Profile(true).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if searchResult.Profile == nil {
		t.Fatal("Profiled MatchAll query did not return profiling data with results")
	}
}

func TestSearchResultEach(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	all := NewMatchAllQuery()
	searchResult, err := client.Search().Index(testIndexName).Query(all).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Iterate over non-ptr type
	var aTweet tweet
	count := 0
	for _, item := range searchResult.Each(reflect.TypeOf(aTweet)) {
		count++
		_, ok := item.(tweet)
		if !ok {
			t.Fatalf("expected hit to be serialized as tweet; got: %v", reflect.ValueOf(item))
		}
	}
	if count == 0 {
		t.Errorf("expected to find some hits; got: %d", count)
	}

	// Iterate over ptr-type
	count = 0
	var aTweetPtr *tweet
	for _, item := range searchResult.Each(reflect.TypeOf(aTweetPtr)) {
		count++
		tw, ok := item.(*tweet)
		if !ok {
			t.Fatalf("expected hit to be serialized as tweet; got: %v", reflect.ValueOf(item))
		}
		if tw == nil {
			t.Fatal("expected hit to not be nil")
		}
	}
	if count == 0 {
		t.Errorf("expected to find some hits; got: %d", count)
	}

	// Iterate over ptr-type
	count = 0
	var aTweetPtrWithID *tweetWithID
	for _, item := range searchResult.Each(reflect.TypeOf(aTweetPtrWithID)) {
		count++
		tw, ok := item.(*tweetWithID)
		if !ok {
			t.Fatalf("expected hit to be serialized as tweet; got: %v", reflect.ValueOf(item))
		}
		if tw == nil {
			t.Fatal("expected hit to not be nil")
		}
		if tw.ElasticID == "" {
			t.Fatal("No ID setup in the structure")
		}
	}
	if count == 0 {
		t.Errorf("expected to find some hits; got: %d", count)
	}

	// Does not iterate when no hits are found
	searchResult = &SearchResult{Hits: nil}
	count = 0
	for _, item := range searchResult.Each(reflect.TypeOf(aTweet)) {
		count++
		_ = item
	}
	if count != 0 {
		t.Errorf("expected to not find any hits; got: %d", count)
	}
	searchResult = &SearchResult{Hits: &SearchHits{Hits: make([]*SearchHit, 0)}}
	count = 0
	for _, item := range searchResult.Each(reflect.TypeOf(aTweet)) {
		count++
		_ = item
	}
	if count != 0 {
		t.Errorf("expected to not find any hits; got: %d", count)
	}
}

func TestSearchResultEachNoSource(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocsNoSource(t)

	all := NewMatchAllQuery()
	searchResult, err := client.Search().Index(testNoSourceIndexName).Query(all).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Iterate over non-ptr type
	var aTweet tweet
	count := 0
	for _, item := range searchResult.Each(reflect.TypeOf(aTweet)) {
		count++
		tw, ok := item.(tweet)
		if !ok {
			t.Fatalf("expected hit to be serialized as tweet; got: %v", reflect.ValueOf(item))
		}

		if tw.User != "" {
			t.Fatalf("expected no _source hit to be empty tweet; got: %v", reflect.ValueOf(item))
		}
	}
	if count != 2 {
		t.Errorf("expected to find 2 hits; got: %d", count)
	}

	// Iterate over ptr-type
	count = 0
	var aTweetPtr *tweet
	for _, item := range searchResult.Each(reflect.TypeOf(aTweetPtr)) {
		count++
		tw, ok := item.(*tweet)
		if !ok {
			t.Fatalf("expected hit to be serialized as tweet; got: %v", reflect.ValueOf(item))
		}
		if tw != nil {
			t.Fatal("expected hit to be nil")
		}
	}
	if count != 2 {
		t.Errorf("expected to find 2 hits; got: %d", count)
	}
}

func TestSearchSorting(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	// Match all should return all documents
	all := NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(all).
		Sort("created", false).
		Timeout("1s").
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 3, len(searchResult.Hits.Hits))
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchSortingBySorters(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	// Match all should return all documents
	all := NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(all).
		SortBy(NewFieldSort("created").Desc(), NewScoreSort()).
		Timeout("1s").
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 3, len(searchResult.Hits.Hits))
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchSpecificFields(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Retweets: 1, Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Retweets: 2, Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Retweets: 3, Message: "Cycling is fun."}
	tweets := []tweet{
		tweet1,
		tweet2,
		tweet3,
	}

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

	// Match all should return all documents
	all := NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(all).
		StoredFields("message").
		DocvalueFields("retweets").
		Sort("_id", true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 3, len(searchResult.Hits.Hits))
	}

	// Manually inspect the fields
	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		if hit.Source != nil {
			t.Fatalf("expected SearchResult.Hits.Hit.Source to be nil; got: %v", hit.Source)
		}
		if hit.Fields == nil {
			t.Fatal("expected SearchResult.Hits.Hit.Fields to be != nil")
		}
		field, found := hit.Fields["message"]
		if !found {
			t.Errorf("expected SearchResult.Hits.Hit.Fields[%s] to be found", "message")
		}
		fields, ok := field.([]interface{})
		if !ok {
			t.Errorf("expected []interface{}; got: %v", reflect.TypeOf(fields))
		}
		if len(fields) != 1 {
			t.Errorf("expected a field with 1 entry; got: %d", len(fields))
		}
		message, ok := fields[0].(string)
		if !ok {
			t.Errorf("expected a string; got: %v", reflect.TypeOf(fields[0]))
		}
		if message == "" {
			t.Errorf("expected a message; got: %q", message)
		}
	}

	// With the new helper method for fields
	for i, hit := range searchResult.Hits.Hits {
		// Field: message
		items, ok := hit.Fields.Strings("message")
		if !ok {
			t.Fatalf("expected SearchResult.Hits.Hit.Fields[%s] to be found", "message")
		}
		if want, have := 1, len(items); want != have {
			t.Fatalf("expected a field with %d entries; got %d", want, have)
		}
		if want, have := tweets[i].Message, items[0]; want != have {
			t.Fatalf("expected message[%d]=%q; got %q", i, want, have)
		}

		// Field: retweets
		retweets, ok := hit.Fields.Float64s("retweets")
		if !ok {
			t.Fatalf("expected SearchResult.Hits.Hit.Fields[%s] to be found", "retweets")
		}
		if want, have := 1, len(retweets); want != have {
			t.Fatalf("expected a field with %d entries; got %d", want, have)
		}
		if want, have := tweets[i].Retweets, int(retweets[0]); want != have {
			t.Fatalf("expected retweets[%d]=%q; got %q", i, want, have)
		}

		// Field should not exist
		numbers, ok := hit.Fields.Float64s("score")
		if ok {
			t.Fatalf("expected SearchResult.Hits.Hit.Fields[%s] to NOT be found", "numbers")
		}
		if numbers != nil {
			t.Fatalf("expected no field %q; got %+v", "numbers", numbers)
		}
	}
}

func TestSearchExplain(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	// Match all should return all documents
	all := NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(all).
		Explain(true).
		Timeout("1s").
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 3, len(searchResult.Hits.Hits))
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		if hit.Explanation == nil {
			t.Fatal("expected search explanation")
		}
		if hit.Explanation.Value <= 0.0 {
			t.Errorf("expected explanation value to be > 0.0; got: %v", hit.Explanation.Value)
		}
		if hit.Explanation.Description == "" {
			t.Errorf("expected explanation description != %q; got: %q", "", hit.Explanation.Description)
		}
	}
}

func TestSearchSource(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	// Set up the request JSON manually to pass to the search service via Source()
	source := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	searchResult, err := client.Search().
		Index(testIndexName).
		Source(source). // sets the JSON request
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
}

func TestSearchSourceWithString(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	searchResult, err := client.Search().
		Index(testIndexName).
		Source(`{"query":{"match_all":{}}}`). // sets the JSON request
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
}

func TestSearchRawString(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	query := RawStringQuery(`{"match_all":{}}`)
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(query).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
}

func TestSearchSearchSource(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	// Set up the search source manually and pass it to the search service via SearchSource()
	ss := NewSearchSource().
		Query(NewMatchAllQuery()).
		IndexBoost(testIndexName, 1.0).
		IndexBoosts(IndexBoost{Index: testIndexName2, Boost: 2.0}).
		From(0).Size(2)

	// One can use ss.Source() to get to the raw interface{} that will be used
	// as the search request JSON by the SearchService.

	searchResult, err := client.Search().
		Index(testIndexName).
		SearchSource(ss). // sets the SearchSource
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 2 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 2, len(searchResult.Hits.Hits))
	}
}

func TestSearchInnerHitsOnHasChild(t *testing.T) {
	// client := setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	ctx := context.Background()

	// Create join index
	createIndex, err := client.CreateIndex(testJoinIndex).Body(testJoinMapping).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Add documents
	// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/parent-join.html for example code.
	doc1 := joinDoc{
		Message:   "This is a question",
		JoinField: &joinField{Name: "question"},
	}
	_, err = client.Index().Index(testJoinIndex).Id("1").BodyJson(&doc1).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc2 := joinDoc{
		Message:   "This is another question",
		JoinField: "question",
	}
	_, err = client.Index().Index(testJoinIndex).Id("2").BodyJson(&doc2).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc3 := joinDoc{
		Message: "This is an answer",
		JoinField: &joinField{
			Name:   "answer",
			Parent: "1",
		},
	}
	_, err = client.Index().Index(testJoinIndex).Id("3").BodyJson(&doc3).Routing("1").Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc4 := joinDoc{
		Message: "This is another answer",
		JoinField: &joinField{
			Name:   "answer",
			Parent: "1",
		},
	}
	_, err = client.Index().Index(testJoinIndex).Id("4").BodyJson(&doc4).Routing("1").Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testJoinIndex).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Search for all documents that have an answer, and return those answers as inner hits
	bq := NewBoolQuery()
	bq = bq.Must(NewMatchAllQuery())
	bq = bq.Filter(NewHasChildQuery("answer", NewMatchAllQuery()).
		InnerHit(NewInnerHit().Name("answers")))

	searchResult, err := client.Search().
		Index(testJoinIndex).
		Query(bq).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 1 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 2, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 1 {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", 2, len(searchResult.Hits.Hits))
	}

	hit := searchResult.Hits.Hits[0]
	if want, have := "1", hit.Id; want != have {
		t.Fatalf("expected tweet %q; got: %q", want, have)
	}
	if hit.InnerHits == nil {
		t.Fatalf("expected inner hits; got: %v", hit.InnerHits)
	}
	if want, have := 1, len(hit.InnerHits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	innerHits, found := hit.InnerHits["answers"]
	if !found {
		t.Fatalf("expected inner hits for name %q", "answers")
	}
	if innerHits == nil || innerHits.Hits == nil {
		t.Fatal("expected inner hits != nil")
	}
	if want, have := 2, len(innerHits.Hits.Hits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	if want, have := "3", innerHits.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}
	if want, have := "4", innerHits.Hits.Hits[1].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}
}

func TestSearchInnerHitsOnHasParent(t *testing.T) {
	// client := setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	ctx := context.Background()

	// Create join index
	createIndex, err := client.CreateIndex(testJoinIndex).Body(testJoinMapping).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Add documents
	// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/parent-join.html for example code.
	doc1 := joinDoc{
		Message:   "This is a question",
		JoinField: &joinField{Name: "question"},
	}
	_, err = client.Index().Index(testJoinIndex).Id("1").BodyJson(&doc1).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc2 := joinDoc{
		Message:   "This is another question",
		JoinField: "question",
	}
	_, err = client.Index().Index(testJoinIndex).Id("2").BodyJson(&doc2).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc3 := joinDoc{
		Message: "This is an answer",
		JoinField: &joinField{
			Name:   "answer",
			Parent: "1",
		},
	}
	_, err = client.Index().Index(testJoinIndex).Id("3").BodyJson(&doc3).Routing("1").Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	doc4 := joinDoc{
		Message: "This is another answer",
		JoinField: &joinField{
			Name:   "answer",
			Parent: "1",
		},
	}
	_, err = client.Index().Index(testJoinIndex).Id("4").BodyJson(&doc4).Routing("1").Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Refresh().Index(testJoinIndex).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Search for all documents that have an answer, and return those answers as inner hits
	bq := NewBoolQuery()
	bq = bq.Must(NewMatchAllQuery())
	bq = bq.Filter(NewHasParentQuery("question", NewMatchAllQuery()).
		InnerHit(NewInnerHit().Name("answers")))

	searchResult, err := client.Search().
		Index(testJoinIndex).
		Query(bq).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if want, have := int64(2), searchResult.TotalHits(); want != have {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, have)
	}
	if want, have := 2, len(searchResult.Hits.Hits); want != have {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, have)
	}

	hit := searchResult.Hits.Hits[0]
	if want, have := "3", hit.Id; want != have {
		t.Fatalf("expected tweet %q; got: %q", want, have)
	}
	if hit.InnerHits == nil {
		t.Fatalf("expected inner hits; got: %v", hit.InnerHits)
	}
	if want, have := 1, len(hit.InnerHits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	innerHits, found := hit.InnerHits["answers"]
	if !found {
		t.Fatalf("expected inner hits for name %q", "tweets")
	}
	if innerHits == nil || innerHits.Hits == nil {
		t.Fatal("expected inner hits != nil")
	}
	if want, have := 1, len(innerHits.Hits.Hits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	if want, have := "1", innerHits.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}

	hit = searchResult.Hits.Hits[1]
	if want, have := "4", hit.Id; want != have {
		t.Fatalf("expected tweet %q; got: %q", want, have)
	}
	if hit.InnerHits == nil {
		t.Fatalf("expected inner hits; got: %v", hit.InnerHits)
	}
	if want, have := 1, len(hit.InnerHits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	innerHits, found = hit.InnerHits["answers"]
	if !found {
		t.Fatalf("expected inner hits for name %q", "tweets")
	}
	if innerHits == nil || innerHits.Hits == nil {
		t.Fatal("expected inner hits != nil")
	}
	if want, have := 1, len(innerHits.Hits.Hits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	if want, have := "1", innerHits.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}
}

func TestSearchInnerHitsOnNested(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	ctx := context.Background()

	// Create index
	createIndex, err := client.CreateIndex(testIndexName5).Body(`{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings": {
			"properties": {
			  	"comments": {
					"type": "nested"
				}
			}
		}
	}`).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Add documents
	// See https://www.elastic.co/guide/en/elasticsearch/reference/7.9/inner-hits.html#nested-inner-hits for example code.
	type comment struct {
		Author string `json:"author"`
		Number int    `json:"number"`
	}
	type doc struct {
		Title    string    `json:"title"`
		Comments []comment `json:"comments"`
	}
	doc1 := doc{
		Title: "Test title",
		Comments: []comment{
			{Author: "kimchy", Number: 1},
			{Author: "nik9000", Number: 2},
		},
	}
	_, err = client.Index().Index(testIndexName5).Id("1").BodyJson(&doc1).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Search for all documents that have an answer, and return those answers as inner hits
	q := NewNestedQuery("comments", NewMatchQuery("comments.number", 2)).InnerHit(NewInnerHit())
	searchResult, err := client.Search().
		Index(testIndexName5).
		Query(q).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 1 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 2, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 1 {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", 2, len(searchResult.Hits.Hits))
	}

	hit := searchResult.Hits.Hits[0]
	if want, have := "1", hit.Id; want != have {
		t.Fatalf("expected tweet %q; got: %q", want, have)
	}
	if hit.InnerHits == nil {
		t.Fatalf("expected inner hits; got: %v", hit.InnerHits)
	}
	if want, have := 1, len(hit.InnerHits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	innerHits, found := hit.InnerHits["comments"]
	if !found {
		t.Fatalf("expected inner hits for name %q", "comments")
	}
	if innerHits == nil || innerHits.Hits == nil {
		t.Fatal("expected inner hits != nil")
	}
	if want, have := 1, len(innerHits.Hits.Hits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	if want, have := "1", innerHits.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}
}

func TestSearchInnerHitsOnNestedHierarchy(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	ctx := context.Background()

	// Create index
	createIndex, err := client.CreateIndex(testIndexName5).Body(`{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings": {
			"properties": {
			  	"comments": {
					"type": "nested",
					"properties": {
				  		"votes": {
							"type": "nested"
				  		}
					}
			  	}
			}
		}
	}`).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Add documents
	// See https://www.elastic.co/guide/en/elasticsearch/reference/7.9/inner-hits.html#hierarchical-nested-inner-hits for example code.
	type vote struct {
		Voter string `json:"voter"`
		Value int    `json:"value"`
	}
	type comment struct {
		Author string `json:"author"`
		Text   string `json:"text"`
		Votes  []vote `json:"votes"`
	}
	type doc struct {
		Title    string    `json:"title"`
		Comments []comment `json:"comments"`
	}
	doc1 := doc{
		Title: "Test title",
		Comments: []comment{
			{Author: "kimchy", Text: "words words words", Votes: []vote{}},
			{Author: "nik9000", Text: "words words words", Votes: []vote{{Voter: "kimchy", Value: 1}, {Voter: "other", Value: -1}}},
		},
	}
	_, err = client.Index().Index(testIndexName5).Id("1").BodyJson(&doc1).Refresh("true").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Search for all documents that have an answer, and return those answers as inner hits
	q := NewNestedQuery("comments.votes", NewMatchQuery("comments.votes.voter", "kimchy")).InnerHit(NewInnerHit())
	searchResult, err := client.Search().
		Index(testIndexName5).
		Query(q).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 1 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 2, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 1 {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", 2, len(searchResult.Hits.Hits))
	}

	hit := searchResult.Hits.Hits[0]
	if want, have := "1", hit.Id; want != have {
		t.Fatalf("expected tweet %q; got: %q", want, have)
	}
	if hit.InnerHits == nil {
		t.Fatalf("expected inner hits; got: %v", hit.InnerHits)
	}
	if want, have := 1, len(hit.InnerHits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	innerHits, found := hit.InnerHits["comments.votes"]
	if !found {
		t.Fatalf("expected inner hits for name %q", "comments.votes")
	}
	if innerHits == nil || innerHits.Hits == nil {
		t.Fatal("expected inner hits != nil")
	}
	if want, have := 1, len(innerHits.Hits.Hits); want != have {
		t.Fatalf("expected %d inner hits; got: %d", want, have)
	}
	if want, have := "1", innerHits.Hits.Hits[0].Id; want != have {
		t.Fatalf("expected inner hit with id %q; got: %q", want, have)
	}
}

func TestSearchBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Indices  []string
		Types    []string
		Expected string
	}{
		{
			[]string{},
			[]string{},
			"/_search",
		},
		{
			[]string{"index1"},
			[]string{},
			"/index1/_search",
		},
		{
			[]string{"index1", "index2"},
			[]string{},
			"/index1%2Cindex2/_search",
		},
		{
			[]string{},
			[]string{"type1"},
			"/_all/type1/_search",
		},
		{
			[]string{"index1"},
			[]string{"type1"},
			"/index1/type1/_search",
		},
		{
			[]string{"index1", "index2"},
			[]string{"type1", "type2"},
			"/index1%2Cindex2/type1%2Ctype2/_search",
		},
		{
			[]string{},
			[]string{"type1", "type2"},
			"/_all/type1%2Ctype2/_search",
		},
	}

	for i, test := range tests {
		path, _, err := client.Search().Index(test.Indices...).Type(test.Types...).buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}

func TestSearchFilterPath(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Match all should return all documents
	all := NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(all).
		FilterPath(
			"took",
			"hits.hits._id",
			"hits.hits._source.user",
			"hits.hits._source.message",
		).
		Timeout("1s").
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Fatalf("expected SearchResult.Hits != nil; got nil")
	}
	// 0 because it was filtered out
	if want, got := int64(0), searchResult.TotalHits(); want != got {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if want, got := 3, len(searchResult.Hits.Hits); want != got {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}

	for _, hit := range searchResult.Hits.Hits {
		if want, got := "", hit.Index; want != got {
			t.Fatalf("expected index %q, got %q", want, got)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
		// user field
		v, found := item["user"]
		if !found {
			t.Fatalf("expected SearchResult.Hits.Hit[%q] to be found", "user")
		}
		if v == "" {
			t.Fatalf("expected user field, got %v (%T)", v, v)
		}
		// No retweets field
		v, found = item["retweets"]
		if found {
			t.Fatalf("expected SearchResult.Hits.Hit[%q] to not be found, got %v", "retweets", v)
		}
		if v == "" {
			t.Fatalf("expected user field, got %v (%T)", v, v)
		}
	}
}

func TestSearchAfter(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndLog(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User: "olivere", Retweets: 108,
		Message: "Welcome to Golang and Elasticsearch.",
		Created: time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User: "olivere", Retweets: 0,
		Message: "Another unrelated topic.",
		Created: time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User: "sandrae", Retweets: 12,
		Message: "Cycling is fun.",
		Created: time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

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

	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		SearchAfter("olivere").
		Sort("user", true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", 3, searchResult.TotalHits())
	}
	if want, got := 1, len(searchResult.Hits.Hits); want != got {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got: %d", want, got)
	}
	hit := searchResult.Hits.Hits[0]
	if want, got := "3", hit.Id; want != got {
		t.Fatalf("expected tweet %q; got: %q", want, got)
	}
}

func TestSearchResultWithFieldCollapsing(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Collapse(NewCollapseBuilder("user")).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if searchResult.Hits == nil {
		t.Fatalf("expected SearchResult.Hits != nil; got nil")
	}
	if got := searchResult.TotalHits(); got == 0 {
		t.Fatalf("expected SearchResult.TotalHits() > 0; got %d", got)
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
		if len(hit.Fields) == 0 {
			t.Fatal("expected fields in SearchResult")
		}
		usersVal, ok := hit.Fields["user"]
		if !ok {
			t.Fatalf("expected %q field in fields of SearchResult", "user")
		}
		users, ok := usersVal.([]interface{})
		if !ok {
			t.Fatalf("expected slice of strings in field of SearchResult, got %T", usersVal)
		}
		if len(users) != 1 {
			t.Fatalf("expected 1 entry in users slice, got %d", len(users))
		}
	}
}

func TestSearchResultWithFieldCollapsingAndInnerHits(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Collapse(
			NewCollapseBuilder("user").
				InnerHit(
					NewInnerHit().Name("last_tweets").Size(5).Sort("created", true),
				).
				MaxConcurrentGroupRequests(4)).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if searchResult.Hits == nil {
		t.Fatalf("expected SearchResult.Hits != nil; got nil")
	}
	if got := searchResult.TotalHits(); got == 0 {
		t.Fatalf("expected SearchResult.TotalHits() > 0; got %d", got)
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
		if len(hit.Fields) == 0 {
			t.Fatal("expected fields in SearchResult")
		}
		usersVal, ok := hit.Fields["user"]
		if !ok {
			t.Fatalf("expected %q field in fields of SearchResult", "user")
		}
		users, ok := usersVal.([]interface{})
		if !ok {
			t.Fatalf("expected slice of strings in field of SearchResult, got %T", usersVal)
		}
		if len(users) != 1 {
			t.Fatalf("expected 1 entry in users slice, got %d", len(users))
		}
		lastTweets, ok := hit.InnerHits["last_tweets"]
		if !ok {
			t.Fatalf("expected inner_hits named %q in SearchResult", "last_tweets")
		}
		if lastTweets == nil {
			t.Fatal("expected inner_hits in SearchResult")
		}
	}
}

func TestSearchScriptQuery(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	// ES uses Painless as default scripting engine in 6.x
	// Another example of using painless would be:
	//
	//	script := NewScript(`
	//	String username = doc['user'].value;
	//	return username == 'olivere'
	//`)
	// See https://www.elastic.co/guide/en/elasticsearch/painless/6.7/painless-examples.html
	script := NewScript("doc['user'].value == 'olivere'")
	query := NewScriptQuery(script)

	searchResult, err := client.Search().
		Index(testIndexName).
		Query(query).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if want, have := int64(2), searchResult.TotalHits(); want != have {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, have)
	}
	if want, have := 2, len(searchResult.Hits.Hits); want != have {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, have)
	}
}

func TestSearchWithDocvalueFields(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		DocvalueFields("user", "retweets").
		DocvalueFieldsWithFormat(
			// DocvalueField{Field: "user"},
			// DocvalueField{Field: "retweets", Format: "long"},
			DocvalueField{Field: "created", Format: "epoch_millis"},
		).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(3); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 3; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			t.Errorf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
		}
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchWithDateMathIndices(t *testing.T) {
	client := setupTestClient(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	ctx := context.Background()
	now := time.Now().UTC()
	indexNameToday := fmt.Sprintf("elastic-trail-%s", now.Format("2006.01.02"))
	indexNameYesterday := fmt.Sprintf("elastic-trail-%s", now.AddDate(0, 0, -1).Format("2006.01.02"))
	indexNameTomorrow := fmt.Sprintf("elastic-trail-%s", now.AddDate(0, 0, +1).Format("2006.01.02"))

	const mapping = `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	}
}`

	// Create indices
	for i, indexName := range []string{indexNameToday, indexNameTomorrow, indexNameYesterday} {
		_, err := client.CreateIndex(indexName).Body(mapping).Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteIndex(indexName).Do(ctx)

		// Add a document
		id := fmt.Sprintf("%d", i+1)
		_, err = client.Index().Index(indexName).Id(id).BodyJson(map[string]interface{}{
			"index": indexName,
		}).Refresh("wait_for").Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Count total
	cnt, err := client.
		Count(indexNameYesterday, indexNameToday, indexNameTomorrow).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if cnt != 3 {
		t.Fatalf("expected Count=%d; got %d", 3, cnt)
	}

	// Match all should return all documents
	res, err := client.Search().
		Index("<elastic-trail-{now/d}>", "<elastic-trail-{now-1d/d}>").
		Query(NewMatchAllQuery()).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := res.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(res.Hits.Hits), 2; got != want {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}
}

func TestSearchResultDecode(t *testing.T) {
	tests := []struct {
		Body string
	}{
		// #0 With _shards.failures
		{
			Body: `{
				"took":1146,
				"timed_out":false,
				"_shards":{
				   "total":8,
				   "successful":6,
				   "skipped":0,
				   "failed":2,
				   "failures":[
					  {
						 "shard":1,
						 "index":"l9leakip-0000001",
						 "node":"AsQq1Dh2QxCSTRSLTg0vFw",
						 "reason":{
							"type":"illegal_argument_exception",
							"reason":"The length [1119437] of field [events.summary] in doc[2524900]/index[l9leakip-0000001] exceeds the [index.highlight.max_analyzed_offset] limit [1000000]. To avoid this error, set the query parameter [max_analyzed_offset] to a value less than index setting [1000000] and this will tolerate long field values by truncating them."
						 }
					  },
					  {
						 "shard":3,
						 "index":"l9leakip-0000001",
						 "node":"AsQq1Dh2QxCSTRSLTg0vFw",
						 "reason":{
							"type":"illegal_argument_exception",
							"reason":"The length [1023566] of field [events.summary] in doc[2168434]/index[l9leakip-0000001] exceeds the [index.highlight.max_analyzed_offset] limit [1000000]. To avoid this error, set the query parameter [max_analyzed_offset] to a value less than index setting [1000000] and this will tolerate long field values by truncating them."
						 }
					  }
				   ]
				},
				"hits":{}
			 }`,
		},
	}

	for i, tt := range tests {
		var resp SearchResult
		if err := json.Unmarshal([]byte(tt.Body), &resp); err != nil {
			t.Fatalf("case #%d: expected no error, got %v", i, err)
		}
	}
}
