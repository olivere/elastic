package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityAudit(t *testing.T) {

	client := setupTestClient(t)
	var err error

	// Get current audit
	currentAudit, err := client.SecurityGetAudit().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, currentAudit)

	// Update audit
	currentAudit.Config.Audit.IgnoreUsers = []string{"test", "kibanaserver", "admin"}
	res, err := client.SecurityPutAudit().Body(currentAudit.Config).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, res)
	currentAudit, err = client.SecurityGetAudit().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"test", "kibanaserver", "admin"}, currentAudit.Config.Audit.IgnoreUsers)

}
