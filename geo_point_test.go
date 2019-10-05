// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
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

func TestGeoPointMarshalJSON(t *testing.T) {
	pt := GeoPoint{Lat: 40, Lon: -70}

	data, err := json.Marshal(pt)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"lat":40,"lon":-70}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestGeoPointIndexAndSearch(t *testing.T) {
	client := setupTestClient(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	// Create index
	mapping := `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"doc":{
				"properties":{
					"name":{
						"type":"keyword"
					},
					"location":{
						"type":"geo_point"
					}
				}
			}
		}
	}
`
	createIndex, err := client.CreateIndex(testIndexName).Body(mapping).IncludeTypeName(true).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Add document
	type City struct {
		Name     string    `json:"name"`
		Location *GeoPoint `json:"location"`
	}
	munich := &City{
		Name:     "München",
		Location: GeoPointFromLatLon(48.137154, 11.576124),
	}
	_, err = client.Index().Index(testIndexName).Type("doc").Id("1").BodyJson(&munich).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Flush
	_, err = client.Flush().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Get document
	q := NewGeoDistanceQuery("location")
	q = q.GeoPoint(GeoPointFromLatLon(48, 11))
	q = q.Distance("50km")
	res, err := client.
		Search(testIndexName).
		Type("doc").
		Query(q).
		Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := int64(1), res.TotalHits(); want != have {
		t.Fatalf("TotalHits: want %d, have %d", want, have)
	}
	var doc City
	if err := json.Unmarshal(*res.Hits.Hits[0].Source, &doc); err != nil {
		t.Fatal(err)
	}
	if want, have := munich.Name, doc.Name; want != have {
		t.Fatalf("Name: want %q, have %q", want, have)
	}
	if want, have := munich.Location.Lat, doc.Location.Lat; want != have {
		t.Fatalf("Lat: want %v, have %v", want, have)
	}
	if want, have := munich.Location.Lon, doc.Location.Lon; want != have {
		t.Fatalf("Lon: want %v, have %v", want, have)
	}
}
