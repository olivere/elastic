package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoHashGridAggregation(t *testing.T) {
	agg := NewGeoHashGridAggregation().Field("location").Precision(5)
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoHashGridAggregationWithMetaData(t *testing.T) {
	agg := NewGeoHashGridAggregation().Field("location").Precision(5)
	agg = agg.Meta(map[string]interface{}{"name": "Oliver"})
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5},"meta":{"name":"Oliver"}}`

	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoHashGridAggregationWithSize(t *testing.T) {
	agg := NewGeoHashGridAggregation().Field("location").Precision(5).Size(5)
	agg = agg.Meta(map[string]interface{}{"name": "Oliver"})
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5,"size":5},"meta":{"name":"Oliver"}}`

	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoHashGridAggregationWithShardSize(t *testing.T) {
	agg := NewGeoHashGridAggregation().Field("location").Precision(5).ShardSize(5)
	agg = agg.Meta(map[string]interface{}{"name": "Oliver"})
	data, err := json.Marshal(agg.Source())
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geohash_grid":{"field":"location","precision":5,"shard_size":5},"meta":{"name":"Oliver"}}`

	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
