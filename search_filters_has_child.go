// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// The has_child query works the same as the has_child filter,
// by automatically wrapping the filter with a constant_score
// (when using the default score type).
// For more details, see
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-has-child-query.html
type HasChildFilter struct {
	filter             Filter
	query              Query
	childType          string
	filterName         string
	shortCircuitCutoff *int
	minChildren        *int
	maxChildren        *int
	innerHit           *InnerHit
}

// NewHasChildFilter creates a new has_child query.
func NewHasChildFilter(childType string) HasChildFilter {
	f := HasChildFilter{
		childType: childType,
	}
	return f
}

func (f HasChildFilter) Query(query Query) HasChildFilter {
	f.query = query
	return f
}

func (f HasChildFilter) Filter(filter Filter) HasChildFilter {
	f.filter = filter
	return f
}

func (f HasChildFilter) FilterName(filterName string) HasChildFilter {
	f.filterName = filterName
	return f
}

func (f HasChildFilter) ShortCircuitCutoff(shortCircuitCutoff int) HasChildFilter {
	f.shortCircuitCutoff = &shortCircuitCutoff
	return f
}

func (f HasChildFilter) MinChildren(minChildren int) HasChildFilter {
	f.minChildren = &minChildren
	return f
}

func (f HasChildFilter) MaxChildren(maxChildren int) HasChildFilter {
	f.maxChildren = &maxChildren
	return f
}

func (f HasChildFilter) InnerHit(innerHit *InnerHit) HasChildFilter {
	f.innerHit = innerHit
	return f
}

// Source returns the JSON document for the filter.
func (f HasChildFilter) Source() (interface{}, error) {
	// {
	//   "has_child" : {
	//       "type" : "blog_tag",
	//       "query" : {
	//           "term" : {
	//               "tag" : "something"
	//           }
	//       }
	//   }
	// }

	source := make(map[string]interface{})

	filter := make(map[string]interface{})
	source["has_child"] = filter

	if f.query != nil {
		src, err := f.query.Source()
		if err != nil {
			return nil, err
		}
		filter["query"] = src
	} else if f.filter != nil {
		src, err := f.filter.Source()
		if err != nil {
			return nil, err
		}
		filter["filter"] = src
	}

	filter["type"] = f.childType
	if f.filterName != "" {
		filter["_name"] = f.filterName
	}
	if f.shortCircuitCutoff != nil {
		filter["short_circuit_cutoff"] = *f.shortCircuitCutoff
	}
	if f.minChildren != nil {
		filter["min_children"] = *f.minChildren
	}
	if f.maxChildren != nil {
		filter["max_children"] = *f.maxChildren
	}
	if f.innerHit != nil {
		src, err := f.innerHit.Source()
		if err != nil {
			return nil, err
		}
		filter["inner_hits"] = src
	}
	return source, nil
}
