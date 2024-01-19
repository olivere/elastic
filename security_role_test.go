package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRole(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expectedRole := &SecurityRole{
		ClusterPermissions: []string{"*"},
		IndexPermissions: []SecurityIndexPermissions{
			{
				IndexPatterns:  []string{"*"},
				AllowedActions: []string{"*"},
				MaskedFields:   []string{},
				FieldLevelSecurity:            []string{},
			},
		},
		TenantPermissions: []SecurityTenantPermissions{
			{
				TenantPatterns: []string{"*"},
				AllowedAction:  []string{"*"},
			},
		},
	}

	// Create role
	resPut, err := client.SecurityPutRole("superuser").Body(expectedRole).Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, resPut)

	// Get role
	resGet, err := client.SecurityGetRole("superuser").Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, resGet)
	assert.Equal(t, *expectedRole, (*resGet)["superuser"])

	// Update role
	expectedRole.ClusterPermissions = []string{"cluster:admin/opendistro/alerting/alerts/get"}
	_, err = client.SecurityPutRole("superuser").Body(expectedRole).Do(context.Background())
	assert.NoError(t, err)
	resGet, err = client.SecurityGetRole("superuser").Do(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, *expectedRole, (*resGet)["superuser"])

	// Delete role
	resDelete, err := client.SecurityDeleteRole("superuser").Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetRole("superuser").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
