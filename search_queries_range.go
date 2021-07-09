// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// RangeQuery matches documents with fields that have terms within a certain range.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/7.0/query-dsl-range-query.html
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
	relation  string
}

// NewRangeQuery creates and initializes a new RangeQuery.
func NewRangeQuery(name string) *RangeQuery {
	return &RangeQuery{name: name}
}

// From indicates the from part of the RangeQuery.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) From(from interface{}) *RangeQuery {
	return q.Gte(from)
}

// Gt indicates a greater-than value for the from part.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) Gt(from interface{}) *RangeQuery {
	q.gt = from
	return q
}

// Gte indicates a greater-than-or-equal value for the from part.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) Gte(from interface{}) *RangeQuery {
	q.gte = from
	return q
}

// To indicates the to part of the RangeQuery.
// Use nil to indicate an unbounded to part.
func (q *RangeQuery) To(to interface{}) *RangeQuery {
	return q.Lte(to)
}

// Lt indicates a less-than value for the to part.
// Use nil to indicate an unbounded to part.
func (q *RangeQuery) Lt(to interface{}) *RangeQuery {
	q.lt = to
	return q
}

// Lte indicates a less-than-or-equal value for the to part.
// Use nil to indicate an unbounded to part.
func (q *RangeQuery) Lte(to interface{}) *RangeQuery {
	q.lte = to
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

// Relation is used for range fields. which can be one of
// "within", "contains", "intersects" (default) and "disjoint".
func (q *RangeQuery) Relation(relation string) *RangeQuery {
	q.relation = relation
	return q
}

// Source returns JSON for the query.
func (q *RangeQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	rangeQ := make(map[string]interface{})
	source["range"] = rangeQ

	params := make(map[string]interface{})
	rangeQ[q.name] = params

	if nil != q.gt {
		params["gt"] = q.gt
	}

	if nil != q.gte {
		params["gte"] = q.gte
	}

	if nil != q.lt {
		params["lt"] = q.lt
	}

	if nil != q.lte {
		params["lte"] = q.lte
	}

	if q.timeZone != "" {
		params["time_zone"] = q.timeZone
	}
	if q.format != "" {
		params["format"] = q.format
	}
	if q.relation != "" {
		params["relation"] = q.relation
	}
	if q.boost != nil {
		params["boost"] = *q.boost
	}

	if q.queryName != "" {
		rangeQ["_name"] = q.queryName
	}

	return source, nil
}
