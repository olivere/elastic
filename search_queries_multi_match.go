// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
)

// The multi_match query builds further on top of the match query by allowing multiple fields to be specified.
// For more details, see:
// http://www.elasticsearch.org/guide/reference/query-dsl/multi-match-query.html
type MultiMatchQuery struct {
	Query
	text               interface{}
	fields             []string
	fieldBoosts        map[string]*float32
	_type              *MatchQueryType
	operator           string
	analyzer           string
	boost              *float32
	slop               *int
	fuzziness          string
	prefixLength       *int
	maxExpansions      *int
	minimumShouldMatch string
	rewrite            string
	fuzzyRewrite       string
	useDisMax          *bool
	tieBreaker         *float32
	lenient            *bool
}

func NewMultiMatchQuery(text interface{}, fields ...string) MultiMatchQuery {
	q := MultiMatchQuery{
		text:        text,
		fields:      make([]string, 0),
		fieldBoosts: make(map[string]*float32),
	}
	q.fields = append(q.fields, fields...)
	return q
}

func (q MultiMatchQuery) Field(field string) MultiMatchQuery {
	q.fields = append(q.fields, field)
	return q
}

func (q MultiMatchQuery) FieldWithBoost(field string, boost float32) MultiMatchQuery {
	q.fields = append(q.fields, field)
	q.fieldBoosts[field] = &boost
	return q
}

func (q MultiMatchQuery) Type(_type MatchQueryType) MultiMatchQuery {
	q._type = &_type
	return q
}

func (q MultiMatchQuery) Operator(operator string) MultiMatchQuery {
	q.operator = operator
	return q
}

func (q MultiMatchQuery) Analyzer(analyzer string) MultiMatchQuery {
	q.analyzer = analyzer
	return q
}

func (q MultiMatchQuery) Boost(boost float32) MultiMatchQuery {
	q.boost = &boost
	return q
}

func (q MultiMatchQuery) Slop(slop int) MultiMatchQuery {
	q.slop = &slop
	return q
}

func (q MultiMatchQuery) Fuzziness(fuzziness string) MultiMatchQuery {
	q.fuzziness = fuzziness
	return q
}

func (q MultiMatchQuery) PrefixLength(prefixLength int) MultiMatchQuery {
	q.prefixLength = &prefixLength
	return q
}

func (q MultiMatchQuery) MaxExpansions(maxExpansions int) MultiMatchQuery {
	q.maxExpansions = &maxExpansions
	return q
}

func (q MultiMatchQuery) MinimumShouldMatch(minimumShouldMatch string) MultiMatchQuery {
	q.minimumShouldMatch = minimumShouldMatch
	return q
}

func (q MultiMatchQuery) Rewrite(rewrite string) MultiMatchQuery {
	q.rewrite = rewrite
	return q
}

func (q MultiMatchQuery) FuzzyRewrite(fuzzyRewrite string) MultiMatchQuery {
	q.fuzzyRewrite = fuzzyRewrite
	return q
}

func (q MultiMatchQuery) UseDisMax(useDisMax bool) MultiMatchQuery {
	q.useDisMax = &useDisMax
	return q
}

func (q MultiMatchQuery) TieBreaker(tieBreaker float32) MultiMatchQuery {
	q.tieBreaker = &tieBreaker
	return q
}

func (q MultiMatchQuery) Lenient(lenient bool) MultiMatchQuery {
	q.lenient = &lenient
	return q
}

func (q MultiMatchQuery) Source() interface{} {
	//
	// {
	//   "multi_match" : {
	//     "query" : "this is a test",
	//     "fields" : [ "subject", "message" ]
	//   }
	// }

	source := make(map[string]interface{})

	multiMatch := make(map[string]interface{})
	source["multi_match"] = multiMatch

	multiMatch["query"] = q.text

	if len(q.fields) > 0 {
		fields := make([]string, 0)
		for _, field := range q.fields {
			if boost, found := q.fieldBoosts[field]; found {
				if boost != nil {
					fields = append(fields, fmt.Sprintf("%s^%f", field, *boost))
				} else {
					fields = append(fields, field)
				}
			} else {
				fields = append(fields, field)
			}
		}
		multiMatch["fields"] = fields
	}

	if q._type != nil {
		if *q._type == Boolean {
			multiMatch["type"] = "boolean"
		} else if *q._type == Phrase {
			multiMatch["type"] = "phrase"
		} else if *q._type == PhrasePrefix {
			multiMatch["type"] = "phrase_prefix"
		}
	}

	if q.operator != "" {
		multiMatch["operator"] = q.operator
	}

	if q.analyzer != "" {
		multiMatch["analyzer"] = q.analyzer
	}

	if q.boost != nil {
		multiMatch["boost"] = *q.boost
	}

	if q.slop != nil {
		multiMatch["slop"] = *q.slop
	}

	if q.fuzziness != "" {
		multiMatch["fuzziness"] = q.fuzziness
	}

	if q.prefixLength != nil {
		multiMatch["prefix_length"] = *q.prefixLength
	}

	if q.maxExpansions != nil {
		multiMatch["max_expansions"] = *q.maxExpansions
	}

	if q.minimumShouldMatch != "" {
		multiMatch["minimum_should_match"] = q.minimumShouldMatch
	}

	if q.rewrite != "" {
		multiMatch["rewrite"] = q.rewrite
	}

	if q.fuzzyRewrite != "" {
		multiMatch["fuzzy_rewrite"] = q.fuzzyRewrite
	}

	if q.useDisMax != nil {
		multiMatch["use_dis_max"] = *q.useDisMax
	}

	if q.tieBreaker != nil {
		multiMatch["tie_breaker"] = *q.tieBreaker
	}

	if q.lenient != nil {
		multiMatch["lenient"] = *q.lenient
	}

	return source
}
