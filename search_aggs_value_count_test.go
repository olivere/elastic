package elastic

import (
	"encoding/json"
	"testing"
)

func TestValueCountAggregation(t *testing.T) {
	agg := NewValueCountAggregation().Field("grade")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"value_count":{"field":"grade"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
