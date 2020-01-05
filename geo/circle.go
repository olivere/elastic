package geo

import (
	"encoding/json"
)

// Circle specified by a center point and radius with units,
// which default to meters.
type Circle struct {
	Type        string    `json:"type"`
	Radius      string    `json:"radius"`
	Coordinates []float64 `json:"coordinates"`
}

// NewCircle creates and initializes a new Circle.
func NewCircle(radius string, coordinates []float64) *Circle {
	return &Circle{
		Type:        TypeCircle,
		Radius:      radius,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the circle data into a GeoJSON geometry.
func (m *Circle) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Circle) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	m.Radius = raw.Radius
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
