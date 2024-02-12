package opensearch

import (
	"testing"

	"k8s.io/utils/ptr"
)

func TestAlertingSearchMonitorBuildURL(t *testing.T) {

	client := setupTestClient(t)

	tests := []struct {
		Name         *string
		Search       *string
		ExpectedPath string
		ExpectErr    bool
	}{
		{
			nil,
			nil,
			"",
			true,
		},
		{
			ptr.To[string]("my-monitor"),
			nil,
			"/_plugins/_alerting/monitors/_search",
			false,
		},
		{
			nil,
			ptr.To[string]("my-search"),
			"/_plugins/_alerting/monitors/_search",
			false,
		},
	}

	var builder *AlertingSearchMonitorService
	for i, test := range tests {
		builder = client.AlertingSearchMonitor()
		if test.Name != nil {
			builder.SearchByName(*test.Name)
		}
		if test.Search != nil {
			builder.Search(test.Search)
		}

		err := builder.Validate()
		if err != nil {
			if !test.ExpectErr {
				t.Errorf("case #%d: %v", i+1, err)
				continue
			}
		} else {
			// err == nil
			if test.ExpectErr {
				t.Errorf("case #%d: expected error", i+1)
				continue
			}
			path, _, _ := builder.buildURL()
			if path != test.ExpectedPath {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.ExpectedPath, path)
			}
		}
	}
}
