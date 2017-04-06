package elastic

import "testing"

func TestSnapshotVerifyRepositoryURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository string
		Expected   string
	}{
		{
			"repo",
			"/_snapshot/repo/_verify",
		},
	}

	for _, test := range tests {
		path, _, err := client.SnapshotVerifyRepository(test.Repository).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}
