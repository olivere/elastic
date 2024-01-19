package opensearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityCache(t *testing.T) {
	client := setupTestClient(t)
	var err error

	res, err := client.SecurityFlushCache().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, res)

}
