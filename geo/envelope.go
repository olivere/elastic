package geo

import "encoding/json"

// Envelope consists of coordinates for upper left and lower right points of the shape
// to represent a bounding rectangle.
type Envelope struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// NewEnvelope creates and initializes a new Rectangle.
func NewEnvelope(coordinates [][]float64) *Envelope {
	return &Envelope{
		Type:        TypeEnvelope,
		Coordinates: coordinates,
	}
}

// UnmarshalJSON decodes the envelope data into a GeoJSON geometry.
func (m *Envelope) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *Envelope) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	return json.Unmarshal(raw.Coordinates, &m.Coordinates)
}
