// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatShards(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatShards().Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].Index; have == "" {
		t.Fatalf("Index[0]: want != %q, have %q", "", have)
	}
}

// TestCatShardsResponseRowAliasesMap tests if catIndicesResponseRowAliasesMap is declared
func TestCatShardsResponseRowAliasesMap(t *testing.T) {
	if catIndicesResponseRowAliasesMap == nil {
		t.Fatal("want catIndicesResponseRowAliasesMap to be not nil")
	}
	if len(catIndicesResponseRowAliasesMap) == 0 {
		t.Fatal("want catIndicesResponseRowAliasesMap to be not empty")
	}
}

// TestCatShardsWithSpecificColumns makes a simple test with specific column names.
func TestCatShardsWithSpecificColumns(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatShards().Columns("index", "shard").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].Index; have == "" {
		t.Fatalf("Index[0]: want != %q, have %q", "", have)
	}
}
