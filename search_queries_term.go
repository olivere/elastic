// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A term query matches documents that contain
// a term (not analyzed). For more details, see
// http://www.elasticsearch.org/guide/reference/query-dsl/term-query.html
type TermQuery struct {
	Query
	name  string
	value interface{}
}

// Creates a new term query.
func NewTermQuery(name string, value interface{}) TermQuery {
	t := TermQuery{name: name, value: value}
	return t
}

// Creates the query source for the term query.
func (q TermQuery) Source() interface{} {
	// {"term":{"name":"value"}}
	source := make(map[string]interface{})
	tq := make(map[string]interface{})
	tq[q.name] = q.value
	source["term"] = tq
	return source
}
