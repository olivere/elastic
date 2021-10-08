// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestAggsBucketTermsIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	resp, err := client.Search().
		Index(testIndexName).
		Query(NewMatchAllQuery()).
		Aggregation("retweets", NewTermsAggregation().Field("retweets")).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Hits == nil {
		t.Errorf("expected Hits != nil")
	}
	if want, have := int64(3), resp.TotalHits(); want != have {
		t.Errorf("TotalHits(): want %d, have %d", want, have)
	}

	aggs := resp.Aggregations
	if aggs == nil {
		t.Fatalf("expected Aggregations != nil")
	}

	agg, found := aggs.Terms("retweets")
	if !found {
		t.Fatal("expected to find aggregation by name")
	}
	if agg == nil {
		t.Fatalf("expected != nil")
	}
	if want, have := 3, len(agg.Buckets); want != have {
		t.Fatalf("len(Buckets): want %d, have %d", want, have)
	}

	// Element #0
	if want, have := float64(0), agg.Buckets[0].Key; want != have {
		t.Errorf("Buckets[0].Key: want %v, have %v", want, have)
	}
	if want, have := json.Number("0"), agg.Buckets[0].KeyNumber; want != have {
		t.Errorf("agg.Buckets[0].KeyNumber: want %q, have %q", want, have)
	}
	if have, err := agg.Buckets[0].KeyNumber.Int64(); err != nil {
		t.Errorf("Buckets[0].KeyNumber.Int64(): %v", err)
	} else if want := int64(0); want != have {
		t.Errorf("Buckets[0].KeyNumber.Int64(): want %d, have %d", want, have)
	}
	if want, have := "0", agg.Buckets[0].KeyNumber.String(); want != have {
		t.Errorf("Buckets[0].KeyNumber.String(): want %q, have %q", want, have)
	}
	if want, have := int64(1), agg.Buckets[0].DocCount; want != have {
		t.Errorf("Buckets[0].DocCount: want %d, have %d", want, have)
	}

	// Element #1
	if want, have := float64(12), agg.Buckets[1].Key; want != have {
		t.Errorf("Buckets[1].Key: want %v, have %v", want, have)
	}
	if have, err := agg.Buckets[1].KeyNumber.Int64(); err != nil {
		t.Errorf("Buckets[1].KeyNumber.Int64(): %v", err)
	} else if want := int64(12); want != have {
		t.Errorf("Buckets[1].KeyNumber.Int64(): want %d, have %d", want, have)
	}
	if want, have := "12", agg.Buckets[1].KeyNumber.String(); want != have {
		t.Errorf("Buckets[1].KeyNumber.String(): want %q, have %q", want, have)
	}
	if want, have := int64(1), agg.Buckets[1].DocCount; want != have {
		t.Errorf("Buckets[1].DocCount: want %d, have %d", want, have)
	}

	// Element #2
	if want, have := float64(108), agg.Buckets[2].Key; want != have {
		t.Errorf("Buckets[2].Key: want %v, have %v", want, have)
	}
	if want, have := json.Number("108"), agg.Buckets[2].KeyNumber; want != have {
		t.Errorf("Buckets[2].KeyNumber: want %q, have %q", want, have)
	}
	if have, err := agg.Buckets[2].KeyNumber.Int64(); err != nil {
		t.Errorf("Buckets[2].KeyNumber.Int64(): %v", err)
	} else if want := int64(108); want != have {
		t.Errorf("Buckets[2].KeyNumber.Int64(): want %d, have %d", want, have)
	}
	if want, have := "108", agg.Buckets[2].KeyNumber.String(); want != have {
		t.Errorf("Buckets[2].KeyNumber.String(): want %q, have %q", want, have)
	}
	if want, have := int64(1), agg.Buckets[2].DocCount; want != have {
		t.Errorf("Buckets[2].DocCount: want %d, have %d", want, have)
	}

}
