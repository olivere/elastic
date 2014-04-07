package elastic

import (
	"encoding/json"
	"testing"
)

func TestStatsAggregation(t *testing.T) {
	agg := NewStatsAggregation().Field("grade")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"stats":{"field":"grade"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
