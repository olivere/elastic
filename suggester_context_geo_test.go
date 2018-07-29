// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestSuggesterGeoMapping(t *testing.T) {
	q := NewSuggesterGeoMapping("location").
		Precision("1km", "5m").
		Neighbors(true).
		FieldName("pin").
		DefaultLocations(GeoPointFromLatLon(0.0, 0.0))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"location":{"default":{"lat":0,"lon":0},"neighbors":true,"path":"pin","precision":["1km","5m"],"type":"geo"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterGeoQuery(t *testing.T) {
	q := NewSuggesterGeoQuery("location", GeoPointFromLatLon(11.5, 62.71)).Precision("1km").
		Neighbours("2km", "3km").Boost(2)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expectedOutcomes := []string{
		`{"location":{"context":{"lat":11.5,"lon":62.71},"precision":"1km","boost":2,"neighbours":["2km","3km"]}}`,
		`{"location":{"boost":2,"context":{"lat":11.5,"lon":62.71},"neighbours":["2km","3km"],"precision":"1km"}}`,
	}
	var match bool
	for _, expected := range expectedOutcomes {
		if got == expected {
			match = true
			break
		}
	}
	if !match {
		t.Errorf("expected any of %v\n,got:\n%s", expectedOutcomes, got)
	}
}

func TestSuggesterGeoIndex(t *testing.T) {
	in := NewSuggesterGeoIndex("location").Locations(GeoPointFromLatLon(11.5, 62.71))
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"location":{"lat":11.5,"lon":62.71}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSuggesterGeoIndexWithTwoValues(t *testing.T) {
	in := NewSuggesterGeoIndex("location").Locations(GeoPointFromLatLon(11.5, 62.71), GeoPointFromLatLon(31.5, 22.71))
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"location":[{"lat":11.5,"lon":62.71},{"lat":31.5,"lon":22.71}]}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
