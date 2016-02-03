// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// RangeQuery matches documents with fields that have terms within a certain range.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-range-query.html
type RangeQuery struct {
	name      string
	gt        interface{}
	gte       interface{}
	lt        interface{}
	lte       interface{}
	timeZone  string
	boost     *float64
	queryName string
	format    string
}

// NewRangeQuery creates and initializes a new RangeQuery.
func NewRangeQuery(name string) *RangeQuery {
	return &RangeQuery{name: name}
}

// From indicates the `gte` part of the RangeQuery.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) From(from interface{}) *RangeQuery {
	return q.Gte(from)
}

// To indicates the `lte` part of the RangeQuery.
func (q *RangeQuery) To(to interface{}) *RangeQuery {
	return q.Lte(to)
}

// Gt indicates a greater-than value.
// Use nil to indicate an unbounded gt part.
func (q *RangeQuery) Gt(gt interface{}) *RangeQuery {
	q.gt = gt
	return q
}

// Gte indicates a greater-than-equal value.
// Use nil to indicate an unbounded gte part.
func (q *RangeQuery) Gte(gte interface{}) *RangeQuery {
	q.gte = gte
	return q
}

// Lt indicates a less-than value.
// Use nil to indicate an unbounded `lt` part.
func (q *RangeQuery) Lt(lt interface{}) *RangeQuery {
	q.lt = lt
	return q
}

// Lte indicates a less-than-or-equal value.
// Use nil to indicate an unbounded `lte` part.
func (q *RangeQuery) Lte(lte interface{}) *RangeQuery {
	q.lte = lte
	return q
}

// IncludeLower indicates whether the lower bound should be included or not.
// Defaults to true.
func (q *RangeQuery) IncludeLower(includeLower bool) *RangeQuery {
	if includeLower {
		q.lte = q.lt
		q.lt = nil
	} else {
		q.lt = q.lte
		q.lte = nil
	}
	return q
}

// IncludeUpper indicates whether the upper bound should be included or not.
// Defaults to true.
func (q *RangeQuery) IncludeUpper(includeUpper bool) *RangeQuery {
	if includeUpper {
		q.gte = q.gt
		q.gt = nil
	} else {
		q.gt = q.gte
		q.gte = nil
	}
	return q
}

// Boost sets the boost for this query.
func (q *RangeQuery) Boost(boost float64) *RangeQuery {
	q.boost = &boost
	return q
}

// QueryName sets the query name for the filter that can be used when
// searching for matched_filters per hit.
func (q *RangeQuery) QueryName(queryName string) *RangeQuery {
	q.queryName = queryName
	return q
}

// TimeZone is used for date fields. In that case, we can adjust the
// from/to fields using a timezone.
func (q *RangeQuery) TimeZone(timeZone string) *RangeQuery {
	q.timeZone = timeZone
	return q
}

// Format is used for date fields. In that case, we can set the format
// to be used instead of the mapper format.
func (q *RangeQuery) Format(format string) *RangeQuery {
	q.format = format
	return q
}

// Source returns JSON for the query.
func (q *RangeQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	rangeQ := make(map[string]interface{})
	source["range"] = rangeQ

	params := make(map[string]interface{})

	if q.gt != nil {
		params["gt"] = q.gt
	}

	if q.gte != nil {
		params["gte"] = q.gte
	}

	if q.lt != nil {
		params["lt"] = q.lt
	}

	if q.lte != nil {
		params["lte"] = q.lte
	}

	rangeQ[q.name] = params

	if q.timeZone != "" {
		params["time_zone"] = q.timeZone
	}
	if q.format != "" {
		params["format"] = q.format
	}
	if q.boost != nil {
		rangeQ["boost"] = *q.boost
	}

	if q.queryName != "" {
		rangeQ["_name"] = q.queryName
	}

	return source, nil
}
