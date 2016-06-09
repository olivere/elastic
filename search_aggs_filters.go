// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// FiltersAggregation defines a multi bucket aggregations where each bucket
// is associated with a filter. Each bucket will collect all documents that
// match its associated filter.
//
// Notice that the caller has to decide whether to add filters by name
// (using FilterWithName) or unnamed filters (using Filter or Filters). One cannot
// use both named and unnamed filters.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-filters-aggregation.html
type FiltersAggregation struct {
	unnamedFilters  []Filter
	namedFilters    map[string]Filter
	subAggregations map[string]Aggregation
}

// NewFiltersAggregation initializes a new FiltersAggregation.
func NewFiltersAggregation() FiltersAggregation {
	return FiltersAggregation{
		unnamedFilters:  make([]Filter, 0),
		namedFilters:    make(map[string]Filter),
		subAggregations: make(map[string]Aggregation),
	}
}

// Filter adds an unnamed filter. Notice that you can
// either use named or unnamed filters, but not both.
func (a FiltersAggregation) Filter(filter Filter) FiltersAggregation {
	a.unnamedFilters = append(a.unnamedFilters, filter)
	return a
}

// Filters adds one or more unnamed filters. Notice that you can
// either use named or unnamed filters, but not both.
func (a FiltersAggregation) Filters(filters ...Filter) FiltersAggregation {
	if len(filters) > 0 {
		a.unnamedFilters = append(a.unnamedFilters, filters...)
	}
	return a
}

// FilterWithName adds a filter with a specific name. Notice that you can
// either use named or unnamed filters, but not both.
func (a FiltersAggregation) FilterWithName(name string, filter Filter) FiltersAggregation {
	a.namedFilters[name] = filter
	return a
}

// SubAggregation adds a sub-aggregation to this aggregation.
func (a FiltersAggregation) SubAggregation(name string, subAggregation Aggregation) FiltersAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Source returns the a JSON-serializable interface.
func (a FiltersAggregation) Source() interface{} {
	// Example:
	//	{
	//  "aggs" : {
	//    "messages" : {
	//      "filters" : {
	//        "filters" : {
	//          "errors" :   { "term" : { "body" : "error"   }},
	//          "warnings" : { "term" : { "body" : "warning" }}
	//        }
	//      }
	//    }
	//  }
	//	}
	// This method returns only the (outer) { "filters" : {} } part.

	source := make(map[string]interface{})
	filters := make(map[string]interface{})
	source["filters"] = filters

	// TODO We cannot return an error here due to compatibility reasons
	// in elastic.v2, so we're going to return only the unnamed filters for now.

	if len(a.unnamedFilters) > 0 {
		arr := make([]interface{}, len(a.unnamedFilters))
		for i, filter := range a.unnamedFilters {
			arr[i] = filter.Source()
		}
		filters["filters"] = arr
	} else {
		dict := make(map[string]interface{})
		for key, filter := range a.namedFilters {
			dict[key] = filter.Source()
		}
		filters["filters"] = dict
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			aggsMap[name] = aggregate.Source()
		}
	}

	return source
}
