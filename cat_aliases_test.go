// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatAliases(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()

	// Add two aliases
	_, err := client.Alias().
		Add(testIndexName, testAliasName).
		Action(NewAliasAddAction(testAliasName).Index(testIndexName2)).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// Remove aliases
		client.Alias().
			Remove(testIndexName, testAliasName).
			Remove(testIndexName2, testAliasName).
			Do(context.TODO())
	}()

	// Check the response
	res, err := client.CatAliases().Pretty(true).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if want, have := 2, len(res); want != have {
		t.Fatalf("want len=%d, have %d", want, have)
	}
}
