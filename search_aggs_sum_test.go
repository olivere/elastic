package elastic

import (
	"encoding/json"
	"testing"
)

func TestSumAggregation(t *testing.T) {
	agg := NewSumAggregation().Field("price")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"sum":{"field":"price"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
