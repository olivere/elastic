package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityUser(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expecedUser := &SecurityPutUser{
		SecurityUser: SecurityUser{
			BackendRoles:  []string{"admin"},
			SecurityRoles: []string{"all_access"},
			Attributes:    map[string]string{},
		},
		Password: "myverystrongpassword",
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
	assert.Equal(t, expecedUser.SecurityUser, (*resGet)["test"])

	// Update user
	expecedUser.Description = "test"
	_, err = client.SecurityPutUser("test").Body(expecedUser).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SecurityGetUser("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expecedUser.SecurityUser, (*resGet)["test"])

	// Delete user
	resDelete, err := client.SecurityDeleteUser("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)

}
