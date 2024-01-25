package opensearch

import (
	"testing"
)

func TestSmPutPolicyBuildURL(t *testing.T) {

	client := setupTestClient(t)

	tests := []struct {
		Name         string
		Body         any
		ExpectedPath string
		ExpectErr    bool
		SeqNum       int64
		PrimaryTerm  int64
	}{
		{
			"",
			nil,
			"",
			true,
			0,
			1,
		},
		{
			"my-policy",
			nil,
			"",
			true,
			0,
			1,
		},
		{
			"",
			`{}`,
			"",
			true,
			0,
			1,
		},
		{
			"my-policy",
			`{}`,
			"/_plugins/_sm/policies/my-policy?if_seq_no=0&if_primary_term=1",
			true,
			0,
			0,
		},
		{
			"my-policy",
			`{}`,
			"/_plugins/_sm/policies/my-policy?if_seq_no=0&if_primary_term=1",
			false,
			0,
			1,
		},
		{
			"my-policy",
			`{}`,
			"/_plugins/_sm/policies/my-policy?if_seq_no=10&if_primary_term=1",
			false,
			10,
			1,
		},
	}

	for i, test := range tests {
		builder := client.SmPutPolicy(test.Name).Body(test.Body).SequenceNumber(test.SeqNum).PrimaryTerm(test.PrimaryTerm)
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
