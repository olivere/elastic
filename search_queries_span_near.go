// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// SpanNearQuery filters Matches spans which are near one another.
//
// For more details, see
//https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-span-near-query.html
type SpanNearQuery struct {
	clauses []*clause `json:"clauses"`
	slop    int       `json:"slop"`
	inOrder bool      `json:"in_order"`
}

// clause
type clause map[string]interface{}

// NewSpanNearQuery creates a new clause field
func NewClause(field, value string) *clause {
	return &clause{
		"span_term": map[string]interface{}{
			field: value,
		},
	}
}

// NewSpanNearQuery creates and initializes a new SpanNearQuery.
func NewSpanNearQuery(clauses ...*clause) *SpanNearQuery {
	q := &SpanNearQuery{
		clauses: make([]*clause, 0),
	}
	if len(clauses) > 0 {
		q.clauses = append(q.clauses, clauses...)
	}
	return q
}

//AddClause Add a span clause to the current list of clauses
func (q *SpanNearQuery) AddClause(clause *clause) *SpanNearQuery {
	if clause != nil {
		q.clauses = append(q.clauses, clause)
	}
	return q
}

//InOrder must be in the same order as in clauses and must be non-overlapping
func (q *SpanNearQuery) InOrder(inOrder bool) *SpanNearQuery {
	q.inOrder = inOrder
	return q
}

//Slop set the maximum number of intervening unmatched positions permitted
func (q *SpanNearQuery) Slop(slop int) *SpanNearQuery {
	q.slop = slop
	return q
}

// Creates the query source for the span near query.
func (q *SpanNearQuery) Source() (interface{}, error) {
	query := make(map[string]interface{})
	span := make(map[string]interface{})
	query["span_near"] = span

	span["clauses"] = q.clauses
	span["slop"] = q.slop
	span["in_order"] = q.inOrder
	return query, nil
}
