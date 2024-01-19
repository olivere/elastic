package opensearch

// Need auth be certificate admin ...
/*
import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSecurityDN(t *testing.T) {
	client := setupTestClient(t)
	var err error

	expectedDN := &SecurityDistinguishedName{
		NodesDN: []string{".*"},
	}

	// Create dn
	resPut, err := client.SecurityPutDistinguishedName("opensearch").Body(expectedDN).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get DN
	resGet, err := client.SecurityGetDistinguishedName("opensearch").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.Equal(t, *expectedDN, (*resGet)["opensearch"])

	// Delete dn
	resDelete, err := client.SecurityDeleteDistinguishedName("opensearch").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SecurityGetDistinguishedName("opensearch").Do(context.Background())
	assert.True(t, IsNotFound(err))

}

*/
