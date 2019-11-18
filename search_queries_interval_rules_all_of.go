package elastic

type AllOfRule struct {
	intervals []IntervalRule
	maxGaps   *int
	ordered   *bool
	filter    *FilterRule
}

var _ IntervalRule = &AllOfRule{}

func NewAllOfRule(intervals ...IntervalRule) *AllOfRule {
	return &AllOfRule{intervals: intervals}
}

func (r *AllOfRule) MaxGaps(maxGaps int) *AllOfRule {
	r.maxGaps = &maxGaps
	return r
}

func (r *AllOfRule) Ordered(ordered bool) *AllOfRule {
	r.ordered = &ordered
	return r
}

func (r *AllOfRule) Filter(filter *FilterRule) *AllOfRule {
	r.filter = filter
	return r
}

// Source returns JSON for the function score query.
func (r *AllOfRule) Source() (interface{}, error) {
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

func (r *AllOfRule) IsIntervalRule() bool {
	return true
}
