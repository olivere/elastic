// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// GeoBoundingBoxQuery allows to filter hits based on a point location using
// a bounding box.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.0/query-dsl-geo-bounding-box-query.html
type GeoBoundingBoxQuery struct {
	name             string
	topLeft          interface{} // can be a GeoPoint, a GeoHash (string), or a lat/lon pair as float64
	topRight         interface{}
	bottomRight      interface{} // can be a GeoPoint, a GeoHash (string), or a lat/lon pair as float64
	bottomLeft       interface{}
	wkt              interface{}
	typ              string
	validationMethod string
	ignoreUnmapped   *bool
	queryName        string
}

// NewGeoBoundingBoxQuery creates and initializes a new GeoBoundingBoxQuery.
func NewGeoBoundingBoxQuery(name string) *GeoBoundingBoxQuery {
	return &GeoBoundingBoxQuery{
		name: name,
	}
}

// TopLeft position from longitude (left) and latitude (top).
func (q *GeoBoundingBoxQuery) TopLeft(top, left float64) *GeoBoundingBoxQuery {
	q.topLeft = []float64{left, top}
	return q
}

// TopLeftFromGeoPoint from a GeoPoint.
func (q *GeoBoundingBoxQuery) TopLeftFromGeoPoint(point *GeoPoint) *GeoBoundingBoxQuery {
	return q.TopLeft(point.Lat, point.Lon)
}

// TopLeftFromGeoHash from a Geo hash.
func (q *GeoBoundingBoxQuery) TopLeftFromGeoHash(topLeft string) *GeoBoundingBoxQuery {
	q.topLeft = topLeft
	return q
}

// BottomRight position from longitude (right) and latitude (bottom).
func (q *GeoBoundingBoxQuery) BottomRight(bottom, right float64) *GeoBoundingBoxQuery {
	q.bottomRight = []float64{right, bottom}
	return q
}

// BottomRightFromGeoPoint from a GeoPoint.
func (q *GeoBoundingBoxQuery) BottomRightFromGeoPoint(point *GeoPoint) *GeoBoundingBoxQuery {
	return q.BottomRight(point.Lat, point.Lon)
}

// BottomRightFromGeoHash from a Geo hash.
func (q *GeoBoundingBoxQuery) BottomRightFromGeoHash(bottomRight string) *GeoBoundingBoxQuery {
	q.bottomRight = bottomRight
	return q
}

// BottomLeft position from longitude (left) and latitude (bottom).
func (q *GeoBoundingBoxQuery) BottomLeft(bottom, left float64) *GeoBoundingBoxQuery {
	q.bottomLeft = []float64{bottom, left}
	return q
}

// BottomLeftFromGeoPoint from a GeoPoint.
func (q *GeoBoundingBoxQuery) BottomLeftFromGeoPoint(point *GeoPoint) *GeoBoundingBoxQuery {
	return q.BottomLeft(point.Lat, point.Lon)
}

// BottomLeftFromGeoHash from a Geo hash.
func (q *GeoBoundingBoxQuery) BottomLeftFromGeoHash(bottomLeft string) *GeoBoundingBoxQuery {
	q.bottomLeft = bottomLeft
	return q
}

// TopRight position from longitude (right) and latitude (top).
func (q *GeoBoundingBoxQuery) TopRight(top, right float64) *GeoBoundingBoxQuery {
	q.topRight = []float64{right, top}
	return q
}

// TopRightFromGeoPoint from a GeoPoint.
func (q *GeoBoundingBoxQuery) TopRightFromGeoPoint(point *GeoPoint) *GeoBoundingBoxQuery {
	return q.TopRight(point.Lat, point.Lon)
}

// TopRightFromGeoHash from a Geo hash.
func (q *GeoBoundingBoxQuery) TopRightFromGeoHash(topRight string) *GeoBoundingBoxQuery {
	q.topRight = topRight
	return q
}

// WKT initializes the bounding box from Well-Known Text (WKT),
// e.g. "BBOX (-74.1, -71.12, 40.73, 40.01)".
func (q *GeoBoundingBoxQuery) WKT(wkt interface{}) *GeoBoundingBoxQuery {
	q.wkt = wkt
	return q
}

// Type sets the type of executing the geo bounding box. It can be either
// memory or indexed. It defaults to memory.
func (q *GeoBoundingBoxQuery) Type(typ string) *GeoBoundingBoxQuery {
	q.typ = typ
	return q
}

// ValidationMethod accepts IGNORE_MALFORMED, COERCE, and STRICT (default).
// IGNORE_MALFORMED accepts geo points with invalid lat/lon.
// COERCE tries to infer the correct lat/lon.
func (q *GeoBoundingBoxQuery) ValidationMethod(method string) *GeoBoundingBoxQuery {
	q.validationMethod = method
	return q
}

// IgnoreUnmapped indicates whether to ignore unmapped fields (and run a
// MatchNoDocsQuery in place of this).
func (q *GeoBoundingBoxQuery) IgnoreUnmapped(ignoreUnmapped bool) *GeoBoundingBoxQuery {
	q.ignoreUnmapped = &ignoreUnmapped
	return q
}

// QueryName gives the query a name. It is used for caching.
func (q *GeoBoundingBoxQuery) QueryName(queryName string) *GeoBoundingBoxQuery {
	q.queryName = queryName
	return q
}

// Source returns JSON for the function score query.
func (q *GeoBoundingBoxQuery) Source() (interface{}, error) {
	// {
	//   "geo_bounding_box" : {
	//     ...
	//   }
	// }

	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source["geo_bounding_box"] = params

	box := make(map[string]interface{})
	if q.wkt != nil {
		box["wkt"] = q.wkt
	} else {
		if q.topLeft != nil {
			box["top_left"] = q.topLeft
		}
		if q.topRight != nil {
			box["top_right"] = q.topRight
		}
		if q.bottomLeft != nil {
			box["bottom_left"] = q.bottomLeft
		}
		if q.bottomRight != nil {
			box["bottom_right"] = q.bottomRight
		}
	}
	params[q.name] = box

	if q.typ != "" {
		params["type"] = q.typ
	}
	if q.validationMethod != "" {
		params["validation_method"] = q.validationMethod
	}
	if q.ignoreUnmapped != nil {
		params["ignore_unmapped"] = *q.ignoreUnmapped
	}
	if q.queryName != "" {
		params["_name"] = q.queryName
	}

	return source, nil
}
