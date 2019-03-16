package geo

import (
	"encoding/json"
	"fmt"
)

// GeometryCollection is a geometry object which represents a collection of geometry objects.
type GeometryCollection struct {
	Type       string        `json:"type"`
	Geometries []interface{} `json:"geometries"`
}

// NewGeometryCollection creates and initializes a new GeometryCollection.
func NewGeometryCollection(geometries []interface{}) *GeometryCollection {
	return &GeometryCollection{
		Type:       TypeGeometryCollection,
		Geometries: geometries,
	}
}

// UnmarshalJSON decodes the geometry collection data into a GeoJSON geometry.
func (m *GeometryCollection) UnmarshalJSON(data []byte) error {
	raw := &rawGeometry{}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}
	return m.decode(raw)
}

func (m *GeometryCollection) decode(raw *rawGeometry) error {
	m.Type = raw.Type
	m.Geometries = make([]interface{}, 0, len(raw.Geometries))
	for _, geometry := range raw.Geometries {
		switch geometry.Type {
		case TypePoint:
			point := &Point{}
			if err := point.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, point)
		case TypeMultiPoint:
			multipoint := &MultiPoint{}
			if err := multipoint.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, multipoint)
		case TypeLineString:
			lineString := &LineString{}
			if err := lineString.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, lineString)
		case TypeMultiLineString:
			multiLine := &MultiLineString{}
			if err := multiLine.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, multiLine)
		case TypePolygon:
			polygon := &Polygon{}
			if err := polygon.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, polygon)
		case TypeMultiPolygon:
			multiPolygon := &MultiPolygon{}
			if err := multiPolygon.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, multiPolygon)
		case TypeGeometryCollection:
			geometryCollection := &GeometryCollection{}
			if err := geometryCollection.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, geometryCollection)
		case TypeEnvelope:
			envelope := &Envelope{}
			if err := envelope.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, envelope)
		case TypeCircle:
			circle := &Circle{}
			if err := circle.decode(&geometry); err != nil {
				return err
			}
			m.Geometries = append(m.Geometries, circle)
		default:
			return fmt.Errorf("geo: unknown type `%s`", m.Type)
		}
	}
	return nil
}
