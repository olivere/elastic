package elastic

// DistanceFeatureQuery filters documents that include only hits that exists
// within a specific distance from a geo point.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.x/query-dsl-distance-feature-query.html
type DistanceFeatureQuery struct {
	field     string
	origin    interface{}
	pivot     string
	boost     *float64
	queryName string
}

// NewDistanceFeatureQuery creates and initializes a new GeoDistanceQuery.
func NewDistanceFeatureQuery(field string, origin interface{}, pivot string) *DistanceFeatureQuery {
	return &DistanceFeatureQuery{field: field, origin: origin, pivot: pivot}
}

// Boost sets the boost for this query.
func (q *DistanceFeatureQuery) Boost(boost float64) *DistanceFeatureQuery {
	q.boost = &boost
	return q
}

// QueryName sets the query name for the filter.
func (q *DistanceFeatureQuery) QueryName(queryName string) *DistanceFeatureQuery {
	q.queryName = queryName
	return q
}

// Source returns JSON for the function score query.
func (q *DistanceFeatureQuery) Source() (interface{}, error) {
	// 	{
	// 		"distance_feature": {
	// 			"field": "coordinates",
	// 			"pivot": "10km",
	// 			"boost": 3.4,
	// 			"origin": {
	// 				"lat": 11.1111111,
	// 				"lon": 11.1111111
	// 			}
	//   	}
	// 	}

	source := make(map[string]interface{})
	query := make(map[string]interface{})

	query["field"] = q.field
	query["pivot"] = q.pivot
	query["origin"] = q.origin

	if q.boost != nil {
		query["boost"] = *q.boost
	}

	source["distance_feature"] = query

	return source, nil
}
