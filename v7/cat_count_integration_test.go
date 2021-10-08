// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatCountIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	var (
		total   int64
		indices = []string{testIndexName, testIndexName2, testOrderIndex}
	)

	for _, index := range indices {
		count, err := client.Count(index).Do(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		total += count
	}

	resp, err := client.CatCount().
		Index(indices...).
		Columns("*").
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 1, len(resp); want != have {
		t.Fatalf("expected %d response item, got %d", want, have)
	}

	if want, have := total, int64(resp[0].Count); want != have {
		t.Fatalf("expected %d documents, got %d", want, have)
	}
}
