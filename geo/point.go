package geo

import "encoding/json"

// Point is an object that must be a single position.
type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// NewPoint creates and initializes a new Point.
func NewPoint(coordinates []float64) *Point {
	return &Point{
		Type:        TypePoint,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the point data into a GeoJSON geometry.
func (m *Point) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Point) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
