// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/url"
	"testing"
)

func TestSnapshotRepoValidate(t *testing.T) {
	var client *Client

	err := NewSnapshotDeleteService(client).Validate()
	got := err.Error()
	expected := "missing required fields: [Repository Snapshot]"
	if got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}

func TestSnapshotDeleteURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository     string
		Snapshot       string
		ExpectedPath   string
		ExpectedParams url.Values
	}{
		{
			Repository:   "repo",
			Snapshot:     "snapshot_of_sunday",
			ExpectedPath: "/_snapshot/repo/snapshot_of_sunday",
		},
		{
			Repository:   "repö",
			Snapshot:     "001",
			ExpectedPath: "/_snapshot/rep%C3%B6/001",
		},
	}

	for _, tt := range tests {
		path, _, err := client.SnapshotDelete(tt.Repository, tt.Snapshot).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != tt.ExpectedPath {
			t.Errorf("expected %q; got: %q", tt.ExpectedPath, path)
		}
	}
}
