package elastic

import "github.com/olivere/elastic/v7"

// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// SpanOrQuery matches spans which are near one another. One can specify slop,
// the maximum number of intervening unmatched positions, as well as whether
// matches are required to be in-order. The span near query maps to Lucene SpanOrQuery.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.7/query-dsl-span-near-query.html
// for details.
type SpanOrQuery struct {
	clauses   []elastic.Query
	boost     *float64
	queryName string
}

// NewSpanOrQuery creates a new NewSpanOrQuery.
func NewSpanOrQuery(clauses ...elastic.Query) *SpanOrQuery {
	return &SpanOrQuery{
		clauses: clauses,
	}
}

// Add clauses to use in the query.
func (q *SpanOrQuery) Add(clauses ...elastic.Query) *SpanOrQuery {
	q.clauses = append(q.clauses, clauses...)
	return q
}

// Clauses to use in the query.
func (q *SpanOrQuery) Clauses(clauses ...elastic.Query) *SpanOrQuery {
	q.clauses = clauses
	return q
}

// Boost sets the boost for this query.
func (q *SpanOrQuery) Boost(boost float64) *SpanOrQuery {
	q.boost = &boost
	return q
}

// QueryName sets the query name for the filter that can be used when
// searching for matched_filters per hit.
func (q *SpanOrQuery) QueryName(queryName string) *SpanOrQuery {
	q.queryName = queryName
	return q
}

// Source returns the JSON body.
func (q *SpanOrQuery) Source() (interface{}, error) {
	m := make(map[string]interface{})
	c := make(map[string]interface{})

	if len(q.clauses) > 0 {
		var clauses []interface{}
		for _, clause := range q.clauses {
			src, err := clause.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		c["clauses"] = clauses
	}

	if v := q.boost; v != nil {
		c["boost"] = *v
	}
	if v := q.queryName; v != "" {
		c["query_name"] = v
	}
	m["span_or"] = c
	return m, nil
}
