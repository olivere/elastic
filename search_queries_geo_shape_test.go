// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoShapeQuery_Point_WithRelation(t *testing.T) {
	q := NewGeoShapeQuery("person.location")
	q = q.SetPoint(100.1, 0.1)
	q = q.SetRelation("intersects")
	q = q.QueryName("unit test")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"_name":"unit test","person.location":{"relation":"intersects","shape":{"type":"Point","coordinates":[100.1,0.1]}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoShapeQuery_Point_WithoutRelation(t *testing.T) {
	q := NewGeoShapeQuery("person.location")
	q = q.SetPoint(100.1, 0.1)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"person.location":{"shape":{"type":"Point","coordinates":[100.1,0.1]}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoShapeQuery_Polygon_WithRelation(t *testing.T) {
	q := NewGeoShapeQuery("person.location")
	q = q.SetPolygon(
		[][][]float64{
			[][]float64{
				[]float64{100.1, 0.1},
			},
		},
	)
	q = q.SetRelation("intersects")
	q = q.QueryName("unit test")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"_name":"unit test","person.location":{"relation":"intersects","shape":{"type":"Polygon","coordinates":[[[100.1,0.1]]]}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoShapeQuery_Polygon_WithoutRelation(t *testing.T) {
	q := NewGeoShapeQuery("person.location")
	q = q.SetPolygon(
		[][][]float64{
			[][]float64{
				[]float64{100.1, 0.1},
			},
		},
	)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"person.location":{"shape":{"type":"Polygon","coordinates":[[[100.1,0.1]]]}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
