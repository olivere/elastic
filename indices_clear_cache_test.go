package elastic

import (
	"context"
	"testing"
)

func TestClearCache(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	_, err := client.ClearCache().Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndicesClearCacheBuildURL(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		Indices  []string
		Expected string
	}{
		{
			[]string{},
			"/_cache/clear",
		},
		{
			[]string{"index1"},
			"/index1/_cache/clear",
		},
		{
			[]string{"index1", "index2"},
			"/index1%2Cindex2/_cache/clear",
		},
	}

	for i, test := range tests {
		path, _, err := client.ClearCache().Index(test.Indices...).buildURL()
		if err != nil {
			t.Errorf("case #%d: %v", i+1, err)
			continue
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}
