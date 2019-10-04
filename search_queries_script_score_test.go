// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestScriptScoreQuery(t *testing.T) {
	q := NewScriptScoreQuery(
		NewMatchQuery("message", "elasticsearch"),
		NewScript("doc['likes'].value / 10"),
	).MinScore(1.1).Boost(5.0).QueryName("my_query")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"script_score":{"_name":"my_query","boost":5,"min_score":1.1,"query":{"match":{"message":{"query":"elasticsearch"}}},"script":{"source":"doc['likes'].value / 10"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestScriptScoreQueryIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	res, err := client.Search().
		Index(testIndexName).
		Query(
			NewScriptScoreQuery(
				NewMatchQuery("message", "Golang"),
				NewScript("(1 + doc['retweets'].value) * 10"),
			),
		).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.Hits == nil {
		t.Errorf("expected Hits != nil; got nil")
	}
	if want, have := int64(1), res.TotalHits(); want != have {
		t.Errorf("expected TotalHits() = %d; got %d", want, have)
	}
	if want, have := 1, len(res.Hits.Hits); want != have {
		t.Errorf("expected len(Hits.Hits) = %d; got %d", want, have)
	}

	hit := res.Hits.Hits[0]

	if want, have := testIndexName, hit.Index; want != have {
		t.Fatalf("expected Hits.Hit.Index = %q; got %q", want, have)
	}
	if want, have := "1", hit.Id; want != have {
		t.Fatalf("expected Hits.Hit.Id = %q; got %q", want, have)
	}
	if hit.Score == nil {
		t.Fatal("expected Hits.Hit.Score != nil")
	}
	if want, have := 10.0, *hit.Score; want != have {
		t.Fatalf("expected Hits.Hit.Score = %v; got %v", want, have)
	}
	var tw tweet
	if err := json.Unmarshal(hit.Source, &tw); err != nil {
		t.Fatal(err)
	}
	if want, have := "olivere", tw.User; want != have {
		t.Fatalf("expected User = %q; got %q", want, have)
	}
}
