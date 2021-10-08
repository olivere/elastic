// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndicesClearCache(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	res, err := client.ClearCache().Do(context.Background())
	if err != nil {
		t.Fatalf("expected ClearCache to succeed, got: %v", err)
	}
	if res == nil {
		t.Fatalf("expected result to be != nil; got: %v", res)
	}
}

func TestIndicesClearCacheBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Indices  []string
		Expected string
	}{
		{
			[]string{},
			"/_cache/clear",
		},
		{
			[]string{"index1"},
			"/index1/_cache/clear",
		},
		{
			[]string{"index1", "index2"},
			"/index1%2Cindex2/_cache/clear",
		},
	}

	for i, test := range tests {
		path, _, err := client.ClearCache().Index(test.Indices...).buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}
