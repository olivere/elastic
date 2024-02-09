package opensearch

import (
	"testing"
)

func TestAlertingPostMonitorBuildURL(t *testing.T) {

	client := setupTestClient(t)

	tests := []struct {
		Body         any
		ExpectedPath string
		ExpectErr    bool
	}{
		{
			nil,
			"",
			true,
		},
		{
			`{}`,
			"/_plugins/_alerting/monitors",
			false,
		},
	}

	for i, test := range tests {
		builder := client.AlertingPostMonitor().Body(test.Body)

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
