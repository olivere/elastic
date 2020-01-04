package elastic

type IntervalQueryRuleFilter struct {
	after          Query
	before         Query
	containedBy    Query
	containing     Query
	overlapping    Query
	notContainedBy Query
	notContaining  Query
	notOverlapping Query
	script         *Script
}

var _ IntervalQueryRule = &IntervalQueryRuleFilter{}

func NewIntervalQueryRuleFilter() *IntervalQueryRuleFilter {
	return &IntervalQueryRuleFilter{}
}

func (r *IntervalQueryRuleFilter) After(after Query) *IntervalQueryRuleFilter {
	r.after = after
	return r
}
func (r *IntervalQueryRuleFilter) Before(before Query) *IntervalQueryRuleFilter {
	r.before = before
	return r
}
func (r *IntervalQueryRuleFilter) ContainedBy(containedBy Query) *IntervalQueryRuleFilter {
	r.containedBy = containedBy
	return r
}
func (r *IntervalQueryRuleFilter) Containing(containing Query) *IntervalQueryRuleFilter {
	r.containing = containing
	return r
}
func (r *IntervalQueryRuleFilter) Overlapping(overlapping Query) *IntervalQueryRuleFilter {
	r.overlapping = overlapping
	return r
}
func (r *IntervalQueryRuleFilter) NotContainedBy(notContainedBy Query) *IntervalQueryRuleFilter {
	r.notContainedBy = notContainedBy
	return r
}
func (r *IntervalQueryRuleFilter) NotContaining(notContaining Query) *IntervalQueryRuleFilter {
	r.notContaining = notContaining
	return r
}
func (r *IntervalQueryRuleFilter) NotOverlapping(notOverlapping Query) *IntervalQueryRuleFilter {
	r.notOverlapping = notOverlapping
	return r
}
func (r *IntervalQueryRuleFilter) Script(script *Script) *IntervalQueryRuleFilter {
	r.script = script
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleFilter) Source() (interface{}, error) {
	source := make(map[string]interface{})

	if r.before != nil {
		src, err := r.before.Source()
		if err != nil {
			return nil, err
		}
		source["before"] = src
	}

	if r.after != nil {
		src, err := r.after.Source()
		if err != nil {
			return nil, err
		}
		source["after"] = src
	}

	if r.containedBy != nil {
		src, err := r.containedBy.Source()
		if err != nil {
			return nil, err
		}
		source["contained_by"] = src
	}

	if r.containing != nil {
		src, err := r.containing.Source()
		if err != nil {
			return nil, err
		}
		source["containing"] = src
	}

	if r.overlapping != nil {
		src, err := r.overlapping.Source()
		if err != nil {
			return nil, err
		}
		source["overlapping"] = src
	}

	if r.notContainedBy != nil {
		src, err := r.notContainedBy.Source()
		if err != nil {
			return nil, err
		}
		source["not_contained_by"] = src
	}

	if r.notContaining != nil {
		src, err := r.notContaining.Source()
		if err != nil {
			return nil, err
		}
		source["not_containing"] = src
	}

	if r.notOverlapping != nil {
		src, err := r.notOverlapping.Source()
		if err != nil {
			return nil, err
		}
		source["not_overlapping"] = src
	}

	if r.script != nil {
		src, err := r.script.Source()
		if err != nil {
			return nil, err
		}
		source["script"] = src
	}

	// todo: not so clear from docs, if filter can be a top-level rule, and so
	// 	if it does need a wrapper map[string]interface{}{"filter": ..} like other rules do

	return source, nil
}

func (r *IntervalQueryRuleFilter) isIntervalQueryRule() bool {
	return true
}
