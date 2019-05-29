package elastic

import (
	"net/url"
	"reflect"
	"testing"
)

func TestSnapshotRestoreValidate(t *testing.T) {
	expected := "missing required fields: [Repository Snapshot]"

	if got := NewSnapshotRestoreService(new(Client)).Validate().Error(); got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}

func TestSnapshotRestorePostURL(t *testing.T) {
	test := struct {
		Repository        string
		Snapshot          string
		Pretty            bool
		MasterTimeout     string
		WaitForCompletion bool
		IgnoreUnavailable bool
		ExpectedPath      string
		ExpectedParams    url.Values
	}{

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
	}

	path, params, err := new(Client).SnapshotRestore(test.Repository, test.Snapshot).
		Pretty(test.Pretty).
		MasterTimeout(test.MasterTimeout).
		WaitForCompletion(test.WaitForCompletion).
		IgnoreUnavailable(test.IgnoreUnavailable).
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

func TestSnapshotRestoreBuildBody(t *testing.T) {
	test := struct {
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
	}

	body := new(Client).SnapshotRestore(test.Repository, test.Snapshot).
		Partial(test.Partial).
		IncludeAliases(test.IncludeAliases).
		IncludeGlobalState(test.IncludeGlobalState).
		RenamePattern(test.RenamePattern).
		RenameReplacement(test.RenameReplacement).
		Indices(test.Indices...).
		IndexSettings(test.IndexSettings).
		buildBody()

	if !reflect.DeepEqual(body, test.ExpectedBody) {
		t.Errorf("expected %q; got: %q", test.ExpectedBody, body)
	}
}
