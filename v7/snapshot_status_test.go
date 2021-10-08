// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/url"
	"reflect"
	"testing"
)

func TestSnapshotStatusURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository     string
		Snapshot       []string
		MasterTimeout  string
		ExpectedPath   string
		ExpectedParams url.Values
	}{
		{
			Repository:    "repo",
			Snapshot:      []string{},
			MasterTimeout: "60s",
			ExpectedPath:  "/_snapshot/repo/_status",
			ExpectedParams: url.Values{
				"master_timeout": []string{"60s"},
			},
		},
		{
			Repository:    "repo",
			Snapshot:      []string{"snapA", "snapB"},
			MasterTimeout: "30s",
			ExpectedPath:  "/_snapshot/repo/snapA%2CsnapB/_status",
			ExpectedParams: url.Values{
				"master_timeout": []string{"30s"},
			},
		},
	}

	for _, test := range tests {
		path, params, err := client.SnapshotStatus().
			MasterTimeout(test.MasterTimeout).
			Repository(test.Repository).
			Snapshot(test.Snapshot...).
			buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.ExpectedPath {
			t.Errorf("expected %q; got: %q", test.ExpectedPath, path)
		}
		if !reflect.DeepEqual(params, test.ExpectedParams) {
			t.Errorf("expected %q; got: %q", test.ExpectedParams, params)
		}
	}
}
