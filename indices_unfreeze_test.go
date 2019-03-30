// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestIndicesUnfreezeBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Index    string
		Expected string
	}{
		{
			"index1",
			"/index1/_unfreeze",
		},
	}

	for i, test := range tests {
		path, _, err := client.UnfreezeIndex(test.Index).buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}
