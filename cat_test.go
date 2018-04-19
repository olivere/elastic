// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatIndices(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	res, err := client.Cat().Indices().Pretty(true).Do(context.TODO())
	if err != nil {
		t.Fatalf("expected to not get an error, got %v", err)
	}

	if res == "" {
		t.Fatalf("expected res to not be an empty string")
	}
}
