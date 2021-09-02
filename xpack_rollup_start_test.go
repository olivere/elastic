// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestXPackRollupStartBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Expected  string
		ExpectErr bool
	}{
		{
			"/_rollup/job/my-job/_start",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackRollupStart("my-job")
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
