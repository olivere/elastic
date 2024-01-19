package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestSecurityConfig(t *testing.T) {
	client := setupTestClient(t)
	var err error

	// Get current config
	currentConfig, err := client.SecurityGetConfig().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, currentConfig)

	// Update config
	currentConfig.Config.Dynamic.DoNotFailOnForbidden = ptr.To[bool](true)
	res, err := client.SecurityPutConfig().Body(currentConfig.Config).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, res)
	currentConfig, err = client.SecurityGetConfig().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, *currentConfig.Config.Dynamic.DoNotFailOnForbidden)

}
