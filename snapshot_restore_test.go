// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSnapshotRestoreValidate(t *testing.T) {
	expected := "missing required fields: [Repository Snapshot]"
	if got := NewSnapshotRestoreService(nil).Validate().Error(); got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}

func TestSnapshotRestorePostURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository        string
		Snapshot          string
		Pretty            bool
		MasterTimeout     string
		WaitForCompletion bool
		IgnoreUnavailable bool
		ExpectedPath      string
		ExpectedParams    url.Values
	}{
		{
			Repository:        "repo",
			Snapshot:          "snapshot_of_sunday",
			Pretty:            true,
			MasterTimeout:     "60s",
			WaitForCompletion: true,
			IgnoreUnavailable: true,
			ExpectedPath:      "/_snapshot/repo/snapshot_of_sunday/_restore",
			ExpectedParams: url.Values{
				"pretty":              []string{"true"},
				"master_timeout":      []string{"60s"},
				"wait_for_completion": []string{"true"},
				"ignore_unavailable":  []string{"true"},
			},
		},
	}

	for _, tt := range tests {
		path, params, err := client.SnapshotRestore(tt.Repository, tt.Snapshot).
			Pretty(tt.Pretty).
			MasterTimeout(tt.MasterTimeout).
			WaitForCompletion(tt.WaitForCompletion).
			IgnoreUnavailable(tt.IgnoreUnavailable).
			buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != tt.ExpectedPath {
			t.Errorf("expected Path=%q; got: %q", tt.ExpectedPath, path)
		}
		if want, have := tt.ExpectedParams, params; !cmp.Equal(want, have) {
			t.Errorf("expected Params=%#v; got: %#v\ndiff: %s", want, have, cmp.Diff(want, have))
		}
	}
}

func TestSnapshotRestoreBuildBody(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository         string
		Snapshot           string
		Partial            bool
		IncludeAliases     bool
		IncludeGlobalState bool
		RenamePattern      string
		RenameReplacement  string
		Indices            []string
		IndexSettings      map[string]interface{}
		ExpectedBody       map[string]interface{}
	}{
		{
			Repository:         "repo",
			Snapshot:           "snapshot_of_sunday",
			Partial:            true,
			IncludeAliases:     true,
			IncludeGlobalState: true,
			RenamePattern:      "index_(.+)",
			RenameReplacement:  "restored_index_$1",
			Indices:            []string{"index_1", "indexe_2", "index_3"},
			IndexSettings: map[string]interface{}{
				"index.number_of_replicas": 0,
			},
			ExpectedBody: map[string]interface{}{
				"partial":              true,
				"include_aliases":      true,
				"include_global_state": true,
				"rename_pattern":       "index_(.+)",
				"rename_replacement":   "restored_index_$1",
				"indices":              "index_1,indexe_2,index_3",
				"index_settings": map[string]interface{}{
					"index.number_of_replicas": 0,
				},
			},
		},
	}

	for _, tt := range tests {
		body := client.SnapshotRestore(tt.Repository, tt.Snapshot).
			Partial(tt.Partial).
			IncludeAliases(tt.IncludeAliases).
			IncludeGlobalState(tt.IncludeGlobalState).
			RenamePattern(tt.RenamePattern).
			RenameReplacement(tt.RenameReplacement).
			Indices(tt.Indices...).
			IndexSettings(tt.IndexSettings).
			buildBody()

		if want, have := tt.ExpectedBody, body; !cmp.Equal(want, have) {
			t.Errorf("expected Body=%s; got: %s\ndiff: %s", want, have, cmp.Diff(want, have))
		}
	}
}
