package elastic

type IntervalQueryRuleMatch struct {
	query    string
	maxGaps  *int
	ordered  *bool
	analyzer string
	useField string
	filter   *IntervalQueryRuleFilter
}

var _ IntervalQueryRule = &IntervalQueryRuleMatch{}

func NewIntervalQueryRuleMatch(query string) *IntervalQueryRuleMatch {
	return &IntervalQueryRuleMatch{query: query}
}

func (r *IntervalQueryRuleMatch) MaxGaps(maxGaps int) *IntervalQueryRuleMatch {
	r.maxGaps = &maxGaps
	return r
}
func (r *IntervalQueryRuleMatch) Ordered(ordered bool) *IntervalQueryRuleMatch {
	r.ordered = &ordered
	return r
}
func (r *IntervalQueryRuleMatch) Analyzer(analyzer string) *IntervalQueryRuleMatch {
	r.analyzer = analyzer
	return r
}
func (r *IntervalQueryRuleMatch) UseField(useField string) *IntervalQueryRuleMatch {
	r.useField = useField
	return r
}
func (r *IntervalQueryRuleMatch) Filter(filter *IntervalQueryRuleFilter) *IntervalQueryRuleMatch {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleMatch) Source() (interface{}, error) {
	source := make(map[string]interface{})

	source["query"] = r.query

	if r.ordered != nil {
		source["ordered"] = *r.ordered
	}
	if r.maxGaps != nil {
		source["max_gaps"] = *r.maxGaps
	}
	if r.analyzer != "" {
		source["analyzer"] = r.analyzer
	}
	if r.useField != "" {
		source["use_field"] = r.useField
	}
	if r.filter != nil {
		filterRuleSource, err := r.filter.Source()
		if err != nil {
			return nil, err
		}

		source["filter"] = filterRuleSource
	}

	return map[string]interface{}{"match": source}, nil
}

func (r *IntervalQueryRuleMatch) isIntervalQueryRule() bool {
	return true
}
