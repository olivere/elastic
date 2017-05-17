package elastic

import (
	"reflect"
	"testing"
)

func TestSnapshotValidate(t *testing.T) {
	var client *Client

	err := NewSnapshotCreateService(client).Validate()
	got := err.Error()
	expected := "missing required fields: [Repository Name Indices]"
	if got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}

func TestSnapshotPutURL(t *testing.T) {
	var client *Client

	tests := []struct {
		Repository string
		Name       string
		Expected   string
	}{
		{
			"repo",
			"backup_of_sunday",
			"/_snapshot/repo/backup_of_sunday",
		},
	}

	for _, test := range tests {
		service := NewSnapshotCreateService(client).
			Repository(test.Repository).
			Name(test.Name)

		path, _, err := service.buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}

func TestSnapshotPutBody(t *testing.T) {
	var client *Client

	setupTestService := func(client *Client) *SnapshotCreateService {
		return NewSnapshotCreateService(client).
			IgnoreUnavailable(true).
			Partial(true)
	}

	tests := []struct {
		Service  *SnapshotCreateService
		Expected interface{}
	}{
		{
			setupTestService(client).Indices(testIndexName),
			map[string]interface{}{
				"indices":            testIndexName,
				"ignore_unavailable": true,
				"partial":            true,
			},
		},
		{
			setupTestService(client).Indices(testIndexName2).IncludeGlobalState(false),
			map[string]interface{}{
				"indices":              testIndexName2,
				"ignore_unavailable":   true,
				"partial":              true,
				"include_global_state": false,
			},
		},
	}

	for _, test := range tests {
		got, err := test.Service.buildBody()
		if err != nil {
			t.Fatalf("buildBody failed: %v", err)
		}
		if !reflect.DeepEqual(got, test.Expected) {
			t.Fatalf("got %v, want %v", got, test.Expected)
		}
	}
}
