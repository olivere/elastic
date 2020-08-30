// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-range-query.html#querying-range-fields
const (
	// (Default) Matches documents with a range field value that intersects the query’s range.
	RelationIntersects string = "INTERSECTS"
	// Matches documents with a range field value that entirely contains the query’s range.
	RelationContains string = "CONTAINS"
	// Matches documents with a range field value entirely within the query’s range.
	RelationWithin string = "WITHIN"
)

// RangeQuery matches documents with fields that have terms within a certain range.
//
// For details, see Elastic Documentation (6.8):
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-range-query.html
type RangeQuery struct {
	name         string
	gt           interface{}
	gte          interface{}
	lt           interface{}
	lte          interface{}
	timeZone     string
	includeLower bool
	includeUpper bool
	boost        *float64
	queryName    string
	format       string
	relation     string
}

// NewRangeQuery creates and initializes a new RangeQuery.
func NewRangeQuery(name string) *RangeQuery {
	return &RangeQuery{name: name, includeLower: true, includeUpper: true}
}

// Gt indicates a greater-than value for the from part.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) Gt(gt interface{}) *RangeQuery {
	q.gt = gt
	return q
}

// Gte indicates a greater-than-or-equal value for the from part.
// Use nil to indicate an unbounded from part.
func (q *RangeQuery) Gte(gte interface{}) *RangeQuery {
	q.gte = gte
	return q
}

// Lt indicates a less-than value for the to part.
// Use nil to indicate an unbounded to part.
func (q *RangeQuery) Lt(lt interface{}) *RangeQuery {
	q.lt = lt
	return q
}

// Lte indicates a less-than-or-equal value for the to part.
// Use nil to indicate an unbounded to part.
func (q *RangeQuery) Lte(lte interface{}) *RangeQuery {
	q.lte = lte
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
// "within", "contains" and "intersects" (default).
func (q *RangeQuery) Relation(relation string) *RangeQuery {
	q.relation = relation
	return q
}

// From Deprecated use Gt or Gte
func (q *RangeQuery) From(from interface{}) *RangeQuery {
	if q.includeLower {
		q.gte = from
		q.gt = nil
		return q
	}
	q.gte = nil
	q.gt = from
	return q
}

// To Deprecated use Lt or Lte
func (q *RangeQuery) To(to interface{}) *RangeQuery {
	if q.includeUpper {
		q.lte = to
		q.lt = nil
		return q
	}
	q.lte = nil
	q.lt = to
	return q
}

// IncludeLower Deprecated use Gt or Gte
func (q *RangeQuery) IncludeLower(includeLower bool) *RangeQuery {
	if includeLower && q.gt != nil {
		q.gte = q.gt
		q.gt = nil
	}
	if !includeLower && q.gte != nil {
		q.gt = q.gte
		q.gte = nil
	}
	q.includeLower = includeLower
	return q
}

// IncludeUpper Deprecated use Lt or Lte
func (q *RangeQuery) IncludeUpper(includeUpper bool) *RangeQuery {
	if includeUpper && q.lt != nil {
		q.lte = q.lt
		q.lt = nil
	}
	if !includeUpper && q.lte != nil {
		q.lt = q.lte
		q.lte = nil
	}

	q.includeUpper = includeUpper
	return q
}

// Source returns JSON for the query.
func (q *RangeQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	rangeQ := make(map[string]interface{})
	source["range"] = rangeQ

	params := make(map[string]interface{})
	rangeQ[q.name] = params

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
