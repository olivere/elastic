// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestSyncedFlush(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)
	//client := setupTestClientAndCreateIndexAndLog(t)

	// Sync Flush all indices
	res, err := client.SyncedFlush().Pretty(true).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Errorf("expected res to be != nil; got: %v", res)
	}
}

func TestSyncedFlushBuildURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		Indices               []string
		Expected              string
		ExpectValidateFailure bool
	}{
		{
			[]string{},
			"/_flush/synced",
			false,
		},
		{
			[]string{"index1"},
			"/index1/_flush/synced",
			false,
		},
		{
			[]string{"index1", "index2"},
			"/index1%2Cindex2/_flush/synced",
			false,
		},
	}

	for i, test := range tests {
		err := NewIndicesSyncedFlushService(client).Index(test.Indices...).Validate()
		if err == nil && test.ExpectValidateFailure {
			t.Errorf("case #%d: expected validate to fail", i+1)
			continue
		}
		if err != nil && !test.ExpectValidateFailure {
			t.Errorf("case #%d: expected validate to succeed", i+1)
			continue
		}
		if !test.ExpectValidateFailure {
			path, _, err := NewIndicesSyncedFlushService(client).Index(test.Indices...).buildURL()
			if err != nil {
				t.Fatalf("case #%d: %v", i+1, err)
			}
			if path != test.Expected {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
			}
		}
	}
}

func TestSyncedFlushResponse(t *testing.T) {
	js := `{
		"_shards": {
		   "total": 4,
		   "successful": 1,
		   "failed": 1
		},
		"twitter": {
		   "total": 4,
		   "successful": 3,
		   "failed": 1,
		   "failures": [
			  {
				 "shard": 1,
				 "reason": "unexpected error",
				 "routing": {
					"state": "STARTED",
					"primary": false,
					"node": "SZNr2J_ORxKTLUCydGX4zA",
					"relocating_node": null,
					"shard": 1,
					"index": "twitter"
				 }
			  }
		   ]
		}
	 }`

	var resp IndicesSyncedFlushResponse
	if err := json.Unmarshal([]byte(js), &resp); err != nil {
		t.Fatal(err)
	}
	if want, have := 4, resp.Shards.Total; want != have {
		t.Fatalf("want Shards.Total = %v, have %v", want, have)
	}
	if want, have := 1, resp.Shards.Successful; want != have {
		t.Fatalf("want Shards.Successful = %v, have %v", want, have)
	}
	if want, have := 1, resp.Shards.Failed; want != have {
		t.Fatalf("want Shards.Failed = %v, have %v", want, have)
	}

	{
		indexName := "twitter"
		index, found := resp.Index[indexName]
		if !found {
			t.Fatalf("want index %q", indexName)
		}
		if index == nil {
			t.Fatalf("want index %q", indexName)
		}
		if want, have := 4, index.Total; want != have {
			t.Fatalf("want Index[%q].Total = %v, have %v", indexName, want, have)
		}
		if want, have := 3, index.Successful; want != have {
			t.Fatalf("want Index[%q].Successful = %v, have %v", indexName, want, have)
		}
		if want, have := 1, index.Failed; want != have {
			t.Fatalf("want Index[%q].Failed = %v, have %v", indexName, want, have)
		}
		if want, have := 1, len(index.Failures); want != have {
			t.Fatalf("want len(Index[%q].Failures) = %v, have %v", indexName, want, have)
		}
		failure := index.Failures[0]
		if want, have := 1, failure.Shard; want != have {
			t.Fatalf("want Index[%q].Failures[0].Shard = %v, have %v", indexName, want, have)
		}
		if want, have := "unexpected error", failure.Reason; want != have {
			t.Fatalf("want Index[%q].Failures[0].Reason = %q, have %q", indexName, want, have)
		}
		if want, have := false, failure.Routing.Primary; want != have {
			t.Fatalf("want Index[%q].Failures[0].Routing.Primary = %v, have %v", indexName, want, have)
		}
		if want, have := "SZNr2J_ORxKTLUCydGX4zA", failure.Routing.Node; want != have {
			t.Fatalf("want Index[%q].Failures[0].Routing.Node = %q, have %q", indexName, want, have)
		}
		if have := failure.Routing.RelocatingNode; have != nil {
			t.Fatalf("want Index[%q].Failures[0].Routing.RelocatingNode = nil, have %v", indexName, have)
		}
	}

}
