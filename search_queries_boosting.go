// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A boosting query can be used to effectively
// demote results that match a given query.
// For more details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-boosting-query.html
type BoostingQuery struct {
	Query
	positiveClauses []Query
	negativeClauses []Query
	negativeBoost   *float64
	boost           *float64
}

// Creates a new boosting query.
func NewBoostingQuery() BoostingQuery {
	return BoostingQuery{}
}

func (q BoostingQuery) Positive(queries ...Query) BoostingQuery {
	q.positiveClauses = append(q.positiveClauses, queries...)
	return q
}

func (q BoostingQuery) Negative(queries ...Query) BoostingQuery {
	q.negativeClauses = append(q.negativeClauses, queries...)
	return q
}

func (q BoostingQuery) NegativeBoost(negativeBoost float64) BoostingQuery {
	q.negativeBoost = &negativeBoost
	return q
}

func (q BoostingQuery) Boost(boost float64) BoostingQuery {
	q.boost = &boost
	return q
}

// Creates the query source for the boosting query.
func (q BoostingQuery) Source() interface{} {
	// {
	//     "boosting" : {
	//         "positive" : {
	//             "term" : {
	//                 "field1" : "value1"
	//             }
	//         },
	//         "negative" : {
	//             "term" : {
	//                 "field2" : "value2"
	//             }
	//         },
	//         "negative_boost" : 0.2
	//     }
	// }

	query := make(map[string]interface{})

	boostingClause := make(map[string]interface{})
	query["boosting"] = boostingClause

	// Negative and positive clause as well as negative boost
	// are mandatory in the Java client.

	// positive
	if len(q.positiveClauses) == 1 {
		boostingClause["positive"] = q.positiveClauses[0].Source()
	} else if len(q.positiveClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.positiveClauses {
			clauses = append(clauses, subQuery.Source())
		}
		boostingClause["positive"] = clauses
	}

	// negative
	if len(q.negativeClauses) == 1 {
		boostingClause["negative"] = q.negativeClauses[0].Source()
	} else if len(q.negativeClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.negativeClauses {
			clauses = append(clauses, subQuery.Source())
		}
		boostingClause["negative"] = clauses
	}

	if q.negativeBoost != nil {
		boostingClause["negative_boost"] = *q.negativeBoost
	}

	if q.boost != nil {
		boostingClause["boost"] = *q.boost
	}

	return query
}

