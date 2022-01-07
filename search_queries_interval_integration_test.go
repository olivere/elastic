// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestIntervalQuery_Integration(t *testing.T) {
	// client := setupTestClientAndCreateIndexAndAddDocs(t, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	t.Run("Match", func(t *testing.T) {
		testIntervalQueryMatch(client, t)
	})
	t.Run("Prefix", func(t *testing.T) {
		testIntervalQueryPrefix(client, t)
	})
	t.Run("Wildcard", func(t *testing.T) {
		testIntervalQueryWildcard(client, t)
	})
	t.Run("Fuzzy", func(t *testing.T) {
		testIntervalQueryFuzzy(client, t)
	})
}

func testIntervalQueryMatch(client *Client, t *testing.T) {
	q := NewIntervalQuery(
		"message",
		NewIntervalQueryRuleAllOf(
			NewIntervalQueryRuleAnyOf(
				NewIntervalQueryRuleMatch("Golang").Ordered(true),
				NewIntervalQueryRuleMatch("Cycling").MaxGaps(0).Filter(
					NewIntervalQueryFilter().NotContaining(
						NewIntervalQueryRuleMatch("Hockey"),
					),
				),
			),
		).Ordered(true),
	)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(q).
		Size(10).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 2; got != want {
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

func testIntervalQueryPrefix(client *Client, t *testing.T) {
	q := NewIntervalQuery(
		"message",
		NewIntervalQueryRuleAllOf(
			NewIntervalQueryRuleAnyOf(
				NewIntervalQueryRuleMatch("Golang").Ordered(true),
				NewIntervalQueryRuleMatch("Cycling").MaxGaps(0).Filter(
					NewIntervalQueryFilter().NotContaining(
						NewIntervalQueryRulePrefix("Hockey"),
					),
				),
			),
		).Ordered(true),
	)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(q).
		Size(10).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 2; got != want {
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

func testIntervalQueryWildcard(client *Client, t *testing.T) {
	q := NewIntervalQuery(
		"message",
		NewIntervalQueryRuleAllOf(
			NewIntervalQueryRuleAnyOf(
				NewIntervalQueryRuleMatch("Golang").Ordered(true),
				NewIntervalQueryRuleMatch("Cycling").MaxGaps(0).Filter(
					NewIntervalQueryFilter().NotContaining(
						NewIntervalQueryRuleWildcard("Hockey*"),
					),
				),
			),
		).Ordered(true),
	)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(q).
		Size(10).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 2; got != want {
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

func testIntervalQueryFuzzy(client *Client, t *testing.T) {
	q := NewIntervalQuery(
		"message",
		NewIntervalQueryRuleAllOf(
			NewIntervalQueryRuleAnyOf(
				NewIntervalQueryRuleMatch("Golang").Ordered(true),
				NewIntervalQueryRuleMatch("Cycling").MaxGaps(0).Filter(
					NewIntervalQueryFilter().NotContaining(
						NewIntervalQueryRuleFuzzy("Hocky").Fuzziness("auto"),
					),
				),
			),
		).Ordered(true),
	)

	// Match all should return all documents
	searchResult, err := client.Search().
		Index(testIndexName).
		Query(q).
		Size(10).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := searchResult.TotalHits(), int64(2); got != want {
		t.Errorf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(searchResult.Hits.Hits), 2; got != want {
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
