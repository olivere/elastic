// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A filter that filters out matched documents using a query. Can be placed
// within queries that accept a filter.
// For details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-not-filter.html#query-dsl-not-filter.
type NotFilter struct {
	filter     Filter
	filterName string
}

func NewNotFilter(filter Filter) NotFilter {
	return NotFilter{
		filter: filter,
	}
}

func (f NotFilter) FilterName(filterName string) NotFilter {
	f.filterName = filterName
	return f
}

func (f NotFilter) Source() (interface{}, error) {
	// {
	//   "not" : {
	//      "filter" : { ... }
	//   }
	// }

	source := make(map[string]interface{})

	params := make(map[string]interface{})
	source["not"] = params
	src, err := f.filter.Source()
	if err != nil {
		return nil, err
	}
	params["filter"] = src

	if f.filterName != "" {
		params["_name"] = f.filterName
	}
	return source, nil
}
