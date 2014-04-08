// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
)

// Aggregations can be seen as a unit-of-work that build
// analytic information over a set of documents. It is
// (in many senses) the follow-up of facets in Elasticsearch.
// For more details about aggregations, visit:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-aggregations.html
type Aggregation interface {
	Source() interface{}
}

// SearchAggregation is a generic aggregation from a search result.
// As aggregations in Elasticsearch are recursive, it could have
// any number of sub-aggregations; use GetAggregation to return the
// sub-aggregations.
type SearchAggregation struct {
	name         string
	raw          json.RawMessage
	aggregations map[string]json.RawMessage
}

// NewSearchAggregation initializes a new SearchAggregation.
func NewSearchAggregation(name string, raw json.RawMessage) *SearchAggregation {
	return &SearchAggregation{
		name: name,
		raw:  raw,
	}
}

// GetAggregation returns a sub-aggregation of this aggregation.
func (sa *SearchAggregation) GetAggregation(name string) (*SearchAggregation, bool) {
	if len(sa.aggregations) == 0 {
		pairs := make(map[string]json.RawMessage)
		if err := json.Unmarshal(sa.raw, &pairs); err != nil {
			return nil, false
		}
		sa.aggregations = pairs
	}

	if raw, found := sa.aggregations[name]; found {
		return NewSearchAggregation(name, raw), true
	}
	return nil, false
}

// Min treats this aggregation as a min aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-min-aggregation.html
func (sa *SearchAggregation) Min() (*SearchAggregationMin, bool) {
	agg := new(SearchAggregationMin)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Max treats this aggregation as a max aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-max-aggregation.html
func (sa *SearchAggregation) Max() (*SearchAggregationMax, bool) {
	agg := new(SearchAggregationMax)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Sum treats this aggregation as a sum aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-sum-aggregation.html
func (sa *SearchAggregation) Sum() (*SearchAggregationSum, bool) {
	agg := new(SearchAggregationSum)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Avg treats this aggregation as an avg aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-avg-aggregation.html
func (sa *SearchAggregation) Avg() (*SearchAggregationAvg, bool) {
	agg := new(SearchAggregationAvg)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Stats treats this aggregation as a stats aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-stats-aggregation.html
func (sa *SearchAggregation) Stats() (*SearchAggregationStats, bool) {
	agg := new(SearchAggregationStats)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// ExtendedStats treats this aggregation as an extended stats aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-extendedstats-aggregation.html
func (sa *SearchAggregation) ExtendedStats() (*SearchAggregationExtendedStats, bool) {
	agg := new(SearchAggregationExtendedStats)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// ValueCount treats this aggregation as a value count aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-valuecount-aggregation.html
func (sa *SearchAggregation) ValueCount() (*SearchAggregationValueCount, bool) {
	agg := new(SearchAggregationValueCount)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Percentiles treats this aggregation as a percentiles aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-percentile-aggregation.html
func (sa *SearchAggregation) Percentiles() (*SearchAggregationPercentiles, bool) {
	values := make(map[string]interface{})
	if err := json.Unmarshal(sa.raw, &values); err != nil {
		return nil, false
	}
	return &SearchAggregationPercentiles{Values: values}, true
}

// Cardinality treats this aggregation as a cardinality aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-metrics-cardinality-aggregation.html
func (sa *SearchAggregation) Cardinality() (*SearchAggregationCardinality, bool) {
	agg := new(SearchAggregationCardinality)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Global treats this aggregation as a global aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-global-aggregation.html
func (sa *SearchAggregation) Global() (*SearchAggregationGlobal, bool) {
	agg := new(SearchAggregationGlobal)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Global treats this aggregation as a filter aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-filter-aggregation.html
func (sa *SearchAggregation) Filter() (*SearchAggregationFilter, bool) {
	agg := new(SearchAggregationFilter)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Missing treats this aggregation as a missing aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-missing-aggregation.html
func (sa *SearchAggregation) Missing() (*SearchAggregationMissing, bool) {
	agg := new(SearchAggregationMissing)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Nested treats this aggregation as a nested aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-nested-aggregation.html
func (sa *SearchAggregation) Nested() (*SearchAggregationNested, bool) {
	agg := new(SearchAggregationNested)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Terms treats this aggregation as a terms aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-terms-aggregation.html
func (sa *SearchAggregation) Terms() (*SearchAggregationTerms, bool) {
	agg := new(SearchAggregationTerms)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// SignificantTerms treats this aggregation as a significant terms aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-significantterms-aggregation.html
func (sa *SearchAggregation) SignificantTerms() (*SearchAggregationSignificantTerms, bool) {
	agg := new(SearchAggregationSignificantTerms)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Range treats this aggregation as a range aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-range-aggregation.html
func (sa *SearchAggregation) Range() (*SearchAggregationRange, bool) {
	agg := new(SearchAggregationRange)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// DateRange treats this aggregation as a date range aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-daterange-aggregation.html
func (sa *SearchAggregation) DateRange() (*SearchAggregationDateRange, bool) {
	agg := new(SearchAggregationDateRange)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// Histogram treats this aggregation as a histogram aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-histogram-aggregation.html
func (sa *SearchAggregation) Histogram() (*SearchAggregationHistogram, bool) {
	agg := new(SearchAggregationHistogram)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// DateHistogram treats this aggregation as a date histogram aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-datehistogram-aggregation.html
func (sa *SearchAggregation) DateHistogram() (*SearchAggregationDateHistogram, bool) {
	agg := new(SearchAggregationDateHistogram)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

// GeoDistance treats this aggregation as a geo distance aggregation.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/master/search-aggregations-bucket-geodistance-aggregation.html
func (sa *SearchAggregation) GeoDistance() (*SearchAggregationGeoDistance, bool) {
	agg := new(SearchAggregationGeoDistance)
	if err := json.Unmarshal(sa.raw, &agg); err != nil {
		return nil, false
	}
	return agg, true
}

type SearchAggregationMin struct {
	Value         interface{} `json:"value"`
	ValueAsString string      `json:"value_as_string,omitempty"`
}

type SearchAggregationMax struct {
	Value         interface{} `json:"value"`
	ValueAsString string      `json:"value_as_string,omitempty"`
}

type SearchAggregationSum struct {
	Value         interface{} `json:"value"`
	ValueAsString string      `json:"value_as_string,omitempty"`
}

type SearchAggregationAvg struct {
	Value         interface{} `json:"value"`
	ValueAsString string      `json:"value_as_string,omitempty"`
}

type SearchAggregationStats struct {
	Count int         `json:"count,omitempty"`
	Min   interface{} `json:"min,omitempty"`
	Max   interface{} `json:"max,omitempty"`
	Avg   interface{} `json:"avg,omitempty"`
	Sum   interface{} `json:"sum,omitempty"`
}

type SearchAggregationExtendedStats struct {
	Count        int         `json:"count,omitempty"`
	Min          interface{} `json:"min,omitempty"`
	Max          interface{} `json:"max,omitempty"`
	Avg          interface{} `json:"avg,omitempty"`
	Sum          interface{} `json:"sum,omitempty"`
	SumOfSquares interface{} `json:"sum_of_squares,omitempty"`
	Variance     interface{} `json:"variance,omitempty"`
	StdDeviation interface{} `json:"std_deviation,omitempty"`
}

type SearchAggregationValueCount struct {
	Value         int    `json:"value"`
	ValueAsString string `json:"value_as_string,omitempty"`
}

type SearchAggregationPercentiles struct {
	Values map[string]interface{}
}

type SearchAggregationCardinality struct {
	Value         int    `json:"value"`
	ValueAsString string `json:"value_as_string,omitempty"`
}

type SearchAggregationGlobal struct {
	DocCount int `json:"doc_count,omitempty"`
}

type SearchAggregationFilter struct {
	DocCount int `json:"doc_count,omitempty"`
}

type SearchAggregationMissing struct {
	DocCount int `json:"doc_count,omitempty"`
}

type SearchAggregationNested struct {
	Value         int    `json:"value"`
	ValueAsString string `json:"value_as_string,omitempty"`
}

type SearchAggregationTerms struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationSignificantTerms struct {
	DocCount int                        `json:"doc_count,omitempty"`
	Buckets  []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationRange struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationDateRange struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationHistogram struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationDateHistogram struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type SearchAggregationGeoDistance struct {
	Buckets []*searchAggregationBucket `json:"buckets,omitempty"`
}

type searchAggregationBucket struct {
	Key          interface{} `json:"key,omitempty"`
	KeyAsString  *string     `json:"key_as_string,omitempty"`
	DocCount     int         `json:"doc_count,omitempty"`
	From         *float64    `json:"from,omitempty"`
	FromAsString *string     `json:"from_as_string,omitempty"`
	To           *float64    `json:"to,omitempty"`
	ToAsString   *string     `json:"to_as_string,omitempty"`
	Unit         string      `json:"unit,omitempty"`
	Score        *float64    `json:"score,omitempty"`    // significant_terms
	BgCount      *int        `json:"bg_count,omitempty"` // significant_terms
}
