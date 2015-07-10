// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A filter that matches documents matching boolean combinations
// of other queries. Similar in concept to Boolean query,
// except that the clauses are other filters.
// Can be placed within queries that accept a filter.
// For more details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-bool-filter.html
type BoolFilter struct {
	mustClauses    []Filter
	shouldClauses  []Filter
	mustNotClauses []Filter
	filterName     string
}

// NewBoolFilter creates a new bool filter.
func NewBoolFilter() BoolFilter {
	f := BoolFilter{
		mustClauses:    make([]Filter, 0),
		shouldClauses:  make([]Filter, 0),
		mustNotClauses: make([]Filter, 0),
	}
	return f
}

func (f BoolFilter) Must(filters ...Filter) BoolFilter {
	f.mustClauses = append(f.mustClauses, filters...)
	return f
}

func (f BoolFilter) MustNot(filters ...Filter) BoolFilter {
	f.mustNotClauses = append(f.mustNotClauses, filters...)
	return f
}

func (f BoolFilter) Should(filters ...Filter) BoolFilter {
	f.shouldClauses = append(f.shouldClauses, filters...)
	return f
}

func (f BoolFilter) FilterName(filterName string) BoolFilter {
	f.filterName = filterName
	return f
}

// Creates the query source for the bool query.
func (f BoolFilter) Source() (interface{}, error) {
	// {
	//	"bool" : {
	//		"must" : {
	//			"term" : { "user" : "kimchy" }
	//		},
	//		"must_not" : {
	//			"range" : {
	//				"age" : { "from" : 10, "to" : 20 }
	//			}
	//		},
	//		"should" : [
	//			{
	//				"term" : { "tag" : "wow" }
	//			},
	//			{
	//				"term" : { "tag" : "elasticsearch" }
	//			}
	//		]
	//	}
	// }

	source := make(map[string]interface{})

	boolClause := make(map[string]interface{})
	source["bool"] = boolClause

	// must
	if len(f.mustClauses) == 1 {
		src, err := f.mustClauses[0].Source()
		if err != nil {
			return nil, err
		}
		boolClause["must"] = src
	} else if len(f.mustClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range f.mustClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["must"] = clauses
	}

	// must_not
	if len(f.mustNotClauses) == 1 {
		src, err := f.mustNotClauses[0].Source()
		if err != nil {
			return nil, err
		}
		boolClause["must_not"] = src
	} else if len(f.mustNotClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range f.mustNotClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["must_not"] = clauses
	}

	// should
	if len(f.shouldClauses) == 1 {
		src, err := f.shouldClauses[0].Source()
		if err != nil {
			return nil, err
		}
		boolClause["should"] = src
	} else if len(f.shouldClauses) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range f.shouldClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["should"] = clauses
	}

	if f.filterName != "" {
		boolClause["_name"] = f.filterName
	}

	return source, nil
}
