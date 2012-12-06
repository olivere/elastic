// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A query that applies a filter to the results of another query.
// For more details, see
// http://www.elasticsearch.org/guide/reference/query-dsl/filtered-query.html
type FilteredQuery struct {
	Query
	query  Query
	filter *Filter
	boost  *float32
}

// Creates a new filtered query.
func NewFilteredQuery(query Query, filter *Filter) FilteredQuery {
	q := FilteredQuery{query: query, filter: filter}
	return q
}

func (q FilteredQuery) Boost(boost float32) FilteredQuery {
	q.boost = &boost
	return q
}

// Creates the query source for the filtered query.
func (q FilteredQuery) Source() interface{} {
	// {
	//     "filtered" : {
	//         "query" : {
	//             "term" : { "tag" : "wow" }
	//         },
	//         "filter" : {
	//             "range" : {
	//                 "age" : { "from" : 10, "to" : 20 }
	//             }
	//         }
	//     }
	// }

	source := make(map[string]interface{})

	filtered := make(map[string]interface{})
	source["filtered"] = filtered

	filtered["query"] = q.query.Source()

	if q.filter != nil {
		filtered["filter"] = (*q.filter).Source()
	}

	if q.boost != nil {
		filtered["boost"] = *q.boost
	}

	return source
}
