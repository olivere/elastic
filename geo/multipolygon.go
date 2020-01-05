package geo

import "encoding/json"

// MultiPolygon represents a GeoJSON object of multiple Polygons.
type MultiPolygon struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

// NewMultiPolygon creates and initializes a new MultiPolygon.
func NewMultiPolygon(coordinates [][][][]float64) *MultiPolygon {
	return &MultiPolygon{
		Type:        TypeMultiPolygon,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the multi polygon data into a GeoJSON geometry.
func (m *MultiPolygon) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *MultiPolygon) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
