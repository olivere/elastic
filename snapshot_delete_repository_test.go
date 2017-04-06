package elastic

import "testing"

func TestSnapshotDeleteRepositoryURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Repository []string
		Expected   string
	}{
		{
			[]string{"repo1"},
			"/_snapshot/repo1",
		},
		{
			[]string{"repo1", "repo2"},
			"/_snapshot/repo1%2Crepo2",
		},
	}

	for _, test := range tests {
		path, _, err := client.SnapshotDeleteRepository(test.Repository...).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}
