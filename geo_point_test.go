package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoPointSource(t *testing.T) {
	pt := GeoPoint{Lat: 40, Lon: -70}

	data, err := json.Marshal(pt.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"lat":40,"lon":-70}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
