// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Query Facet
// See: http://www.elasticsearch.org/guide/reference/api/search/facets/query-facet.html
type QueryFacet struct {
	Facet
	global   *bool
	scope    string
	query    Query
	nested   string
	filters  []Filter
	size     *int
	order    string
}

func NewQueryFacet(query Query) QueryFacet {
	f := QueryFacet{
		query: query,
		filters: make([]Filter, 0),
	}
	return f
}

func (f QueryFacet) Global(global bool) QueryFacet {
	f.global = &global
	return f
}

func (f QueryFacet) Size(size int) QueryFacet {
	f.size = &size
	return f
}

// Valid order options are: "count" (default), "term",
// "reverse_count", and "reverse_term".
func (f QueryFacet) Order(order string) QueryFacet {
	f.order = order
	return f
}

func (f QueryFacet) Scope(scope string) QueryFacet {
	f.scope = scope
	return f
}

func (f QueryFacet) Nested(nested string) QueryFacet {
	f.nested = nested
	return f
}

func (f QueryFacet) Query(query Query) QueryFacet {
	f.query = query
	return f
}

func (f QueryFacet) Filter(filter Filter) QueryFacet {
	f.filters = append(f.filters, filter)
	return f
}

func (f QueryFacet) Source() interface{} {
	source := make(map[string]interface{})
	source["query"] = f.query.Source()

	if f.global != nil {
		source["global"] = *f.global
	}

	if f.size != nil {
		source["size"] = *f.size
	}

	if f.order != "" {
		source["order"] = f.order
	}

	if f.nested != "" {
		source["nested"] = f.nested
	}

	if f.scope != "" {
		source["scope"] = f.scope
	}

	if len(f.filters) == 1 {
		source["facet_filter"] = f.filters[0].Source()
	} else if len(f.filters) > 1 {
		ff := make(map[string]interface{})
		andedFilters := make([]interface{}, 0)
		for _, filter := range f.filters {
			andedFilters = append(andedFilters, filter.Source())
		}
		ff["and"] = andedFilters
		source["facet_filter"] = ff
	}

	return source
}
