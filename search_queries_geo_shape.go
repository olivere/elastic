// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// GeoShapeQuery allows to include hits that match a geoshape
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-geo-shape-query.html
type GeoShapeQuery struct {
	name      string
	geoJSON   interface{}
	relation  string
	queryName string
}

// NewGeoShapeQuery creates and initializes a new GeoShapeQuery.
func NewGeoShapeQuery(name string) *GeoShapeQuery {
	return &GeoShapeQuery{
		name: name,
	}
}

// SetPoint adds a point from latitude and longitude.
func (q *GeoShapeQuery) SetPoint(lat, lon float64) *GeoShapeQuery {
	var geoJSON struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}

	geoJSON.Type = "Point"
	geoJSON.Coordinates = []float64{lat, lon}

	q.geoJSON = geoJSON

	return q
}

// SetPolygon adds a polygon from a list of latitude and longitude.
func (q *GeoShapeQuery) SetPolygon(points [][][]float64) *GeoShapeQuery {
	var geoJSON struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}

	geoJSON.Type = "Polygon"
	geoJSON.Coordinates = points

	q.geoJSON = geoJSON

	return q
}

// SetRelation sets the geoJSON relation for the query
func (q *GeoShapeQuery) SetRelation(rel string) *GeoShapeQuery {
	q.relation = rel
	return q
}

func (q *GeoShapeQuery) QueryName(queryName string) *GeoShapeQuery {
	q.queryName = queryName
	return q
}

// Source returns JSON for the function score query.
func (q *GeoShapeQuery) Source() (interface{}, error) {
	// "geo_shape" : {
	//  	"person.location" : {
	//			// GeoJSON data...
	//     }
	// }
	source := make(map[string]interface{})
	params := make(map[string]interface{})

	if q.relation != "" {
		params[q.name] = map[string]interface{}{
			"shape":    q.geoJSON,
			"relation": q.relation,
		}
	} else {
		params[q.name] = map[string]interface{}{
			"shape": q.geoJSON,
		}
	}

	source["geo_shape"] = params

	if q.queryName != "" {
		params["_name"] = q.queryName
	}

	return source, nil
}
