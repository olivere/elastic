package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeohashGridAggregation(t *testing.T) {
	agg := NewGeohashGridAggregation().Field("location").Precision(5)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeohashGridAggregationWithMetaData(t *testing.T) {
	agg := NewGeohashGridAggregation().Field("location").Precision(5)
	agg = agg.Meta(map[string]interface{}{"name": "Oliver"})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5},"meta":{"name":"Oliver"}}`

	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
