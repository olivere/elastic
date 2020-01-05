// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"

	"github.com/olivere/elastic/geo"
	"github.com/stretchr/testify/require"
)

func TestGeoShapeQueryIndexedShape(t *testing.T) {
	var (
		id       = "deu"
		index    = "shapes"
		itype    = "_doc"
		path     = "location"
		routing  = ""
		expected = []byte(`{"geo_shape":{"location":{"indexed_shape":{"index":"shapes","type":"_doc","id":"deu","path":"location"}}}}`)
	)
	assert := require.New(t)

	q := NewGeoShapeQuery("location").IndexedShape(id, index, itype, path, routing)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryPoint(t *testing.T) {
	var (
		coordinates = []float64{13.400544, 52.530286}
		relation    = "within"
		expected    = []byte(`{"geo_shape":{"location":{"shape":{"type":"point","coordinates":[13.400544,52.530286]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location")
	q = q.AddPoint(coordinates)
	q = q.Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryLineString(t *testing.T) {
	var (
		coordinates = [][]float64{{-77.03653, 38.897676}, {-77.009051, 38.889939}}
		relation    = "within"
		expected    = []byte(`{"geo_shape":{"location":{"shape":{"type":"linestring","coordinates":[[-77.03653, 38.897676],[-77.009051, 38.889939]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddLineString(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryPolygon(t *testing.T) {
	var (
		coordinates = [][][]float64{
			[][]float64{
				{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
			},
		}
		relation = "within"
		expected = []byte(`{"geo_shape":{"location":{"shape":{"type":"polygon","coordinates":[[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddPolygon(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryMultiPoint(t *testing.T) {
	var (
		coordinates = [][]float64{{102.0, 2.0}, {103.0, 2.0}}
		relation    = "within"
		expected    = []byte(`{"geo_shape":{"location":{"shape":{"type":"multipoint","coordinates":[[102.0, 2.0],[103.0, 2.0]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddMultiPoint(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryMultiLineString(t *testing.T) {
	var (
		coordinates = [][][]float64{
			[][]float64{
				{100.0, 0.0}, {101.0, 1.0},
			},
			[][]float64{
				{102.0, 2.0}, {103.0, 3.0},
			},
		}
		relation = "within"
		expected = []byte(`{"geo_shape":{"location":{"shape":{"type":"multilinestring","coordinates":[[[100.0,0.0],[101.0,1.0]],[[102.0,2.0],[103.0,3.0]]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddMultiLineString(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryMultiPolygon(t *testing.T) {
	var (
		coordinates = [][][][]float64{
			{
				[][]float64{
					{102.0, 2.0}, {103.0, 2.0}, {103.0, 3.0}, {102.0, 3.0}, {102.0, 2.0},
				},
			},
			{
				[][]float64{
					{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
				},
				[][]float64{
					{100.2, 0.2}, {100.8, 0.2}, {100.8, 0.8}, {100.2, 0.8}, {100.2, 0.2},
				},
			},
		}
		relation = "within"
		expected = []byte(`{"geo_shape":{"location":{"shape":{"type":"multipolygon","coordinates":[[[[102.0,2.0],[103.0,2.0],[103.0,3.0],[102.0,3.0],[102.0,2.0]]],[[[100.0,0.0],[101.0,0.0],[101.0,1.0],[100.0,1.0],[100.0,0.0]],[[100.2,0.2],[100.8,0.2],[100.8,0.8],[100.2,0.8],[100.2,0.2]]]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddMultiPolygon(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryGeometryCollection(t *testing.T) {
	var (
		geometries = []interface{}{
			geo.NewPoint([]float64{100.0, 0.0}),
			geo.NewLineString([][]float64{{101.0, 0.0}, {102.0, 1.0}}),
		}
		relation = "within"
		expected = []byte(`{"geo_shape":{"location":{"shape":{"type":"geometrycollection","geometries":[{"type":"point","coordinates":[100.0,0.0]},{"type":"linestring","coordinates":[[101.0,0.0],[102.0,1.0]]}]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddGeometryCollection(geometries).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryEnvelope(t *testing.T) {
	var (
		coordinates = [][]float64{{100.0, 1.0}, {101.0, 0.0}}
		relation    = "within"
		expected    = []byte(`{"geo_shape":{"location":{"shape":{"type":"envelope","coordinates":[[100.0,1.0],[101.0,0.0]]},"relation":"within"}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddEnvelope(coordinates).Relation(relation)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}

func TestGeoShapeQueryCircle(t *testing.T) {
	var (
		coordinates = []float64{-109.874838, 44.439550}
		radius      = "25m"
		expected    = []byte(`{"geo_shape":{"location":{"shape":{"type":"circle","radius":"25m","coordinates":[-109.874838,44.439550]}}}}`)
	)
	assert := require.New(t)
	q := NewGeoShapeQuery("location").AddCircle(radius, coordinates)

	src, err := q.Source()
	assert.NoError(err)

	actual, err := json.Marshal(src)
	assert.NoError(err)
	assert.JSONEq(string(expected), string(actual))
}
