package elastic

type IntervalQueryRuleAllOf struct {
	intervals []IntervalQueryRule
	maxGaps   *int
	ordered   *bool
	filter    *IntervalQueryRuleFilter
}

var _ IntervalQueryRule = &IntervalQueryRuleAllOf{}

func NewIntervalQueryRuleAllOf(intervals ...IntervalQueryRule) *IntervalQueryRuleAllOf {
	return &IntervalQueryRuleAllOf{intervals: intervals}
}

func (r *IntervalQueryRuleAllOf) MaxGaps(maxGaps int) *IntervalQueryRuleAllOf {
	r.maxGaps = &maxGaps
	return r
}

func (r *IntervalQueryRuleAllOf) Ordered(ordered bool) *IntervalQueryRuleAllOf {
	r.ordered = &ordered
	return r
}

func (r *IntervalQueryRuleAllOf) Filter(filter *IntervalQueryRuleFilter) *IntervalQueryRuleAllOf {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleAllOf) Source() (interface{}, error) {
	source := make(map[string]interface{})

	intervalSources := make([]interface{}, 0)
	for _, interval := range r.intervals {
		src, err := interval.Source()
		if err != nil {
			return nil, err
		}

		intervalSources = append(intervalSources, src)
	}
	source["intervals"] = intervalSources

	if r.ordered != nil {
		source["ordered"] = *r.ordered
	}
	if r.maxGaps != nil {
		source["max_gaps"] = *r.maxGaps
	}
	if r.filter != nil {
		src, err := r.filter.Source()
		if err != nil {
			return nil, err
		}

		source["filter"] = src
	}

	return map[string]interface{}{"all_of": source}, nil
}

func (r *IntervalQueryRuleAllOf) isIntervalQueryRule() bool {
	return true
}
