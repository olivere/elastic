package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoTileGridAggregation(t *testing.T) {
	bounds := BoundingBox{
		TopLeft: GeoPoint{
			Lat: 55.145984,
			Lon: 82.75195,
		},
		BottomRight: GeoPoint{
			Lat: 54.830199,
			Lon: 83.143839,
		},
	}
	agg := NewGeoTileGridAggregation().
		Field("location").
		Precision(12).
		Size(3).
		ShardSize(5).
		Bounds(bounds).
		Meta(map[string]interface{}{"city": "Novosibirsk"})

	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshalling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geotile_grid":{"bounds":{"top_left":{"lat":55.145984,"lon":82.75195},"bottom_right":{"lat":54.830199,"lon":83.143839}},"field":"location","precision":12,"shard_size":5,"size":3},"meta":{"city":"Novosibirsk"}}`

	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
