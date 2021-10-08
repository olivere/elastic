// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoBoundingBoxQuery(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.TopLeft(40.73, -74.1)
	q = q.BottomRight(40.01, -71.12)
	q = q.Type("memory")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"pin.location":{"bottom_right":[-71.12,40.01],"top_left":[-74.1,40.73]},"type":"memory"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoBoundingBoxQueryWithGeoPoint(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.TopLeftFromGeoPoint(GeoPointFromLatLon(40.73, -74.1))
	q = q.BottomRightFromGeoPoint(GeoPointFromLatLon(40.01, -71.12))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"pin.location":{"bottom_right":[-71.12,40.01],"top_left":[-74.1,40.73]}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoBoundingBoxQueryWithGeoHash(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.TopLeftFromGeoHash("dr5r9ydj2y73")
	q = q.BottomRightFromGeoHash("drj7teegpus6")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"pin.location":{"bottom_right":"drj7teegpus6","top_left":"dr5r9ydj2y73"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoBoundingBoxQueryWithWKT(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.WKT("BBOX (-74.1, -71.12, 40.73, 40.01)")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"pin.location":{"wkt":"BBOX (-74.1, -71.12, 40.73, 40.01)"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoBoundingBoxQueryWithMixed(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.TopLeftFromGeoPoint(GeoPointFromLatLon(40.73, -74.1))
	q = q.BottomRightFromGeoHash("drj7teegpus6")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"pin.location":{"bottom_right":"drj7teegpus6","top_left":[-74.1,40.73]}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoBoundingBoxQueryWithParameters(t *testing.T) {
	q := NewGeoBoundingBoxQuery("pin.location")
	q = q.TopLeftFromGeoHash("dr5r9ydj2y73")
	q = q.BottomRightFromGeoHash("drj7teegpus6")
	q = q.ValidationMethod("IGNORE_MALFORMED")
	q = q.IgnoreUnmapped((true))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_bounding_box":{"ignore_unmapped":true,"pin.location":{"bottom_right":"drj7teegpus6","top_left":"dr5r9ydj2y73"},"validation_method":"IGNORE_MALFORMED"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
