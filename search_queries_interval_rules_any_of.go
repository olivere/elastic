package elastic

type AnyOfRule struct {
	intervals []IntervalRule
	filter    *FilterRule
}

var _ IntervalRule = &AnyOfRule{}

func NewAnyOfRule(intervals ...IntervalRule) *AnyOfRule {
	return &AnyOfRule{intervals: intervals}
}
func (r *AnyOfRule) Filter(filter *FilterRule) *AnyOfRule {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *AnyOfRule) Source() (interface{}, error) {
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

func (r *AnyOfRule) IsIntervalRule() bool {
	return true
}
