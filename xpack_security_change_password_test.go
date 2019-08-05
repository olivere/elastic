// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestXPackSecurityChangePasswordBuildURL(t *testing.T) {
	client := setupTestClientForXpackSecurity(t)

	tests := []struct {
		Username       string
		Password       string
		Refresh        string
		Body           interface{}
		ExpectedPath   string
		ExpectedParams url.Values
		ExpectedErr    bool
	}{
		// #0 No username
		{
			Username:       "",
			Password:       "",
			Refresh:        "",
			Body:           nil,
			ExpectedPath:   "",
			ExpectedParams: url.Values{},
			ExpectedErr:    true,
		},
		// #1 No username (but body)
		{
			Username:       "",
			Password:       "",
			Refresh:        "",
			Body:           `{}`,
			ExpectedPath:   "",
			ExpectedParams: url.Values{},
			ExpectedErr:    true,
		},
		// #2 No body or password
		{
			Username:       "my-user",
			Password:       "",
			Refresh:        "",
			Body:           nil,
			ExpectedPath:   "",
			ExpectedParams: url.Values{},
			ExpectedErr:    true,
		},
		// #3 No password but body
		{
			Username:       "my-user",
			Password:       "",
			Refresh:        "",
			Body:           `{"password":"secret"}`,
			ExpectedPath:   "/_xpack/security/user/my-user/_password",
			ExpectedParams: url.Values{},
			ExpectedErr:    false,
		},
		// #4 No body but password
		{
			Username:       "my-user",
			Password:       "secret",
			Refresh:        "",
			Body:           nil,
			ExpectedPath:   "/_xpack/security/user/my-user/_password",
			ExpectedParams: url.Values{},
			ExpectedErr:    false,
		},
		// #5 With refresh option
		{
			Username:     "my-user",
			Password:     "secret",
			Refresh:      "wait_for",
			Body:         nil,
			ExpectedPath: "/_xpack/security/user/my-user/_password",
			ExpectedParams: url.Values{
				"refresh": []string{"wait_for"},
			},
			ExpectedErr: false,
		},
	}

	for i, tt := range tests {
		builder := client.XPackSecurityChangePassword(tt.Username).
			Password(tt.Password).
			Refresh(tt.Refresh).
			Body(tt.Body)
		err := builder.Validate()
		if err != nil {
			if !tt.ExpectedErr {
				t.Errorf("case #%d: %v", i, err)
				continue
			}
		} else {
			// err == nil
			if tt.ExpectedErr {
				t.Errorf("case #%d: expected error", i)
				continue
			}
			path, params, err := builder.buildURL()
			if err != nil {
				t.Fatalf("case #%d: %v", i, err)
			}
			if path != tt.ExpectedPath {
				t.Errorf("case #%d: expected %q; got: %q", i, tt.ExpectedPath, path)
			}
			if want, have := tt.ExpectedParams, params; !cmp.Equal(want, have) {
				t.Errorf("case #%d: want Params=%#v, have %#v\n\tdiff: %s", i, want, have, cmp.Diff(want, have))
			}
		}
	}
}
