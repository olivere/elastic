// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestXPackRollupPutBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		JobId        string
		Body         interface{}
		ExpectedPath string
		ExpectErr    bool
	}{
		{
			"",
			nil,
			"",
			true,
		},
		{
			"my-job",
			nil,
			"",
			true,
		},
		{
			"",
			`{}`,
			"",
			true,
		},
		{
			"my-job",
			`{}`,
			"/_rollup/job/my-job",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackRollupPut(test.JobId).Body(test.Body)
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
