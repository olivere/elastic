// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A filter that matches documents using OR boolean operator
// on other queries. Can be placed within queries that accept a filter.
// For details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-or-filter.html
type NotFilter struct {
	filters    []Filter
	cache      *bool
	cacheKey   string
	filterName string
}

func NewNotFilter(filters ...Filter) NotFilter {
	f := NotFilter{}
	if len(filters) > 0 {
		f.filters = make([]Filter, 0, len(filters))
		f.filters = append(f.filters, filters...)
	} else {
		f.filters = make([]Filter, 0)
	}
	return f
}

func (f NotFilter) Add(filter Filter) NotFilter {
	f.filters = append(f.filters, filter)
	return f
}

func (f NotFilter) Cache(cache bool) NotFilter {
	f.cache = &cache
	return f
}

func (f NotFilter) CacheKey(cacheKey string) NotFilter {
	f.cacheKey = cacheKey
	return f
}

func (f NotFilter) FilterName(filterName string) NotFilter {
	f.filterName = filterName
	return f
}

func (f NotFilter) Source() interface{} {
	// {
	//   "not" : [
	//      ... filters ...
	//   ]
	// }

	source := make(map[string]interface{})

	params := make(map[string]interface{})
	source["not"] = params

	filters := make([]interface{}, len(f.filters))
	params["filters"] = filters
	for i, filter := range f.filters {
		filters[i] = filter.Source()
	}

	if f.cache != nil {
		params["_cache"] = *f.cache
	}
	if f.cacheKey != "" {
		params["_cache_key"] = f.cacheKey
	}
	if f.filterName != "" {
		params["_name"] = f.filterName
	}
	return source
}
