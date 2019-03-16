package geo

// Option allows to configure various aspects of Shape.
type Option func(*Shape)

// WithPoint set the shape as Point.
func WithPoint(coordinates []float64) Option {
	return func(s *Shape) {
		s.Type = TypePoint
		s.Point = NewPoint(coordinates)
	}
}

// WithMultiPoint set the shape as MultiPoint.
func WithMultiPoint(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiPoint
		s.MultiPoint = NewMultiPoint(coordinates)
	}
}

// WithLineString set the shape as LineString.
func WithLineString(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeLineString
		s.LineString = NewLineString(coordinates)
	}
}

// WithMultiLineString set the shape as MultiLineString.
func WithMultiLineString(coordinates [][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiLineString
		s.MultiLineString = NewMultiLineString(coordinates)
	}
}

// WithPolygon set the shape as Polygon.
func WithPolygon(coordinates [][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypePolygon
		s.Polygon = NewPolygon(coordinates)
	}
}

// WithMultiPolygon set the shape as MultiPolygon.
func WithMultiPolygon(coordinates [][][][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeMultiPolygon
		s.MultiPolygon = NewMultiPolygon(coordinates)
	}
}

// WithGeometryCollection set the shape as GeometryCollection.
func WithGeometryCollection(geometries []interface{}) Option {
	return func(s *Shape) {
		s.Type = TypeGeometryCollection
		s.GeometryCollection = NewGeometryCollection(geometries)
	}
}

// WithEnvelope set the shape as Envelope.
func WithEnvelope(coordinates [][]float64) Option {
	return func(s *Shape) {
		s.Type = TypeEnvelope
		s.Envelope = NewEnvelope(coordinates)
	}
}

// WithCircle set the shape as Circle.
func WithCircle(radius string, coordinates []float64) Option {
	return func(s *Shape) {
		s.Type = TypeCircle
		s.Circle = NewCircle(radius, coordinates)
	}
}
