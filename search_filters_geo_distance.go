// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Filters documents that include only hits that exists within a specific distance from a geo point.
//
// For more details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-geo-distance-filter.html
type GeoDistanceFilter struct {
	Filter
	name          string
	distance      string
	lat           float64
	lon           float64
	distance_type string
	optimize_bbox string
	cache         *bool
}

// Creates a new dis_max query.
func NewGeoDistanceFilter(name string) GeoDistanceFilter {
	f := GeoDistanceFilter{name: name}
	return f
}

func (f GeoDistanceFilter) Distance(distance string) GeoDistanceFilter {
	f.distance = distance
	return f
}

func (f GeoDistanceFilter) Lat(lat float64) GeoDistanceFilter {
	f.lat = lat
	return f
}

func (f GeoDistanceFilter) Lon(lon float64) GeoDistanceFilter {
	f.lon = lon
	return f
}

func (f GeoDistanceFilter) DistanceType(distance_type string) GeoDistanceFilter {
	f.distance_type = distance_type
	return f
}

func (f GeoDistanceFilter) OptimizeBbox(optimize_bbox string) GeoDistanceFilter {
	f.optimize_bbox = optimize_bbox
	return f
}

func (f GeoDistanceFilter) Cache(cache bool) GeoDistanceFilter {
	f.cache = &cache
	return f
}

// Creates the query source for the geo_distance filter.
func (f GeoDistanceFilter) Source() interface{} {
	// {
	//   "geo_distance" : {
	//       "distance" : "200km",
	//       "pin.location" : {
	//           "lat" : 40,
	//           "lon" : -70
	//       }
	//   }
	// }

	source := make(map[string]interface{})

	params := make(map[string]interface{})

	if f.distance != "" {
		params["distance"] = f.distance
	}

	if f.distance_type != "" {
		params["distance_type"] = f.distance_type
	}

	if f.optimize_bbox != "" {
		params["optimize_bbox"] = f.optimize_bbox
	}

	if f.cache != nil {
		params["_cache"] = *f.cache
	}

	source["geo_distance"] = params

	location := make(map[string]interface{})
	location["lat"] = f.lat
	location["lon"] = f.lon
	params[f.name] = location

	return source
}
