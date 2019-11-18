package elastic

type MatchRule struct {
	query    string
	maxGaps  *int
	ordered  *bool
	analyzer string
	useField string
	filter   *FilterRule
}

var _ IntervalRule = &MatchRule{}

func NewMatchRule(query string) *MatchRule {
	return &MatchRule{query: query}
}

func (r *MatchRule) MaxGaps(maxGaps int) *MatchRule {
	r.maxGaps = &maxGaps
	return r
}
func (r *MatchRule) Ordered(ordered bool) *MatchRule {
	r.ordered = &ordered
	return r
}
func (r *MatchRule) Analyzer(analyzer string) *MatchRule {
	r.analyzer = analyzer
	return r
}
func (r *MatchRule) UseField(useField string) *MatchRule {
	r.useField = useField
	return r
}
func (r *MatchRule) Filter(filter *FilterRule) *MatchRule {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *MatchRule) Source() (interface{}, error) {
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

func (r *MatchRule) IsIntervalRule() bool {
	return true
}
