// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestRankFeatureQueryTest(t *testing.T) {
	tests := []struct {
		Query    Query
		Expected string
	}{
		// #0
		{
			Query:    NewRankFeatureQuery("pagerank"),
			Expected: `{"rank_feature":{"field":"pagerank"}}`,
		},
		// #1
		{
			Query:    NewRankFeatureQuery("url_length").Boost(0.1),
			Expected: `{"rank_feature":{"boost":0.1,"field":"url_length"}}`,
		},
		// #2
		{
			Query:    NewRankFeatureQuery("pagerank").ScoreFunction(NewRankFeatureSaturationScoreFunction().Pivot(8)),
			Expected: `{"rank_feature":{"field":"pagerank","saturation":{"pivot":8}}}`,
		},
		// #3
		{
			Query:    NewRankFeatureQuery("pagerank").ScoreFunction(NewRankFeatureSaturationScoreFunction()),
			Expected: `{"rank_feature":{"field":"pagerank","saturation":{}}}`,
		},
		// #4
		{
			Query:    NewRankFeatureQuery("pagerank").ScoreFunction(NewRankFeatureLogScoreFunction(4)),
			Expected: `{"rank_feature":{"field":"pagerank","log":{"scaling_factor":4}}}`,
		},
		// #5
		{
			Query:    NewRankFeatureQuery("pagerank").ScoreFunction(NewRankFeatureSigmoidScoreFunction(7, 0.6)),
			Expected: `{"rank_feature":{"field":"pagerank","sigmoid":{"exponent":0.6,"pivot":7}}}`,
		},
		// #6
		{
			Query:    NewRankFeatureQuery("pagerank").ScoreFunction(NewRankFeatureLinearScoreFunction()),
			Expected: `{"rank_feature":{"field":"pagerank","linear":{}}}`,
		},
	}

	for i, tt := range tests {
		src, err := tt.Query.Source()
		if err != nil {
			t.Fatalf("#%d: encoding Source failed: %v", i, err)
		}
		data, err := json.Marshal(src)
		if err != nil {
			t.Fatalf("#%d: marshaling to JSON failed: %v", i, err)
		}
		if want, got := tt.Expected, string(data); want != got {
			t.Fatalf("#%d: expected\n%s\ngot:\n%s", i, want, got)
		}
	}
}
