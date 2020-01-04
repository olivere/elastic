// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatIndices(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatIndices().Columns("*").Do(ctx)
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

// TestCatIndicesResponseRowAliasesMap tests if catIndicesResponseRowAliasesMap is declared
func TestCatIndicesResponseRowAliasesMap(t *testing.T) {
	if catIndicesResponseRowAliasesMap == nil {
		t.Fatal("want catIndicesResponseRowAliasesMap to be not nil")
	}

	if len(catIndicesResponseRowAliasesMap) == 0 {
		t.Fatal("want catIndicesResponseRowAliasesMap to be not empty")
	}
}

// TestCatIndicesWithAliases makes a simple test (if ?h=h will be the same as ?h=health)
func TestCatIndicesWithAliases(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatIndices().Columns("h").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].Health; have == "" {
		t.Fatalf("Index[0]: want != %q, have %q", "", have)
	}
}

// TestCatIndicesWithAliases makes a test with a double-alias
// asking `?h=rti` will fill one of the refresh.external_time/refresh.time fields (depending on elasticsearch version)
func TestCatIndicesWithAliases_Double(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatIndices().Columns("rti").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}

	refreshTime := res[0].RefreshTime
	refreshExternalTime := res[0].RefreshExternalTime

	if refreshTime == "" && refreshExternalTime == "" {
		t.Fatalf("Index[0]: want one of [refreshTime or refreshExternalTime] be not empty")
	}
}
