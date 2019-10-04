// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestDistanceFeatureQuery(t *testing.T) {
	q := NewDistanceFeatureQuery("coordinates", map[string]float64{"lat": 11.111, "lon": 12.121}, "10km")
	q.Boost(3.4)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"distance_feature":{"boost":3.4,"field":"coordinates","origin":{"lat":11.111,"lon":12.121},"pivot":"10km"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
