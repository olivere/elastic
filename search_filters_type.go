// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Filters documents matching the provided document / mapping type.
// For details, see:
// http://www.elasticsearch.org/guide/reference/query-dsl/type-filter.html
type TypeFilter struct {
	Filter
	_type string
}

func NewTypeFilter(_type string) TypeFilter {
	f := TypeFilter{_type: _type}
	return f
}

func (f TypeFilter) Source() interface{} {
	// {
	//   "type" : {
	//     "value" : "..."
	//   }
	// }

	source := make(map[string]interface{})

	params := make(map[string]interface{})
	source["type"] = params

	params["value"] = f._type

	return source
}
