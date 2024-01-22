package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestSecurityTenant(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expectedTenant := &SecurityPutTenant{
		Description: ptr.To[string]("test"),
	}

	// Create tenant
	resPut, err := client.SecurityPutTenant("test").Body(expectedTenant).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get tenant
	resGet, err := client.SecurityGetTenant("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.Equal(t, *expectedTenant, (*resGet)["test"].SecurityPutTenant)

	// Update tenant
	expectedTenant.Description = ptr.To[string]("this is a test")
	_, err = client.SecurityPutTenant("test").Body(expectedTenant).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SecurityGetTenant("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, *expectedTenant, (*resGet)["test"].SecurityPutTenant)

	// Delete tenant
	resDelete, err := client.SecurityDeleteTenant("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetTenant("test").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
