package opensearch

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestSmPolicy(t *testing.T) {
	client := setupTestClient(t)
	var err error

	logrus.SetLevel(logrus.TraceLevel)

	expecedSmPolicy := &SmPutPolicy{
		Description: ptr.To[string]("Daily snapshot policy"),
		Creation: SmPolicyCreation{
			Schedule: map[string]any{
				"cron": map[string]any{
					"expression": "0 8 * * *",
					"timezone":   "UTC",
				},
			},
			TimeLimit: ptr.To[string]("1h"),
		},
		Deletion: &SmPolicyDeletion{
			Schedule: map[string]any{
				"cron": map[string]any{
					"expression": "0 1 * * *",
					"timezone":   "America/Los_Angeles",
				},
			},
			Condition: &SmPolicyDeleteCondition{
				MaxAge:   ptr.To[string]("7d"),
				MaxCount: ptr.To[int64](21),
				MinCount: ptr.To[int64](7),
			},
			TimeLimit: ptr.To[string]("1h"),
		},
		SnapshotConfig: SmPolicySnapshotConfig{
			DateFormat:         ptr.To[string]("yyyy-MM-dd-HH:mm"),
			Timezone:           ptr.To[string]("America/Los_Angeles"),
			Indices:            ptr.To[string]("*"),
			Repository:         "s3-repo",
			IgnoreUnavailable:  ptr.To[bool](true),
			IncludeGlobalState: ptr.To[bool](false),
			Partial:            ptr.To[bool](true),
			Metadata: map[string]any{
				"any_key": "any_value",
			},
		},
		Notification: &SmPolicyNotification{
			Channel: SmPolicyNotificationChannel{
				ID: "NC3OpoEBzEoHMX183R3f",
			},
			Conditions: &SmPolicyNotificationCondition{
				Creation:          ptr.To[bool](true),
				Deletion:          ptr.To[bool](false),
				Failure:           ptr.To[bool](false),
				TimeLimitExceeded: ptr.To[bool](false),
			},
		},
	}

	// Create SM policy
	resPut, err := client.SmPostPolicy("test").Body(expecedSmPolicy).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPut)

	// Get SM policy
	resGet, err := client.SmGetPolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.NotNil(t, resGet.Policy)

	// Update SM policy
	expecedSmPolicy.Description = ptr.To[string]("test")
	_, err = client.SmPutPolicy("test").Body(expecedSmPolicy).SequenceNumber(resGet.SequenceNumber).PrimaryTerm(resGet.PrimaryTerm).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.SmGetPolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", *resGet.Policy.Description)

	// Delete SM policy
	resDelete, err := client.SmDeletePolicy("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.SmGetPolicy("test").Do(context.Background())
	assert.True(t, IsNotFound(err))

}
