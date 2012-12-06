// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Match query types
type MatchQueryType int

const (
	Boolean MatchQueryType = iota // 0
	Phrase
	PhrasePrefix
)

// Zero terms query
type ZeroTermsQuery int

const (
	None ZeroTermsQuery = iota // 0
	All
)

// Match query is a family of match queries that
// accept text/numerics/dates, analyzes it, and
// constructs a query out of it. For more details,
// see http://www.elasticsearch.org/guide/reference/query-dsl/match-query.html
type MatchQuery struct {
	Query
	name                string
	value               interface{}
	_type               *MatchQueryType
	operator            *Operator
	analyzer            string
	boost               *float32
	slop                *int
	fuzziness           string
	prefixLength        *int
	maxExpansions       *int
	minimumShouldMatch  string
	rewrite             string
	fuzzyRewrite        string
	lenient             *bool
	fuzzyTranspositions *bool
	zeroTermsQuery      *ZeroTermsQuery
}

func NewMatchQuery(name string, value interface{}) MatchQuery {
	q := MatchQuery{name: name, value: value}
	return q
}

func (q MatchQuery) Type(_type MatchQueryType) MatchQuery {
	q._type = &_type
	return q
}

func (q MatchQuery) Operator(operator Operator) MatchQuery {
	q.operator = &operator
	return q
}

func (q MatchQuery) Analyzer(analyzer string) MatchQuery {
	q.analyzer = analyzer
	return q
}

func (q MatchQuery) Boost(boost float32) MatchQuery {
	q.boost = &boost
	return q
}

func (q MatchQuery) Slop(slop int) MatchQuery {
	q.slop = &slop
	return q
}

func (q MatchQuery) Fuzziness(fuzziness string) MatchQuery {
	q.fuzziness = fuzziness
	return q
}

func (q MatchQuery) PrefixLength(prefixLength int) MatchQuery {
	q.prefixLength = &prefixLength
	return q
}

func (q MatchQuery) MaxExpansions(maxExpansions int) MatchQuery {
	q.maxExpansions = &maxExpansions
	return q
}

func (q MatchQuery) MinimumShouldMatch(minimumShouldMatch string) MatchQuery {
	q.minimumShouldMatch = minimumShouldMatch
	return q
}

func (q MatchQuery) Rewrite(rewrite string) MatchQuery {
	q.rewrite = rewrite
	return q
}

func (q MatchQuery) FuzzyRewrite(fuzzyRewrite string) MatchQuery {
	q.fuzzyRewrite = fuzzyRewrite
	return q
}

func (q MatchQuery) Lenient(lenient bool) MatchQuery {
	q.lenient = &lenient
	return q
}

func (q MatchQuery) FuzzyTranspositions(fuzzyTranspositions bool) MatchQuery {
	q.fuzzyTranspositions = &fuzzyTranspositions
	return q
}

func (q MatchQuery) ZeroTermsQuery(zeroTermsQuery ZeroTermsQuery) MatchQuery {
	q.zeroTermsQuery = &zeroTermsQuery
	return q
}

func (q MatchQuery) Source() interface{} {
	// {"match":{"name":{"query":"value","type":"boolean/phrase"}}}
	source := make(map[string]interface{})

	match := make(map[string]interface{})
	source["match"] = match

	query := make(map[string]interface{})
	match[q.name] = query

	query["query"] = q.value

	if q._type != nil {
		if *q._type == Boolean {
			query["type"] = "boolean"
		} else if *q._type == Phrase {
			query["type"] = "phrase"
		} else if *q._type == PhrasePrefix {
			query["type"] = "phrase_prefix"
		}
	}

	if q.operator != nil {
		if *q.operator == And {
			query["operator"] = "and"
		} else if *q.operator == Or {
			query["operator"] = "or"
		}
	}

	if q.boost != nil {
		query["boost"] = *q.boost
	}

	if q.slop != nil {
		query["slop"] = *q.slop
	}

	if q.fuzziness != "" {
		query["fuzziness"] = q.fuzziness
	}

	if q.prefixLength != nil {
		query["prefix_length"] = *q.prefixLength
	}

	if q.maxExpansions != nil {
		query["max_expansions"] = *q.maxExpansions
	}

	if q.minimumShouldMatch != "" {
		query["minimum_should_match"] = q.minimumShouldMatch
	}

	if q.rewrite != "" {
		query["rewrite"] = q.rewrite
	}

	if q.fuzzyRewrite != "" {
		query["fuzzy_rewrite"] = q.fuzzyRewrite
	}


	if q.lenient != nil {
		query["lenient"] = *q.lenient
	}

	if q.fuzzyTranspositions != nil {
		query["fuzzy_transpositions"] = *q.fuzzyTranspositions
	}

	if q.zeroTermsQuery != nil {
		if *q.zeroTermsQuery == None {
			query["zero_terms_query"] = "none"
		} else if *q.zeroTermsQuery == All {
			query["zero_terms_query"] = "all"
		}
	}

	return source
}
