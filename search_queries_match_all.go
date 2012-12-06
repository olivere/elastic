// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A query that matches all documents. For more details, see
// http://www.elasticsearch.org/guide/reference/query-dsl/match-all-query.html
type MatchAllQuery struct {
	Query

	normsField string
	boost *float32
}

// Creates a new match all query.
func NewMatchAllQuery() MatchAllQuery {
	q := MatchAllQuery{}
	return q
}

func (q MatchAllQuery) NormsField(normsField string) MatchAllQuery {
	q.normsField = normsField
	return q
}

func (q MatchAllQuery) Boost(boost float32) MatchAllQuery {
	q.boost = &boost
	return q
}

// Creates the query source for the match all query.
func (q MatchAllQuery) Source() interface{} {
	// {
	//   "match_all" : { ... }
	// }

	source := make(map[string]interface{})

	params := make(map[string]interface{})
	source["match_all"] = params

	if q.boost != nil {
		params["boost"] = q.boost
	}

	if q.normsField != "" {
		params["norms_field"] = q.normsField
	}

	return source
}
