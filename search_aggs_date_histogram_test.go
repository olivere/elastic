package elastic

import (
	"encoding/json"
	"testing"
)

func TestDateHistogramAggregation(t *testing.T) {
	agg := NewDateHistogramAggregation().Field("date").Interval("month")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"field":"date","interval":"month"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
