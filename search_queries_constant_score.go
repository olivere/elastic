// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// ConstantScoreQuery wraps a filter or another query and simply returns
// a constant score equal to the query boost for every document in the filter.
//
// For more details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/1.7/query-dsl-constant-score-query.html
type ConstantScoreQuery struct {
	query  Query
	filter Filter
	boost  *float64
}

// NewConstantScoreQuery creates a new constant score query.
func NewConstantScoreQuery() ConstantScoreQuery {
	return ConstantScoreQuery{}
}

// Query to wrap in this constant score query.
func (q ConstantScoreQuery) Query(query Query) ConstantScoreQuery {
	q.query = query
	q.filter = nil
	return q
}

// Filter to wrap in this constant score query.
func (q ConstantScoreQuery) Filter(filter Filter) ConstantScoreQuery {
	q.query = nil
	q.filter = filter
	return q
}

// Boost sets the boost for this query. Documents matching this query
// will (in addition to the normal weightings) have their score multiplied
// by the boost provided.
func (q ConstantScoreQuery) Boost(boost float64) ConstantScoreQuery {
	q.boost = &boost
	return q
}

// Source returns JSON for the function score query.
func (q ConstantScoreQuery) Source() interface{} {
	source := make(map[string]interface{})
	query := make(map[string]interface{})
	source["constant_score"] = query

	if q.query != nil {
		query["query"] = q.query.Source()
	} else if q.filter != nil {
		query["filter"] = q.filter.Source()
	}
	if q.boost != nil {
		query["boost"] = *q.boost
	}
	return source
}
