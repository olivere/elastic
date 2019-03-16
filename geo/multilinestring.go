package geo

import "encoding/json"

// MultiLineString is an object that must be an array of LineString coordinate arrays.
type MultiLineString struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// NewMultiLineString creates and initializes a new MultiLineString.
func NewMultiLineString(coordinates [][][]float64) *MultiLineString {
	return &MultiLineString{
		Type:        TypeMultiLineString,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the multi linestring data into a GeoJSON geometry.
func (m *MultiLineString) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *MultiLineString) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
