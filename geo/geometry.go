package geo

import (
	"encoding/json"
)

// The geometry types supported by Elasticsearch.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/geo-shape.html#spatial-strategy
const (
	TypePoint              = "point"
	TypeMultiPoint         = "multipoint"
	TypeLineString         = "linestring"
	TypeMultiLineString    = "multilinestring"
	TypePolygon            = "polygon"
	TypeMultiPolygon       = "multipolygon"
	TypeGeometryCollection = "geometrycollection"
	TypeEnvelope           = "envelope"
	TypeCircle             = "circle"
)

// rawGeometry holds generic data used to unmarshal GeoJSON information.
type rawGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Geometries  []rawGeometry   `json:"geometries"`
	Radius      string          `json:"radius,omitempty"`
}
