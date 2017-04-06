package elastic

import "testing"

func TestSnapshotPutRepositoryURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository string
		Expected   string
	}{
		{
			"repo",
			"/_snapshot/repo",
		},
	}

	for _, test := range tests {
		path, _, err := client.SnapshotCreateRepository(test.Repository).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}
