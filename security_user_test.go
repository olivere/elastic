package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestSecurityUser(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expecedUser := &SecurityPutUser{
		SecurityUserBase: SecurityUserBase{
			BackendRoles:  []string{"admin"},
			SecurityRoles: []string{"all_access"},
			Attributes:    map[string]string{},
		},
		Password: ptr.To[string]("myverystrongpassword"),
	}

	// Create user
	resPut, err := client.SecurityPutUser("test").Body(expecedUser).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get user
	resGet, err := client.SecurityGetUser("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.Equal(t, expecedUser.SecurityUserBase, (*resGet)["test"].SecurityUserBase)

	// Update user
	expecedUser.Description = ptr.To[string]("test")
	_, err = client.SecurityPutUser("test").Body(expecedUser).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SecurityGetUser("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expecedUser.SecurityUserBase, (*resGet)["test"].SecurityUserBase)

	// Delete user
	resDelete, err := client.SecurityDeleteUser("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetUser("test").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
