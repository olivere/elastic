// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestGeoShapeQueryWithPreIndexedShape(t *testing.T) {
	q := NewGeoShapeQuery("pin.location")
	q = q.IndexedShape("shapes", "_doc", "location", "deu")
	q = q.Relation("contains")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"pin.location":{"indexed_shape":{"id":"deu","index":"shapes","path":"location","type":"_doc"}},"relation":"contains"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoShapeQueryWithPoint(t *testing.T) {
	q := NewGeoShapeQuery("pin.location")
	q = q.Type("point")
	q = q.Coordinates([]float64{13.0, 53.0})
	q = q.Relation("contains")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"pin.location":{"shape":{"coordinates":[13,53],"type":"point"}},"relation":"contains"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoShapeQueryWithEnvelope(t *testing.T) {
	q := NewGeoShapeQuery("pin.location")
	q = q.Type("envelope")
	q = q.Coordinates([][]float64{
		{13.0, 53.0},
		{14.0, 52.0},
	})
	q = q.Relation("within")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"geo_shape":{"pin.location":{"shape":{"coordinates":[[13,53],[14,52]],"type":"envelope"}},"relation":"within"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
