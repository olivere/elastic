// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// RankFeatureQuery boosts the relevance score of documents based on the
// numeric value of a rank_feature or rank_features field.
//
// The RankFeatureQuery is typically used in the should clause of a BoolQuery
// so its relevance scores are added to other scores from the BoolQuery.
//
// For more details, see:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.14/query-dsl-rank-feature-query.html
type RankFeatureQuery struct {
	field     string
	scoreFunc RankFeatureScoreFunction
	boost     *float64
	queryName string
}

// NewRankFeatureQuery creates and initializes a new RankFeatureQuery.
func NewRankFeatureQuery(field string) *RankFeatureQuery {
	return &RankFeatureQuery{
		field: field,
	}
}

// Field name.
func (q *RankFeatureQuery) Field(field string) *RankFeatureQuery {
	q.field = field
	return q
}

// ScoreFunction specifies the score function for the RankFeatureQuery.
func (q *RankFeatureQuery) ScoreFunction(f RankFeatureScoreFunction) *RankFeatureQuery {
	q.scoreFunc = f
	return q
}

// Boost sets the boost for this query.
func (q *RankFeatureQuery) Boost(boost float64) *RankFeatureQuery {
	q.boost = &boost
	return q
}

// QueryName sets the query name for the filter that can be used when
// searching for matched_filters per hit.
func (q *RankFeatureQuery) QueryName(queryName string) *RankFeatureQuery {
	q.queryName = queryName
	return q
}

// Source returns the JSON serializable content for this query.
func (q *RankFeatureQuery) Source() (interface{}, error) {
	// {
	// 	  "rank_feature": {
	// 	  	"field": "pagerank",
	// 	  	"saturation": {
	// 		  "pivot": 8
	// 		}
	//    }
	// }

	query := make(map[string]interface{})
	params := make(map[string]interface{})
	query["rank_feature"] = params
	params["field"] = q.field
	if q.scoreFunc != nil {
		src, err := q.scoreFunc.Source()
		if err != nil {
			return nil, err
		}
		params[q.scoreFunc.Name()] = src
	}
	if q.boost != nil {
		params["boost"] = *q.boost
	}
	if q.queryName != "" {
		params["_name"] = q.queryName
	}

	return query, nil
}

// -- Score functions --

// RankFeatureScoreFunction specifies the interface for score functions
// in the context of a RankFeatureQuery.
type RankFeatureScoreFunction interface {
	Name() string
	Source() (interface{}, error)
}

// -- Log score function --

// RankFeatureLogScoreFunction represents a Logarithmic score function for a
// RankFeatureQuery.
//
// See here for details:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.14/query-dsl-rank-feature-query.html#rank-feature-query-logarithm
type RankFeatureLogScoreFunction struct {
	scalingFactor float64
}

// NewRankFeatureLogScoreFunction returns a new RankFeatureLogScoreFunction
// with the given scaling factor.
func NewRankFeatureLogScoreFunction(scalingFactor float64) *RankFeatureLogScoreFunction {
	return &RankFeatureLogScoreFunction{
		scalingFactor: scalingFactor,
	}
}

// Name of the score function.
func (f *RankFeatureLogScoreFunction) Name() string { return "log" }

// Source returns a serializable JSON object for building the query.
func (f *RankFeatureLogScoreFunction) Source() (interface{}, error) {
	return map[string]interface{}{
		"scaling_factor": f.scalingFactor,
	}, nil
}

// -- Saturation score function --

// RankFeatureSaturationScoreFunction represents a Log score function for a
// RankFeatureQuery.
//
// See here for details:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.14/query-dsl-rank-feature-query.html#rank-feature-query-saturation
type RankFeatureSaturationScoreFunction struct {
	pivot *float64
}

// NewRankFeatureSaturationScoreFunction initializes a new
// RankFeatureSaturationScoreFunction.
func NewRankFeatureSaturationScoreFunction() *RankFeatureSaturationScoreFunction {
	return &RankFeatureSaturationScoreFunction{}
}

// Pivot specifies the pivot to use.
func (f *RankFeatureSaturationScoreFunction) Pivot(pivot float64) *RankFeatureSaturationScoreFunction {
	f.pivot = &pivot
	return f
}

// Name of the score function.
func (f *RankFeatureSaturationScoreFunction) Name() string { return "saturation" }

// Source returns a serializable JSON object for building the query.
func (f *RankFeatureSaturationScoreFunction) Source() (interface{}, error) {
	m := make(map[string]interface{})
	if f.pivot != nil {
		m["pivot"] = *f.pivot
	}
	return m, nil
}

// -- Sigmoid score function --

// RankFeatureSigmoidScoreFunction represents a Sigmoid score function for a
// RankFeatureQuery.
//
// See here for details:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.14/query-dsl-rank-feature-query.html#rank-feature-query-sigmoid
type RankFeatureSigmoidScoreFunction struct {
	pivot    float64
	exponent float64
}

// NewRankFeatureSigmoidScoreFunction returns a new RankFeatureSigmoidScoreFunction
// with the given scaling factor.
func NewRankFeatureSigmoidScoreFunction(pivot, exponent float64) *RankFeatureSigmoidScoreFunction {
	return &RankFeatureSigmoidScoreFunction{
		pivot:    pivot,
		exponent: exponent,
	}
}

// Name of the score function.
func (f *RankFeatureSigmoidScoreFunction) Name() string { return "sigmoid" }

// Source returns a serializable JSON object for building the query.
func (f *RankFeatureSigmoidScoreFunction) Source() (interface{}, error) {
	return map[string]interface{}{
		"pivot":    f.pivot,
		"exponent": f.exponent,
	}, nil
}

// -- Linear score function --

// RankFeatureLinearScoreFunction represents a Linear score function for a
// RankFeatureQuery.
//
// See here for details:
// https://www.elastic.co/guide/en/elasticsearch/reference/7.14/query-dsl-rank-feature-query.html#rank-feature-query-linear
type RankFeatureLinearScoreFunction struct {
}

// NewRankFeatureLinearScoreFunction initializes a new
// RankFeatureLinearScoreFunction.
func NewRankFeatureLinearScoreFunction() *RankFeatureLinearScoreFunction {
	return &RankFeatureLinearScoreFunction{}
}

// Name of the score function.
func (f *RankFeatureLinearScoreFunction) Name() string { return "linear" }

// Source returns a serializable JSON object for building the query.
func (f *RankFeatureLinearScoreFunction) Source() (interface{}, error) {
	return map[string]interface{}{}, nil
}
