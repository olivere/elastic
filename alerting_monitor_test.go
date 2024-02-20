package opensearch

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestAlertingMonitor(t *testing.T) {
	client := setupTestClient(t)
	var err error

	logrus.SetLevel(logrus.TraceLevel)

	expecedAlertingMonitor := &AlertingMonitor{
		Type:        "monitor",
		Name:        "test",
		MonitorType: "query_level_monitor",
		Enabled:     ptr.To[bool](true),
		Schedule: map[string]any{
			"period": map[string]any{
				"interval": 1,
				"unit":     "MINUTES",
			},
		},
		Inputs: []map[string]any{
			{
				"search": map[string]any{
					"indices": []string{"*"},
					"query": map[string]any{
						"query": map[string]any{
							"match_all": map[string]any{},
						},
					},
				},
			},
		},
		Triggers: []map[string]any{
			{
				"name":     "test-trigger",
				"severity": "1",
				"condition": map[string]any{
					"script": map[string]any{
						"source": "ctx.results[0].hits.total.value > 0",
						"lang":   "painless",
					},
				},
				"actions": []map[string]any{
					{
						"name":           "test-action",
						"destination_id": "ld7912sBlQ5JUWWFThoW",
						"message_template": map[string]any{
							"source": "This is my message body.",
						},
						"throttle_enabled": true,
						"throttle": map[string]any{
							"value": 27,
							"unit":  "MINUTES",
						},
						"subject_template": map[string]any{
							"source": "TheSubject",
						},
					},
				},
			},
		},
	}

	// Create monitor
	resPost, err := client.AlertingPostMonitor().Body(expecedAlertingMonitor).Pretty(true).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resPost)

	// Get monitor
	resGet, err := client.AlertingGetMonitor(resPost.Id).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resGet)
	assert.NotNil(t, resGet.Monitor)

	// Search monitor
	resSearch, err := client.AlertingSearchMonitor().SearchByName("test").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, resSearch)
	assert.Equal(t, "test", resSearch[0].Monitor.Name)
	assert.NotEmpty(t, resSearch[0].Id)

	// Update monitor
	expecedAlertingMonitor.Name = "test2"
	_, err = client.AlertingPutMonitor(resPost.Id).Body(expecedAlertingMonitor).SequenceNumber(resGet.SequenceNumber).PrimaryTerm(resGet.PrimaryTerm).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resGet, err = client.AlertingGetMonitor(resPost.Id).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test2", resGet.Monitor.Name)

	// Delete monitor
	resDelete, err := client.AlertingDeleteMonitor(resPost.Id).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resDelete)
	_, err = client.AlertingGetMonitor(resPost.Id).Do(context.Background())
	assert.True(t, IsNotFound(err))

}
