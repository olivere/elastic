package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRoleMapping(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expectedRoleMapping := &SecurityPutRoleMapping{
		BackendRoles:    []string{"admin"},
		Users:           []string{"admin"},
		AndBackendRoles: []string{},
		Hosts:           []string{},
	}

	// Put role mapping
	resPut, err := client.SecurityPutRoleMapping("kibana_user").Body(expectedRoleMapping).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get role mapping
	resGet, err := client.SecurityGetRoleMapping("kibana_user").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.Equal(t, *expectedRoleMapping, (*resGet)["kibana_user"].SecurityPutRoleMapping)

	// Update role mapping
	expectedRoleMapping.AndBackendRoles = []string{"kibanaserver"}
	_, err = client.SecurityPutRoleMapping("kibana_user").Body(expectedRoleMapping).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SecurityGetRoleMapping("kibana_user").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, *expectedRoleMapping, (*resGet)["kibana_user"].SecurityPutRoleMapping)

	// Delete role mapping
	resDelete, err := client.SecurityDeleteRoleMapping("kibana_user").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetRoleMapping("kibana_user").Do(context.Background())
	assert.True(t, IsNotFound(err))
}
