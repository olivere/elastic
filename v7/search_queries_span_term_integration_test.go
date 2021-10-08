// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestSpanTermQueryIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t) //, SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	_, err := client.Search().
		Index(testIndexName).
		Query(NewSpanTermQuery("message", "Golang")).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
