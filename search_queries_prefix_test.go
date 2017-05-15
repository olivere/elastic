// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic_test

import (
	"context"
	"encoding/json"
	"testing"

	"gopkg.in/olivere/elastic.v5"
)

func ExamplePrefixQuery() {
	// Get a client to the local Elasticsearch instance.
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Define wildcard query
	q := elastic.NewPrefixQuery("user", "oli")
	q = q.QueryName("my_query_name")

	searchResult, err := client.Search().
		Index("twitter").
		Query(q).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	_ = searchResult
}

func TestPrefixQuery(t *testing.T) {
	q := elastic.NewPrefixQuery("user", "ki")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"prefix":{"user":"ki"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestPrefixQueryWithOptions(t *testing.T) {
	q := elastic.NewPrefixQuery("user", "ki")
	q = q.QueryName("my_query_name")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"prefix":{"user":{"_name":"my_query_name","value":"ki"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
