package geo

import "encoding/json"

// LineString is an object that must be an array of two or more positions.
type LineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewLineString creates and initializes a new LineString.
func NewLineString(coordinates [][]float64) *LineString {
	return &LineString{
		Type:        TypeLineString,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the linestring data into a GeoJSON geometry.
func (m *LineString) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *LineString) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
