// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestXPackWatcherDeactivateWatchBuildURL(t *testing.T) {
	client := setupTestClient(t) // , SetURL("http://elastic:elastic@localhost:9210"))

	tests := []struct {
		WatchId   string
		Expected  string
		ExpectErr bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"my-watch",
			"/_xpack/watcher/watch/my-watch/_deactivate",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackWatchDeactivate(test.WatchId)
		err := builder.Validate()
		if err != nil {
			if !test.ExpectErr {
				t.Errorf("case #%d: %v", i+1, err)
				continue
			}
		} else {
			// err == nil
			if test.ExpectErr {
				t.Errorf("case #%d: expected error", i+1)
				continue
			}
			path, _, _ := builder.buildURL()
			if path != test.Expected {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
			}
		}
	}
}
