// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import "testing"

func TestXPackRollupGetBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		JobId        string
		ExpectedPath string
		ExpectErr    bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"my-job",
			"/_rollup/job/my-job",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackRollupGet(test.JobId)
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
			if path != test.ExpectedPath {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.ExpectedPath, path)
			}
		}
	}
}
