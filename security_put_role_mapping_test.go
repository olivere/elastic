package opensearch

import (
	"testing"
)

func TestSecurityPutRoleMappingBuildURL(t *testing.T) {

	client := setupTestClient(t)

	tests := []struct {
		Name         string
		Body         interface{}
		ExpectedPath string
		ExpectErr    bool
	}{
		{
			"",
			nil,
			"",
			true,
		},
		{
			"my-role",
			nil,
			"",
			true,
		},
		{
			"",
			nil,
			"`{}`,",
			true,
		},
		{
			"my-role",
			`{}`,
			"/_plugins/_security/api/rolesmapping/my-role",
			false,
		},
	}

	for i, test := range tests {
		builder := client.SecurityPutRoleMapping(test.Name).Body(test.Body)
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
