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
		MasterTimeout     string
		IgnoreUnavailable bool
		Verbose           bool
		ExpectedPath      string
		ExpectedParams    url.Values
	}{
		{
			Repository:        "repo",
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
	}

	for _, test := range tests {
		path, params, err := client.SnapshotGet(test.Repository).
			MasterTimeout(test.MasterTimeout).
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
