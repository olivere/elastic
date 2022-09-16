// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// VariableWidthHistogramAggregation is a multi-bucket values source based aggregation
// that can be applied on numeric values extracted from the documents.
// It dynamically builds buckets over the values given a target number of buckets.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/search-aggregations-bucket-variablewidthhistogram-aggregation.html
type VariableWidthHistogramAggregation struct {
	field           string
	subAggregations map[string]Aggregation
	meta            map[string]interface{}

	buckets       *int64
	initialBuffer *int64
	shardSize     *int64
}

func NewVariableWidthHistogramAggregation() *VariableWidthHistogramAggregation {
	return &VariableWidthHistogramAggregation{
		subAggregations: make(map[string]Aggregation),
	}
}

func (a *VariableWidthHistogramAggregation) Field(field string) *VariableWidthHistogramAggregation {
	a.field = field
	return a
}

func (a *VariableWidthHistogramAggregation) Buckets(buckets int64) *VariableWidthHistogramAggregation {
	a.buckets = &buckets
	return a
}

func (a *VariableWidthHistogramAggregation) InitialBuffer(initialBuffer int64) *VariableWidthHistogramAggregation {
	a.initialBuffer = &initialBuffer
	return a
}

func (a *VariableWidthHistogramAggregation) ShardSize(shardSize int64) *VariableWidthHistogramAggregation {
	a.shardSize = &shardSize
	return a
}

func (a *VariableWidthHistogramAggregation) SubAggregation(name string, subAggregation Aggregation) *VariableWidthHistogramAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *VariableWidthHistogramAggregation) Meta(metaData map[string]interface{}) *VariableWidthHistogramAggregation {
	a.meta = metaData
	return a
}

func (a *VariableWidthHistogramAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "aggs" : {
	//         "prices" : {
	//             "variable_width_histogram" : {
	//                 "field" : "price",
	//                 "buckets" : 2
	//             }
	//         }
	//     }
	// }
	//
	// This method returns only the { "variable_width_histogram" : { ... } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["variable_width_histogram"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}

	if a.buckets != nil {
		opts["buckets"] = *a.buckets
	}
	if a.initialBuffer != nil {
		opts["initial_buffer"] = *a.initialBuffer
	}
	if a.shardSize != nil {
		opts["shard_size"] = *a.shardSize
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
