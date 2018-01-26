// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// CompositeAggregation is a multi-bucket values source based aggregation
// that can be used to calculate unique composite values from source documents
//
// See: https://www.elastic.co/guide/en/elasticsearch/reference/6.1/search-aggregations-bucket-composite-aggregation.html
type CompositeAggregation struct {
	size            *int
	sources         []compositeAggregationSource
	subAggregations map[string]Aggregation
	meta            map[string]interface{}
	after           map[string]interface{}
}

type compositeAggregationSource struct {
	name string
	agg  Aggregation
}

func NewCompositeAggregation() *CompositeAggregation {
	return &CompositeAggregation{
		sources:         make([]compositeAggregationSource, 0),
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

func (a *CompositeAggregation) AddSourceTerms(name string, agg *TermsAggregation) *CompositeAggregation {
	source := compositeAggregationSource{
		name: name,
		agg:  agg,
	}

	a.sources = append(a.sources, source)
	return a
}

func (a *CompositeAggregation) AddSourceHistogram(name string, agg *HistogramAggregation) *CompositeAggregation {
	source := compositeAggregationSource{
		name: name,
		agg:  agg,
	}

	a.sources = append(a.sources, source)
	return a
}

func (a *CompositeAggregation) AddSourceDateHistogram(name string, agg *DateHistogramAggregation) *CompositeAggregation {
	source := compositeAggregationSource{
		name: name,
		agg:  agg,
	}

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

		// Build the aggregation
		sourceAgg, err := s.agg.Source()
		if err != nil {
			return nil, err
		}
		sourceRecord := map[string]interface{}{
			s.name: sourceAgg,
		}
		sources = append(sources, sourceRecord)
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
