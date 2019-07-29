// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// AutoDateHistogramAggregation is a multi-bucket aggregation similar to the
// histogram except it can only be applied on date values, and the buckets num can bin pointed.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/7.2/search-aggregations-bucket-autodatehistogram-aggregation.html
type AutoDateHistogramAggregation struct {
	field           string
	script          *Script
	missing         interface{}
	subAggregations map[string]Aggregation
	meta            map[string]interface{}

	buckets           int
	order             string
	orderAsc          bool
	minDocCount       *int64
	extendedBoundsMin interface{}
	extendedBoundsMax interface{}
	timeZone          string
	format            string
	offset            string
	keyed             *bool
	minimumInterval   string
}

// NewAutoDateHistogramAggregation creates a new AutoDateHistogramAggregation.
func NewAutoDateHistogramAggregation() *AutoDateHistogramAggregation {
	return &AutoDateHistogramAggregation{
		subAggregations: make(map[string]Aggregation),
	}
}

// Field on which the aggregation is processed.
func (a *AutoDateHistogramAggregation) Field(field string) *AutoDateHistogramAggregation {
	a.field = field
	return a
}

// Script on which th
func (a *AutoDateHistogramAggregation) Script(script *Script) *AutoDateHistogramAggregation {
	a.script = script
	return a
}

// Missing configures the value to use when documents miss a value.
func (a *AutoDateHistogramAggregation) Missing(missing interface{}) *AutoDateHistogramAggregation {
	a.missing = missing
	return a
}

// SubAggregation sub aggregation
func (a *AutoDateHistogramAggregation) SubAggregation(name string, subAggregation Aggregation) *AutoDateHistogramAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *AutoDateHistogramAggregation) Meta(metaData map[string]interface{}) *AutoDateHistogramAggregation {
	a.meta = metaData
	return a
}

// Buckets buckets num by which the aggregation gets processed.
func (a *AutoDateHistogramAggregation) Buckets(buckets int) *AutoDateHistogramAggregation {
	a.buckets = buckets
	return a
}

// Order specifies the sort order. Valid values for order are:
// "_key", "_count", a sub-aggregation name, or a sub-aggregation name
// with a metric.
func (a *AutoDateHistogramAggregation) Order(order string, asc bool) *AutoDateHistogramAggregation {
	a.order = order
	a.orderAsc = asc
	return a
}

// OrderByCount specifies the sort order by count.
func (a *AutoDateHistogramAggregation) OrderByCount(asc bool) *AutoDateHistogramAggregation {
	// "order" : { "_count" : "asc" }
	a.order = "_count"
	a.orderAsc = asc
	return a
}

// OrderByCountAsc specifies the sort order by count with asc.
func (a *AutoDateHistogramAggregation) OrderByCountAsc() *AutoDateHistogramAggregation {
	return a.OrderByCount(true)
}

// OrderByCountDesc specifies the sort order by count with desc.
func (a *AutoDateHistogramAggregation) OrderByCountDesc() *AutoDateHistogramAggregation {
	return a.OrderByCount(false)
}

func (a *AutoDateHistogramAggregation) OrderByKey(asc bool) *AutoDateHistogramAggregation {
	// "order" : { "_key" : "asc" }
	a.order = "_key"
	a.orderAsc = asc
	return a
}

func (a *AutoDateHistogramAggregation) OrderByKeyAsc() *AutoDateHistogramAggregation {
	return a.OrderByKey(true)
}

func (a *AutoDateHistogramAggregation) OrderByKeyDesc() *AutoDateHistogramAggregation {
	return a.OrderByKey(false)
}

// OrderByAggregation creates a bucket ordering strategy which sorts buckets
// based on a single-valued calc get.
func (a *AutoDateHistogramAggregation) OrderByAggregation(aggName string, asc bool) *AutoDateHistogramAggregation {
	// {
	//     "aggs" : {
	//         "genders" : {
	//             "terms" : {
	//                 "field" : "gender",
	//                 "order" : { "avg_height" : "desc" }
	//             },
	//             "aggs" : {
	//                 "avg_height" : { "avg" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = aggName
	a.orderAsc = asc
	return a
}

// OrderByAggregationAndMetric creates a bucket ordering strategy which
// sorts buckets based on a multi-valued calc get.
func (a *AutoDateHistogramAggregation) OrderByAggregationAndMetric(aggName, metric string, asc bool) *AutoDateHistogramAggregation {
	// {
	//     "aggs" : {
	//         "genders" : {
	//             "terms" : {
	//                 "field" : "gender",
	//                 "order" : { "height_stats.avg" : "desc" }
	//             },
	//             "aggs" : {
	//                 "height_stats" : { "stats" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = aggName + "." + metric
	a.orderAsc = asc
	return a
}

// MinDocCount sets the minimum document count per bucket.
// Buckets with less documents than this min value will not be returned.
func (a *AutoDateHistogramAggregation) MinDocCount(minDocCount int64) *AutoDateHistogramAggregation {
	a.minDocCount = &minDocCount
	return a
}

// TimeZone sets the timezone in which to translate dates before computing buckets.
func (a *AutoDateHistogramAggregation) TimeZone(timeZone string) *AutoDateHistogramAggregation {
	a.timeZone = timeZone
	return a
}

// Format sets the format to use for dates.
func (a *AutoDateHistogramAggregation) Format(format string) *AutoDateHistogramAggregation {
	a.format = format
	return a
}

// Offset sets the offset of time intervals in the histogram, e.g. "+6h".
func (a *AutoDateHistogramAggregation) Offset(offset string) *AutoDateHistogramAggregation {
	a.offset = offset
	return a
}

// ExtendedBounds accepts int, int64, string, or time.Time values.
// In case the lower value in the histogram would be greater than min or the
// upper value would be less than max, empty buckets will be generated.
func (a *AutoDateHistogramAggregation) ExtendedBounds(min, max interface{}) *AutoDateHistogramAggregation {
	a.extendedBoundsMin = min
	a.extendedBoundsMax = max
	return a
}

// ExtendedBoundsMin accepts int, int64, string, or time.Time values.
func (a *AutoDateHistogramAggregation) ExtendedBoundsMin(min interface{}) *AutoDateHistogramAggregation {
	a.extendedBoundsMin = min
	return a
}

// ExtendedBoundsMax accepts int, int64, string, or time.Time values.
func (a *AutoDateHistogramAggregation) ExtendedBoundsMax(max interface{}) *AutoDateHistogramAggregation {
	a.extendedBoundsMax = max
	return a
}

// Keyed specifies whether to return the results with a keyed response (or not).
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.4/search-aggregations-bucket-datehistogram-aggregation.html#_keyed_response_3.
func (a *AutoDateHistogramAggregation) Keyed(keyed bool) *AutoDateHistogramAggregation {
	a.keyed = &keyed
	return a
}

// MinimumInterval accepted units for minimum_interval are: year/month/day/hour/minute/second
func (a *AutoDateHistogramAggregation) MinimumInterval(interval string) *AutoDateHistogramAggregation {
	a.minimumInterval = interval
	return a
}

// Source source for AutoDateHistogramAggregation
func (a *AutoDateHistogramAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "aggs" : {
	//         "articles_over_time" : {
	//             "auto_date_histogram" : {
	//                 "field" : "date",
	//                 "buckets" : 10
	//             }
	//         }
	//     }
	// }
	//
	// This method returns only the { "auto_date_histogram" : { ... } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["auto_date_histogram"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != nil {
		src, err := a.script.Source()
		if err != nil {
			return nil, err
		}
		opts["script"] = src
	}
	if a.missing != nil {
		opts["missing"] = a.missing
	}

	if a.buckets > 0 {
		opts["buckets"] = a.buckets
	} else {
		opts["buckets"] = 10
	}

	if a.minDocCount != nil {
		opts["min_doc_count"] = *a.minDocCount
	}
	if a.order != "" {
		o := make(map[string]interface{})
		if a.orderAsc {
			o[a.order] = "asc"
		} else {
			o[a.order] = "desc"
		}
		opts["order"] = o
	}
	if a.timeZone != "" {
		opts["time_zone"] = a.timeZone
	}
	if a.offset != "" {
		opts["offset"] = a.offset
	}
	if a.format != "" {
		opts["format"] = a.format
	}
	if a.extendedBoundsMin != nil || a.extendedBoundsMax != nil {
		bounds := make(map[string]interface{})
		if a.extendedBoundsMin != nil {
			bounds["min"] = a.extendedBoundsMin
		}
		if a.extendedBoundsMax != nil {
			bounds["max"] = a.extendedBoundsMax
		}
		opts["extended_bounds"] = bounds
	}
	if a.keyed != nil {
		opts["keyed"] = *a.keyed
	}
	if a.minimumInterval != "" {
		opts["minimumInterval"] = a.minimumInterval
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
