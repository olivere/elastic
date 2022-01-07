package elastic

var (
	_ IntervalQueryRule = (*IntervalQueryRuleFuzzy)(nil)
)

// IntervalQueryRuleFuzzy is an implementation of IntervalQueryRule.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.16/query-dsl-intervals-query.html#intervals-fuzzy
// for details.
type IntervalQueryRuleFuzzy struct {
	term           string
	prefixLength   *int
	transpositions *bool
	fuzziness      interface{}
	analyzer       string
	useField       string
}

// NewIntervalQueryRuleFuzzy initializes and returns a new instance
// of IntervalQueryRuleFuzzy.
func NewIntervalQueryRuleFuzzy(term string) *IntervalQueryRuleFuzzy {
	return &IntervalQueryRuleFuzzy{term: term}
}

// PrefixLength is the number of beginning characters left unchanged when
// creating expansions. Defaults to 0.
func (q *IntervalQueryRuleFuzzy) PrefixLength(prefixLength int) *IntervalQueryRuleFuzzy {
	q.prefixLength = &prefixLength
	return q
}

// Fuzziness is the maximum edit distance allowed for matching.
// It can be integers like 0, 1 or 2 as well as strings
// like "auto", "0..1", "1..4" or "0.0..1.0". Defaults to "auto".
func (q *IntervalQueryRuleFuzzy) Fuzziness(fuzziness interface{}) *IntervalQueryRuleFuzzy {
	q.fuzziness = fuzziness
	return q
}

// Transpositions indicates whether edits include transpositions of two
// adjacent characters (ab -> ba). Defaults to true.
func (q *IntervalQueryRuleFuzzy) Transpositions(transpositions bool) *IntervalQueryRuleFuzzy {
	q.transpositions = &transpositions
	return q
}

// Analyzer specifies the analyzer used to analyze terms in the query.
func (r *IntervalQueryRuleFuzzy) Analyzer(analyzer string) *IntervalQueryRuleFuzzy {
	r.analyzer = analyzer
	return r
}

// UseField, if specified, matches the intervals from this field rather than
// the top-level field.
func (r *IntervalQueryRuleFuzzy) UseField(useField string) *IntervalQueryRuleFuzzy {
	r.useField = useField
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleFuzzy) Source() (interface{}, error) {
	source := make(map[string]interface{})

	source["term"] = r.term

	if r.prefixLength != nil {
		source["prefix_length"] = *r.prefixLength
	}
	if r.transpositions != nil {
		source["transpositions"] = *r.transpositions
	}
	if r.fuzziness != "" {
		source["fuzziness"] = r.fuzziness
	}
	if r.analyzer != "" {
		source["analyzer"] = r.analyzer
	}
	if r.useField != "" {
		source["use_field"] = r.useField
	}

	return map[string]interface{}{
		"fuzzy": source,
	}, nil
}

// isIntervalQueryRule implements the marker interface.
func (r *IntervalQueryRuleFuzzy) isIntervalQueryRule() bool {
	return true
}
