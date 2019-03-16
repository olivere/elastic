package geo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	var (
		expected = NewPoint([]float64{100.0, 0.0})
		actual   = &Point{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Point", func(t *testing.T) {
		err := json.Unmarshal(point, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Point", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(point), string(actual))
	})
}

func TestMultiPoint(t *testing.T) {
	var (
		expected = NewMultiPoint([][]float64{
			{100.0, 0.0}, {101.0, 1.0},
		})
		actual = &MultiPoint{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal MultiPoint", func(t *testing.T) {
		err := json.Unmarshal(multiPoint, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal MultiPoint", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(multiPoint), string(actual))
	})
}

func TestLineString(t *testing.T) {
	var (
		expected = NewLineString([][]float64{
			{100.0, 0.0}, {101.0, 1.0},
		})
		actual = &LineString{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal LineString", func(t *testing.T) {
		err := json.Unmarshal(lineString, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal LineString", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(lineString), string(actual))
	})
}

func TestMultiLineString(t *testing.T) {
	var (
		expected = NewMultiLineString([][][]float64{
			[][]float64{
				{100.0, 0.0}, {101.0, 1.0},
			},
			[][]float64{
				{102.0, 2.0}, {103.0, 3.0},
			},
		})
		actual = &MultiLineString{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal MultiLineString", func(t *testing.T) {
		err := json.Unmarshal(multiLineString, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal MultiLineString", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(multiLineString), string(actual))
	})
}

func TestPolygonNoHoles(t *testing.T) {
	var (
		expected = NewPolygon([][][]float64{
			[][]float64{
				{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
			},
		})
		actual = &Polygon{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Polygon", func(t *testing.T) {
		err := json.Unmarshal(polygon, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Polygon", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(polygon), string(actual))
	})
}

func TestPolygonWithHoles(t *testing.T) {
	var (
		expected = NewPolygon([][][]float64{
			[][]float64{
				{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
			},
			[][]float64{
				{100.2, 0.2}, {100.8, 0.2}, {100.8, 0.8}, {100.2, 0.8}, {100.2, 0.2},
			},
		})
		actual = &Polygon{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Polygon With Holes", func(t *testing.T) {
		err := json.Unmarshal(polygonWithHoles, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Polygon With Holes", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(polygonWithHoles), string(actual))
	})
}

func TestMultiPolygon(t *testing.T) {
	var (
		expected = NewMultiPolygon([][][][]float64{
			{
				[][]float64{
					{102.0, 2.0}, {103.0, 2.0}, {103.0, 3.0}, {102.0, 3.0}, {102.0, 2.0},
				},
			},
			{
				[][]float64{
					{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
				},
				[][]float64{
					{100.2, 0.2}, {100.8, 0.2}, {100.8, 0.8}, {100.2, 0.8}, {100.2, 0.2},
				},
			},
		})
		actual = &MultiPolygon{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal MultiPolygon", func(t *testing.T) {
		err := json.Unmarshal(multiPolygon, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal MultiPolygon", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(multiPolygon), string(actual))
	})
}

func TestGeometryCollection(t *testing.T) {
	var (
		expected = NewGeometryCollection([]interface{}{
			&Point{
				Type:        "point",
				Coordinates: []float64{100.0, 0.0},
			},
			&LineString{
				Type: "linestring",
				Coordinates: [][]float64{
					{101.0, 0.0}, {102.0, 1.0},
				},
			},
		})
		actual = &GeometryCollection{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal GeometryCollection", func(t *testing.T) {
		err := json.Unmarshal(geometryCollection, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal GeometryCollection", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(geometryCollection), string(actual))
	})
}

func TestShape(t *testing.T) {
	var (
		expected = &Shape{
			Type: "polygon",
			Polygon: NewPolygon([][][]float64{
				[][]float64{
					{100.0, 0.0}, {101.0, 0.0}, {101.0, 1.0}, {100.0, 1.0}, {100.0, 0.0},
				},
			}),
		}
		actual = &Shape{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Shape", func(t *testing.T) {
		err := json.Unmarshal(polygon, actual)
		assert.NoError(err)
		assert.True(actual.IsPolygon())
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Shape", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(polygon), string(actual))
	})
}

func TestEnvelope(t *testing.T) {
	var (
		expected = NewEnvelope([][]float64{
			{100.0, 1.0}, {101.0, 0.0},
		})
		actual = &Envelope{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Envelope", func(t *testing.T) {
		err := json.Unmarshal(envelope, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Envelope", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(envelope), string(actual))
	})
}

func TestCircle(t *testing.T) {
	var (
		expected = NewCircle("25m", []float64{-109.874838, 44.439550})
		actual   = &Circle{}
	)
	assert := require.New(t)

	t.Run("Test Unmarshal Circle", func(t *testing.T) {
		err := json.Unmarshal(circle, actual)
		assert.NoError(err)
		assert.Equal(actual, expected)
	})

	t.Run("Test Marshal Circle", func(t *testing.T) {
		actual, err := json.Marshal(expected)
		assert.NoError(err)
		assert.JSONEq(string(circle), string(actual))
	})
}
