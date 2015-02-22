// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// ExtendedExtendedStatsAggregation is a multi-value metrics aggregation that
// computes stats over numeric values extracted from the aggregated documents.
// These values can be extracted either from specific numeric fields
// in the documents, or be generated by a provided script.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-aggregations-metrics-extendedstats-aggregation.html
type ExtendedStatsAggregation struct {
	field           string
	script          string
	lang            string
	format          string
	params          map[string]interface{}
	subAggregations map[string]Aggregation
}

func NewExtendedStatsAggregation() ExtendedStatsAggregation {
	a := ExtendedStatsAggregation{
		params:          make(map[string]interface{}),
		subAggregations: make(map[string]Aggregation),
	}
	return a
}

func (a ExtendedStatsAggregation) Field(field string) ExtendedStatsAggregation {
	a.field = field
	return a
}

func (a ExtendedStatsAggregation) Script(script string) ExtendedStatsAggregation {
	a.script = script
	return a
}

func (a ExtendedStatsAggregation) Lang(lang string) ExtendedStatsAggregation {
	a.lang = lang
	return a
}

func (a ExtendedStatsAggregation) Format(format string) ExtendedStatsAggregation {
	a.format = format
	return a
}

func (a ExtendedStatsAggregation) Param(name string, value interface{}) ExtendedStatsAggregation {
	a.params[name] = value
	return a
}

func (a ExtendedStatsAggregation) SubAggregation(name string, subAggregation Aggregation) ExtendedStatsAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

func (a ExtendedStatsAggregation) Source() interface{} {
	// Example:
	//	{
	//    "aggs" : {
	//      "grades_stats" : { "extended_stats" : { "field" : "grade" } }
	//    }
	//	}
	// This method returns only the { "extended_stats" : { "field" : "grade" } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["extended_stats"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != "" {
		opts["script"] = a.script
	}
	if a.lang != "" {
		opts["lang"] = a.lang
	}
	if a.format != "" {
		opts["format"] = a.format
	}
	if len(a.params) > 0 {
		opts["params"] = a.params
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
