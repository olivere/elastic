package elastic

import (
	"fmt"
	"strconv"
	"strings"
)

// GeoPoint is a geographic position described via latitude and longitude.
type GeoPoint struct {
	Lat, Lon float64
}

// Source returns the object to be serialized in Elasticsearch DSL.
func (pt *GeoPoint) Source() map[string]float64 {
	return map[string]float64{
		"lat": pt.Lat,
		"lon": pt.Lon,
	}
}

// GeoPointFromLatLon initializes a new GeoPoint by latitude and longitude.
func GeoPointFromLatLon(lat, lon float64) *GeoPoint {
	return &GeoPoint{Lat: lat, Lon: lon}
}

// GeoPointFromString initializes a new GeoPoint by a string that is
// formatted as "{latitude},{longitude}", e.g. "40.10210,-70.12091".
func GeoPointFromString(latLon string) (*GeoPoint, error) {
	latlon := strings.SplitN(latLon, ",", 2)
	if len(latlon) != 2 {
		return nil, fmt.Errorf("elastic: %s is not a valid geo point string", latLon)
	}
	lat, err := strconv.ParseFloat(latlon[0], 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(latlon[1], 64)
	if err != nil {
		return nil, err
	}
	return &GeoPoint{Lat: lat, Lon: lon}, nil
}
