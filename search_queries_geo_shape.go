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
	name         string
	shape        Shape
	strategy     string
	relation     string
	indexedShape struct {
		index   string
		typ     string
		path    string
		id      string
		routing string
	}
	ignoreUnmapped bool
	queryName      string
}

// NewGeoShapeQuery creates and initializes a new GeoShapeQuery.
func NewGeoShapeQuery(name string) *GeoShapeQuery {
	return &GeoShapeQuery{
		name: name,
	}
}

func (q *GeoShapeQuery) Shape(s Shape) *GeoShapeQuery {
	q.shape = s
	return q
}

func (q *GeoShapeQuery) Strategy(strategy string) *GeoShapeQuery {
	q.strategy = strategy
	return q
}

func (q *GeoShapeQuery) Relation(relation string) *GeoShapeQuery {
	q.relation = relation
	return q
}

func (q *GeoShapeQuery) IndexedShape(index, typ, path, id, routing string) *GeoShapeQuery {
	q.indexedShape.index = index
	q.indexedShape.typ = typ
	q.indexedShape.path = path
	q.indexedShape.id = id
	q.indexedShape.routing = routing
	return q
}

func (q *GeoShapeQuery) IndexedShapeIndex(index string) *GeoShapeQuery {
	q.indexedShape.index = index
	return q
}

func (q *GeoShapeQuery) IndexedShapeType(typ string) *GeoShapeQuery {
	q.indexedShape.typ = typ
	return q
}

func (q *GeoShapeQuery) IndexedShapePath(path string) *GeoShapeQuery {
	q.indexedShape.path = path
	return q
}

func (q *GeoShapeQuery) IndexedShapeID(id string) *GeoShapeQuery {
	q.indexedShape.id = id
	return q
}

func (q *GeoShapeQuery) IndexedShapeRouting(routing string) *GeoShapeQuery {
	q.indexedShape.routing = routing
	return q
}

func (q *GeoShapeQuery) IgnoreUnmapped(ignoreUnmapped bool) *GeoShapeQuery {
	q.ignoreUnmapped = ignoreUnmapped
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
	source["geo_shape"] = params

	shape := make(map[string]interface{})
	params[q.name] = shape
	params["ignore_unmapped"] = q.ignoreUnmapped

	if q.strategy != "" {
		shape["strategy"] = q.strategy
	}

	if q.relation != "" {
		shape["relation"] = q.relation
	}

	if q.indexedShape.id == "" {
		var err error

		shape["shape"], err = q.shape.Source()
		if err != nil {
			return nil, err
		}
	} else {
		indexedShape := make(map[string]interface{})
		shape["indexed_shape"] = indexedShape
		indexedShape["id"] = q.indexedShape.id
		indexedShape["type"] = q.indexedShape.typ
		if q.indexedShape.index != "" {
			indexedShape["index"] = q.indexedShape.index
		}
		if q.indexedShape.path != "" {
			indexedShape["path"] = q.indexedShape.path
		}
		if q.indexedShape.routing != "" {
			indexedShape["routing"] = q.indexedShape.routing
		}
	}

	if q.queryName != "" {
		params["_name"] = q.queryName
	}

	return source, nil
}

type Shape struct {
}

func (s *Shape) Source() (interface{}, error) {
	return nil, nil
}
