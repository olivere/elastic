package elastic

type IntervalQueryRuleAnyOf struct {
	intervals []IntervalQueryRule
	filter    *IntervalQueryRuleFilter
}

var _ IntervalQueryRule = &IntervalQueryRuleAnyOf{}

func NewIntervalQueryRuleAnyOf(intervals ...IntervalQueryRule) *IntervalQueryRuleAnyOf {
	return &IntervalQueryRuleAnyOf{intervals: intervals}
}
func (r *IntervalQueryRuleAnyOf) Filter(filter *IntervalQueryRuleFilter) *IntervalQueryRuleAnyOf {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleAnyOf) Source() (interface{}, error) {
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

	if r.filter != nil {
		src, err := r.filter.Source()
		if err != nil {
			return nil, err
		}

		source["filter"] = src
	}

	return map[string]interface{}{"any_of": source}, nil
}

func (r *IntervalQueryRuleAnyOf) isIntervalQueryRule() bool {
	return true
}
