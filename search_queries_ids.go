// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Filters documents that only have the provided ids.
// For more details, see
// http://www.elasticsearch.org/guide/reference/query-dsl/ids-query.html
type IdsQuery struct {
	Query
	types       []string
	values      []string
	filterName  string
}

// Creates a new ids query.
func NewIdsQuery(types ...string) IdsQuery {
	q := IdsQuery{
		types:  types,
		values: make([]string, 0),
	}
	return q
}

func (q IdsQuery) Ids(ids ...string) IdsQuery {
	q.values = append(q.values, ids...)
	return q
}

func (q IdsQuery) FilterName(filterName string) IdsQuery {
	q.filterName = filterName
	return q
}

// Creates the query source for the ids query.
func (q IdsQuery) Source() interface{} {
	// {
	//	"ids" : {
	//		"type" : "my_type",
	//		"values" : ["1", "4", "100"]
    //	}
	// }

	source := make(map[string]interface{})

	query := make(map[string]interface{})
	source["ids"] = query

	// type(s)
	if len(q.types) == 1 {
		query["type"] = q.types[0]
	} else if len(q.types) > 1 {
		query["types"] = q.types
	}

	// values
	query["values"] = q.values

	// filter name
	if q.filterName != "" {
		query["_name"] = q.filterName
	}

	return source
}
