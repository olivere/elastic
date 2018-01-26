// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// CompositeAggregation is a multi-bucket values source based aggregation
// that can be used to calculate unique composite values from source documents
//
// See: https://www.co/guide/en/elasticsearch/reference/6.1/search-aggregations-bucket-composite-aggregation.html
type CompositeAggregation struct {
	size            *int
	sources         []CompositeAggregationSource
	subAggregations map[string]Aggregation
	meta            map[string]interface{}
	after           map[string]interface{}
}

// The interface describing Composite Aggregation Options
type CompositeAggregationSource interface {
	Source() (interface{}, error)
	OrderAsc(bool) CompositeAggregationSource
	Missing(string) CompositeAggregationSource
}

func NewCompositeAggregation() *CompositeAggregation {
	return &CompositeAggregation{
		sources:         make([]CompositeAggregationSource, 0),
		subAggregations: make(map[string]Aggregation),
	}
}

func (a *CompositeAggregation) Size(size int) *CompositeAggregation {
	a.size = &size
	return a
}

func (a *CompositeAggregation) After(after map[string]interface{}) *CompositeAggregation {
	a.after = after
	return a
}

func (a *CompositeAggregation) AddSource(source CompositeAggregationSource) *CompositeAggregation {
	a.sources = append(a.sources, source)
	return a
}

func (a *CompositeAggregation) SubAggregation(name string, subAggregation Aggregation) *CompositeAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *CompositeAggregation) Meta(metaData map[string]interface{}) *CompositeAggregation {
	a.meta = metaData
	return a
}

func (a *CompositeAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "aggs" : {
	//         "my_composite_agg" : {
	//             "composite" : {
	//                 "sources": [
	//				      {"my_term": { "terms": { "field": "product" }}},
	//				      {"my_histo": { "histogram": { "field": "price", "interval": 5 }}},
	//				      {"my_date": { "date_histogram": { "field": "timestamp", "interval": "1d" }}},
	//                 ],
	//                 "size" : 10,
	//                 "after" : ["a", 2, "c"]
	//             }
	//         }
	//     }
	// }
	//
	// This method returns only the { "histogram" : { ... } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["composite"] = opts

	sources := make([]interface{}, 0)
	for _, s := range a.sources {
		aggSource, err := s.Source()
		if err != nil {
			return nil, err
		}
		sources = append(sources, aggSource)
	}
	opts["sources"] = sources

	if a.size != nil {
		opts["size"] = *a.size
	}

	if a.after != nil {
		opts["after"] = a.after
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}

// CompositeAggregationSourceTerms is a source for the CompositeAggregation that handles terms
// it works very similar to a terms aggregation with slightly different syntax
type CompositeAggregationSourceTerms struct {
	name     string
	field    string
	orderAsc *bool
	missing  *string
}

func NewCompositeAggregationSourceTerms(name string, field string) *CompositeAggregationSourceTerms {
	return &CompositeAggregationSourceTerms{
		name:  name,
		field: field,
	}
}

func (a *CompositeAggregationSourceTerms) OrderAsc(orderAsc bool) CompositeAggregationSource {
	a.orderAsc = &orderAsc
	return a
}

func (a *CompositeAggregationSourceTerms) Missing(missing string) CompositeAggregationSource {
	a.missing = &missing
	return a
}

func (a *CompositeAggregationSourceTerms) Source() (interface{}, error) {

	source := make(map[string]interface{})
	name := make(map[string]interface{})
	source[a.name] = name
	terms := make(map[string]interface{})
	name["terms"] = terms

	// field
	terms["field"] = a.field
	// order
	if a.orderAsc != nil {
		if *a.orderAsc == true {
			terms["order"] = "asc"
		} else {
			terms["order"] = "desc"
		}
	}
	// missing
	if a.missing != nil {
		terms["missing"] = *a.missing
	}

	return source, nil

}

// CompositeAggregationSourceHistogram is a source for the CompositeAggregation that handles histograms
// it works very similar to a terms histogram with slightly different syntax
type CompositeAggregationSourceHistogram struct {
	name     string
	field    string
	interval int
	orderAsc *bool
	missing  *string
}

func NewCompositeAggregationSourceHistogram(name string, field string, interval int) *CompositeAggregationSourceHistogram {
	return &CompositeAggregationSourceHistogram{
		name:     name,
		field:    field,
		interval: interval,
	}
}

func (a *CompositeAggregationSourceHistogram) OrderAsc(orderAsc bool) CompositeAggregationSource {
	a.orderAsc = &orderAsc
	return a
}

func (a *CompositeAggregationSourceHistogram) Missing(missing string) CompositeAggregationSource {
	a.missing = &missing
	return a
}

func (a *CompositeAggregationSourceHistogram) Source() (interface{}, error) {

	source := make(map[string]interface{})
	name := make(map[string]interface{})
	source[a.name] = name
	histogram := make(map[string]interface{})
	name["histogram"] = histogram

	// base info
	histogram["field"] = a.field
	histogram["interval"] = a.interval
	// order
	if a.orderAsc != nil {
		if *a.orderAsc == true {
			histogram["order"] = "asc"
		} else {
			histogram["order"] = "desc"
		}
	}
	// missing
	if a.missing != nil {
		histogram["missing"] = *a.missing
	}

	return source, nil

}

// CompositeAggregationSourceDateHistogram is a source for the CompositeAggregation that handles date histograms
// it works very similar to a date histogram aggregation with slightly different syntax
type CompositeAggregationSourceDateHistogram struct {
	name     string
	field    string
	interval string
	timeZone *string
	orderAsc *bool
	missing  *string
}

func NewCompositeAggregationSourceDateHistogram(name string, field string, interval string) *CompositeAggregationSourceDateHistogram {
	return &CompositeAggregationSourceDateHistogram{
		name:     name,
		field:    field,
		interval: interval,
	}
}

func (a *CompositeAggregationSourceDateHistogram) OrderAsc(orderAsc bool) CompositeAggregationSource {
	a.orderAsc = &orderAsc
	return a
}

func (a *CompositeAggregationSourceDateHistogram) Missing(missing string) CompositeAggregationSource {
	a.missing = &missing
	return a
}

func (a *CompositeAggregationSourceDateHistogram) TimeZone(timeZone string) *CompositeAggregationSourceDateHistogram {
	a.timeZone = &timeZone
	return a
}

func (a *CompositeAggregationSourceDateHistogram) Source() (interface{}, error) {

	source := make(map[string]interface{})
	name := make(map[string]interface{})
	source[a.name] = name
	dateHistogram := make(map[string]interface{})
	name["date_histogram"] = dateHistogram

	// base info
	dateHistogram["field"] = a.field
	dateHistogram["interval"] = a.interval
	// timeZone
	if a.timeZone != nil {
		dateHistogram["time_zone"] = *a.timeZone
	}
	// order
	if a.orderAsc != nil {
		if *a.orderAsc == true {
			dateHistogram["order"] = "asc"
		} else {
			dateHistogram["order"] = "desc"
		}
	}
	// missing
	if a.missing != nil {
		dateHistogram["missing"] = *a.missing
	}

	return source, nil

}
