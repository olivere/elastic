// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatFielddata(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()

	// generate fielddata by aggregation
	aggRes, err := client.Search(testIndexName).
		Aggregation("gene_message_fielddata", NewTermsAggregation().Field("message")).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if aggRes == nil {
		t.Fatal("want response, have nil")
	}

	res, err := client.CatFielddata().Pretty(true).Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}

	// check fielddata "message" in response
	var found bool
	for _, fielddata := range res {
		if fielddata.Field == "message" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("fielddata message not found")
	}

}
