package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityActionGroup(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expectedActionGroup := &SecurityActionGroup{
		AllowedActions: []string{
			"cluster_all",
		},
	}

	// Create action group
	resPut, err := client.SecurityPutActionGroup("test").Body(expectedActionGroup).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get action group
	resGet, err := client.SecurityGetActionGroup("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.Equal(t, *expectedActionGroup, (*resGet)["test"])

	// Update action group
	expectedActionGroup.Description = "test"
	_, err = client.SecurityPutActionGroup("test").Body(expectedActionGroup).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SecurityGetActionGroup("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, *expectedActionGroup, (*resGet)["test"])

	// Delete action group
	resDelete, err := client.SecurityDeleteActionGroup("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetActionGroup("test").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
