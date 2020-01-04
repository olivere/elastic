// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/url"
	"reflect"
	"testing"
)

func TestSnapshotGetValidate(t *testing.T) {
	var client *Client

	err := NewSnapshotGetService(client).Validate()
	got := err.Error()
	expected := "missing required fields: [Repository]"
	if got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}

func TestSnapshotGetURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository        string
		Snapshot          []string
		MasterTimeout     string
		IgnoreUnavailable bool
		Verbose           bool
		ExpectedPath      string
		ExpectedParams    url.Values
	}{
		{
			Repository:        "repo",
			Snapshot:          []string{},
			MasterTimeout:     "60s",
			IgnoreUnavailable: true,
			Verbose:           true,
			ExpectedPath:      "/_snapshot/repo/_all",
			ExpectedParams: url.Values{
				"master_timeout":     []string{"60s"},
				"ignore_unavailable": []string{"true"},
				"verbose":            []string{"true"},
			},
		},
		{
			Repository:        "repo",
			Snapshot:          []string{"snapA", "snapB"},
			MasterTimeout:     "60s",
			IgnoreUnavailable: true,
			Verbose:           true,
			ExpectedPath:      "/_snapshot/repo/snapA%2CsnapB",
			ExpectedParams: url.Values{
				"master_timeout":     []string{"60s"},
				"ignore_unavailable": []string{"true"},
				"verbose":            []string{"true"},
			},
		},
	}

	for _, test := range tests {
		path, params, err := client.SnapshotGet(test.Repository).
			MasterTimeout(test.MasterTimeout).
			Snapshot(test.Snapshot...).
			IgnoreUnavailable(true).
			Verbose(true).
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
