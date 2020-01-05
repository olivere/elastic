package geo

import (
	"encoding/json"
	"fmt"
)

// Shape is a `generic` object used to unmarshal geometric objects.
type Shape struct {
	Type               string              `json:"type"`
	Point              *Point              `json:"-"`
	MultiPoint         *MultiPoint         `json:"-"`
	LineString         *LineString         `json:"-"`
	MultiLineString    *MultiLineString    `json:"-"`
	Polygon            *Polygon            `json:"-"`
	MultiPolygon       *MultiPolygon       `json:"-"`
	GeometryCollection *GeometryCollection `json:"-"`
	Envelope           *Envelope           `json:"-"`
	Circle             *Circle             `json:"-"`
}

// NewShape creates and initializes a new shape based on given option.
func NewShape(opts ...Option) *Shape {
	shape := &Shape{}
	for _, opt := range opts {
		opt(shape)
	}
	return shape
}

// IsPoint returns true whether valid Point object.
func (m *Shape) IsPoint() bool {
	return m.Type == TypePoint && m.Point != nil
}

// IsMultiPoint returns true whether valid MultiPoint object.
func (m *Shape) IsMultiPoint() bool {
	return m.Type == TypeMultiPoint && m.MultiPoint != nil
}

// IsLineString returns true whether valid LineString object.
func (m *Shape) IsLineString() bool {
	return m.Type == TypeLineString && m.LineString != nil
}

// IsMultiLineString returns true whether valid MultiLineString object.
func (m *Shape) IsMultiLineString() bool {
	return m.Type == TypeMultiLineString && m.MultiLineString != nil
}

// IsPolygon returns true whether valid Polygon object.
func (m *Shape) IsPolygon() bool {
	return m.Type == TypePolygon && m.Polygon != nil
}

// IsMultiPolygon returns true whether valid MultiPolygon.
func (m *Shape) IsMultiPolygon() bool {
	return m.Type == TypeMultiPolygon && m.MultiPolygon != nil
}

// IsGeometryCollection returns true whether valid GeometryCollection.
func (m *Shape) IsGeometryCollection() bool {
	return m.Type == TypeGeometryCollection && m.GeometryCollection != nil
}

// IsEnvelope returns true whether valid Envelope.
func (m *Shape) IsEnvelope() bool {
	return m.Type == TypeEnvelope && m.Envelope != nil
}

// IsCircle returns true whether valid Circle.
func (m *Shape) IsCircle() bool {
	return m.Type == TypeCircle && m.Circle != nil
}

// UnmarshalJSON decodes the geometry data into a GeoJSON.
func (m *Shape) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

// MarshalJSON encodes the geomettry data into a GeoJSON.
func (m *Shape) MarshalJSON() ([]byte, error) {
	switch {
	case m.IsPoint():
		return json.Marshal(m.Point)
	case m.IsMultiPoint():
		return json.Marshal(m.MultiPoint)
	case m.IsLineString():
		return json.Marshal(m.LineString)
	case m.IsMultiLineString():
		return json.Marshal(m.MultiLineString)
	case m.IsPolygon():
		return json.Marshal(m.Polygon)
	case m.IsMultiPolygon():
		return json.Marshal(m.MultiPolygon)
	case m.IsGeometryCollection():
		return json.Marshal(m.GeometryCollection)
	case m.IsEnvelope():
		return json.Marshal(m.Envelope)
	case m.IsCircle():
		return json.Marshal(m.Circle)
	default:
		return nil, fmt.Errorf("geo: unknown type `%s`", m.Type)
	}
}

func (m *Shape) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	switch raw.Type {
	case TypePoint:
		point := &Point{}
		if err := point.decode(raw); err != nil {
			return err
		}
		m.Point = point
	case TypeMultiPoint:
		multipoint := &MultiPoint{}
		if err := multipoint.decode(raw); err != nil {
			return err
		}
		m.MultiPoint = multipoint
	case TypeLineString:
		lineString := &LineString{}
		if err := lineString.decode(raw); err != nil {
			return err
		}
		m.LineString = lineString
	case TypeMultiLineString:
		multiLine := &MultiLineString{}
		if err := multiLine.decode(raw); err != nil {
			return err
		}
		m.MultiLineString = multiLine
	case TypePolygon:
		polygon := &Polygon{}
		if err := polygon.decode(raw); err != nil {
			return err
		}
		m.Polygon = polygon
	case TypeMultiPolygon:
		multiPolygon := &MultiPolygon{}
		if err := multiPolygon.decode(raw); err != nil {
			return err
		}
		m.MultiPolygon = multiPolygon
	case TypeGeometryCollection:
		geometryCollection := &GeometryCollection{}
		if err := geometryCollection.decode(raw); err != nil {
			return err
		}
		m.GeometryCollection = geometryCollection
	case TypeEnvelope:
		envelope := &Envelope{}
		if err := envelope.decode(raw); err != nil {
			return err
		}
		m.Envelope = envelope
	case TypeCircle:
		circle := &Circle{}
		if err := circle.decode(raw); err != nil {
			return err
		}
		m.Circle = circle
	default:
		return fmt.Errorf("geo: unknown type `%s`", m.Type)
	}
	return nil
}
