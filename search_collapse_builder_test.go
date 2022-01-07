// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestCollapseBuilderSource(t *testing.T) {
	b := NewCollapseBuilder("user").
		InnerHit(NewInnerHit().Name("last_tweets").Size(5).Sort("date", true)).
		MaxConcurrentGroupRequests(4)
	src, err := b.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"field":"user","inner_hits":[{"name":"last_tweets","size":5,"sort":[{"date":{"order":"asc"}}]}],"max_concurrent_group_searches":4}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCollapseBuilderSourceMultipleInnerHits(t *testing.T) {
	b := NewCollapseBuilder("user.id").
		InnerHit(NewInnerHit().Name("largest_responses").Size(3).Sort("http.response.bytes", false)).
		InnerHit(NewInnerHit().Name("most_recent").Size(4).Sort("@timestamp", false)).
		MaxConcurrentGroupRequests(5)
	src, err := b.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"field":"user.id","inner_hits":[{"name":"largest_responses","size":3,"sort":[{"http.response.bytes":{"order":"desc"}}]},{"name":"most_recent","size":4,"sort":[{"@timestamp":{"order":"desc"}}]}],"max_concurrent_group_searches":5}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
