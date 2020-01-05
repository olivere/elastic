// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"github.com/olivere/elastic/geo"
)

// GeoShapeQuery allows to find documents that have a shape that intersects with query shape.
//
// For more details, see:
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-geo-shape-query.html
type GeoShapeQuery struct {
	path         string
	relation     string
	queryName    string
	shape        *geo.Shape
	indexedShape *IndexedShape
}

// Pre-Indexed Shape the query also supports using a shape which has already been
// indexed in another index and/or type.
type IndexedShape struct {
	ID      string `json:"id,omitempty"`
	Index   string `json:"index,omitempty"`
	Type    string `json:"type,omitempty"`
	Path    string `json:"path,omitempty"`
	Routing string `json:"routing,omitempty"`
}

// NewGeoShapeQuery creates and initializes a new GeoShapeQuery.
func NewGeoShapeQuery(path string) *GeoShapeQuery {
	return &GeoShapeQuery{
		path: path,
	}
}

// IndexedShape adds a IndexedShape that represents a shape which has already been indexed in
// another index and/or index type.
func (q *GeoShapeQuery) IndexedShape(id, index, typ, path, routing string) *GeoShapeQuery {
	q.indexedShape = &IndexedShape{
		ID:      id,
		Index:   index,
		Type:    typ,
		Path:    path,
		Routing: routing,
	}
	return q
}

// AddPoint adds a Point that represents a single geographic coordinate.
func (q *GeoShapeQuery) AddPoint(coord []float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithPoint(coord))
	return q
}

// AddMultiPoint adds a MultiPoint that represents an array of unconnected
// but likely related points.
func (q *GeoShapeQuery) AddMultiPoint(coord [][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithMultiPoint(coord))
	return q
}

// AddLineString adds a LineString that represents an arbitray line given
// two or more points.
func (q *GeoShapeQuery) AddLineString(coord [][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithLineString(coord))
	return q
}

// AddMultiLineString adds a MultLineString that represents an array of separate
// linestrings.
func (q *GeoShapeQuery) AddMultiLineString(coord [][][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithMultiLineString(coord))
	return q
}

// AddPolygon adds a Polygon that represents a closed polygon whose first and last
// point must match, thus requiring n + 1 vertices to create an n- sided polygon and
// a minimum of 4 vertices.
func (q *GeoShapeQuery) AddPolygon(coord [][][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithPolygon(coord))
	return q
}

// AddMultiPolygon adds a MultiPolyogn that represents an array of separated polygons.
func (q *GeoShapeQuery) AddMultiPolygon(coord [][][][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithMultiPolygon(coord))
	return q
}

// AddGeometryCollection adds a GeometryCollection that represents a geo shape similar
// to the multi* shapes except that multiple types can coexist (e.g.: a Point and a LineString).
func (q *GeoShapeQuery) AddGeometryCollection(geometries []interface{}) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithGeometryCollection(geometries))
	return q
}

// AddEnvelope adds an Envelope that represents an bounding rectangle.
func (q *GeoShapeQuery) AddEnvelope(coord [][]float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithEnvelope(coord))
	return q
}

// AddCircle adds a Circle that represents a circle with radius in meters.
func (q *GeoShapeQuery) AddCircle(radius string, coord []float64) *GeoShapeQuery {
	q.shape = geo.NewShape(geo.WithCircle(radius, coord))
	return q
}

// Relation sets which spatial relation operators may ber used at search time.
func (q *GeoShapeQuery) Relation(relation string) *GeoShapeQuery {
	q.relation = relation
	return q
}

func (q *GeoShapeQuery) QueryName(queryName string) *GeoShapeQuery {
	q.queryName = queryName
	return q
}

// Source returns JSON for the function score query.
func (q *GeoShapeQuery) Source() (interface{}, error) {
	path := make(map[string]interface{})

	if q.indexedShape != nil {
		//	"geo_shape": {
		//		"location": {
		//			"indexed_shape": {
		//				"index": "shapes",
		//				"type": "_doc",
		//				"id": "deu",
		//				"path": "location"
		//			}
		//		}
		//	}
		path["indexed_shape"] = q.indexedShape
	} else if q.shape != nil {
		//	supports:
		//		- point
		//		- multipoint
		//		- linestring
		//		- multilinestring
		//		- polygon
		//		- multipolygon
		//		- envelope
		//		- circle
		//		- geometrycollection
		//	e..g:
		//	"geo_shape" : {
		//		"location": {
		//			"shape": {
		//				"type": "envelope"
		//				"coordinates": [[13.0, 53.0], [14.0, 52.0]]
		//			},
		//			"relation": "within"
		//		}
		//	}
		path["shape"] = q.shape
		if q.relation != "" {
			path["relation"] = q.relation
		}
		if q.queryName != "" {
			path["_name"] = q.queryName
		}
	}
	params := make(map[string]interface{})
	params[q.path] = path

	source := make(map[string]interface{})
	source["geo_shape"] = params

	return source, nil
}
