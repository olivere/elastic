// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.
// Contributor: Yasar Senturk <yasar@senturk.name.tr>, 2014

package elastic

// Matches documents that have fields containing terms
// with a query contains wildcards.
// For more details, see
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-wildcard-query.html
type WildcardQuery struct {
	Query
	name      string
	wildcard  string
	boost     *float32
	rewrite   string
	queryName string
}

// Creates a new wildcard query.
func NewWildcardQuery(name string, wildcard string) WildcardQuery {
	q := WildcardQuery{name: name, wildcard: wildcard}
	return q
}

func (q WildcardQuery) Boost(boost float32) WildcardQuery {
	q.boost = &boost
	return q
}

func (q WildcardQuery) Rewrite(rewrite string) WildcardQuery {
	q.rewrite = rewrite
	return q
}

func (q WildcardQuery) QueryName(queryName string) WildcardQuery {
	q.queryName = queryName
	return q
}

// Creates the query source for the wildcard query.
func (q WildcardQuery) Source() interface{} {
	// {
	//    "wildcard" : { "user" : { "wildcard" : "ki*y", "boost" : 2.0 } }
	// }

	source := make(map[string]interface{})

	query := make(map[string]interface{})
	source["wildcard"] = query

	if q.boost == nil && q.rewrite == "" && q.queryName == "" {
		query[q.name] = q.wildcard
	} else {
		subQuery := make(map[string]interface{})
		subQuery["wildcard"] = q.wildcard
		if q.boost != nil {
			subQuery["boost"] = *q.boost
		}
		if q.rewrite != "" {
			subQuery["rewrite"] = q.rewrite
		}
		if q.queryName != "" {
			subQuery["_name"] = q.queryName
		}
		query[q.name] = subQuery
	}

	return source
}
