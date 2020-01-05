package geo

import "encoding/json"

// Polygon is an object that must be an array of LinearRing coordinate arrays.
type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// NewPolygon creates and initializes a new Polygon.
func NewPolygon(coordinates [][][]float64) *Polygon {
	return &Polygon{
		Type:        TypePolygon,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the polygon data into a GeoJSON geometry.
func (m *Polygon) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Polygon) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
