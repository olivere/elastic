package geo

import "encoding/json"

// MultiPoint is an object that must be an array of positions.
type MultiPoint struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewMultiPoint creates and initializes a new MultiPoint.
func NewMultiPoint(coordinates [][]float64) *MultiPoint {
	return &MultiPoint{
		Type:        TypeMultiPoint,
		Coordinates: [][]float64(coordinates),
	}
}

// UnmarshalJSON decodes the multipoint data into a GeoJSON geometry.
func (m *MultiPoint) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *MultiPoint) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
