// Copyright 2005-2018 Coveo Solution. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestXPackSecurityPutRoleBuildURL(t *testing.T) {
	client := setupTestClientForXpackSecurity(t)

	tests := []struct {
		Name         string
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
			"my-role",
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
			"my-role",
			`{}`,
			"/_xpack/security/role/my-role",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackSecurityPutRole(test.Name).Body(test.Body)
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
