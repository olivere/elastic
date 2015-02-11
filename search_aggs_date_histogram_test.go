package elastic

import (
	"encoding/json"
	"testing"
)

func TestDateHistogramAggregation(t *testing.T) {
	agg := NewDateHistogramAggregation().Field("date").Interval("month").Format("YYYY-MM")
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"field":"date","format":"YYYY-MM","interval":"month"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
