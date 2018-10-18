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
	var client *Client

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
	}

	for _, test := range tests {
		path := NewSnapshotDeleteService(client).client.SnapshotDelete(
			test.Repository, test.Snapshot,
		).buildURL()

		if path != test.ExpectedPath {
			t.Errorf("expected %q; got: %q", test.ExpectedPath, path)
		}
	}
}
