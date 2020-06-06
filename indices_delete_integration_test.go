// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndicesDeleteIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndex(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	before, err := client.IndexNames()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.DeleteIndex(testIndexName, testIndexNameEmpty).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	after, err := client.IndexNames()
	if err != nil {
		t.Fatal(err)
	}

	if want, have := len(after), len(before)-2; want != have {
		t.Fatalf("expected %d indices, got %d", want, have)
	}
}
