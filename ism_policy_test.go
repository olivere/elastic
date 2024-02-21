package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestIsmPolicy(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expecedIsmPolicy := &IsmPutPolicy{
		Policy: IsmPolicyBase{
			Description:  ptr.To[string]("ingesting logs"),
			DefaultState: ptr.To[string]("ingest"),
			States: []IsmPolicyState{
				{
					Name: "ingest",
					Actions: []map[string]any{
						{
							"rollover": map[string]any{
								"min_doc_count": 5,
							},
						},
					},
					Transitions: []IsmPolicyStateTransition{
						{
							StateName: "search",
						},
					},
				},
				{
					Name: "search",
					Transitions: []IsmPolicyStateTransition{
						{
							StateName: "delete",
							Conditions: map[string]any{
								"min_index_age": "5m",
							},
						},
					},
				},
				{
					Name: "delete",
					Actions: []map[string]any{
						{
							"delete": map[string]any{},
						},
					},
				},
			},
		},
	}

	// Create ISM policy
	resPut, err := client.IsmPutPolicy("test").Body(expecedIsmPolicy).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get ISM policy
	resGet, err := client.IsmGetPolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.NotNil(t, resGet.Policy)

	// Update ISM policy
	expecedIsmPolicy.Policy.Description = ptr.To[string]("test")
	_, err = client.IsmPutPolicy("test").Body(expecedIsmPolicy).SequenceNumber(resGet.SequenceNumber).PrimaryTerm(resGet.PrimaryTerm).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.IsmGetPolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", *resGet.Policy.Description)

	// Delete ISM policy
	resDelete, err := client.IsmDeletePolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.IsmGetPolicy("test").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
