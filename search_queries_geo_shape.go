// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// GeoShapeQuery allows to find documents that have a shape that intersects with the query shape
// a bounding box.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/6.2/query-dsl-geo-shape-query.html
type GeoShapeQuery struct {
	name      string
	index     string
	typ       string
	path      string
	id        interface{}
	coord     interface{}
	relation  string
	queryName string
}

// NewGeoShapeQuery creates and initializes a new GeoShapeQuery.
func NewGeoShapeQuery(name string) *GeoShapeQuery {
	return &GeoShapeQuery{
		name: name,
	}
}

func (q *GeoShapeQuery) IndexedShape(index, typ, path string, id interface{}) *GeoShapeQuery {
	q.index = index
	q.typ = typ
	q.path = path
	q.id = id
	return q
}

// Type sets the type of coordinates that you will provide, it can be
// point, polygon, envelope, etc.
func (q *GeoShapeQuery) Type(typ string) *GeoShapeQuery {
	q.typ = typ
	return q
}

// Coordinates depends of your type, can be a point ([]float64), can be
// a linestring ([][]float64) a polygon ([][][]float64), etc.
func (q *GeoShapeQuery) Coordinates(coord interface{}) *GeoShapeQuery {
	q.coord = coord
	return q
}

// Relation sets which spatial relation operators may be used at search time.
// Operators available are intersects, disjoint, within and contains.
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
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	shape := make(map[string]interface{})
	source["geo_shape"] = params

	if q.index != "" {
		// Pre-Indexed Shape
		// {
		// 	 "geo_shape" : {
		//     "location": {
		//       "indexed_shape": {
		//         "index": "shapes",
		//         "type": "_doc",
		//         "id": "deu",
		//         "path": "location"
		//       }
		//     }
		//   }
		// }
		indexedShape := make(map[string]interface{})
		indexedShape["index"] = q.index
		indexedShape["type"] = q.typ
		indexedShape["path"] = q.path
		indexedShape["id"] = q.id
		shape["indexed_shape"] = indexedShape
	} else {
		// Inline Shape
		// {
		// 	 "geo_shape" : {
		//     "location": {
		//       "shape": {
		//         "type": "envelope",
		//         "coordinates": [[13.0, 53.0], [14.0, 52.0]]
		//       },
		//       "relation": "within"
		//     }
		//   }
		// }
		inlineShape := make(map[string]interface{})
		inlineShape["type"] = q.typ
		inlineShape["coordinates"] = q.coord
		shape["shape"] = inlineShape
	}

	params[q.name] = shape

	if q.relation != "" {
		params["relation"] = q.relation
	}
	if q.queryName != "" {
		params["_name"] = q.queryName
	}

	return source, nil
}
