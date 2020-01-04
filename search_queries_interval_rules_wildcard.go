package elastic

type IntervalQueryRuleWildcard struct {
	pattern  string
	analyzer string
	useField string
}

var _ IntervalQueryRule = &IntervalQueryRuleWildcard{}

func NewIntervalQueryRuleWildcard(pattern string) *IntervalQueryRuleWildcard {
	return &IntervalQueryRuleWildcard{pattern: pattern}
}

func (r *IntervalQueryRuleWildcard) Analyzer(analyzer string) *IntervalQueryRuleWildcard {
	r.analyzer = analyzer
	return r
}
func (r *IntervalQueryRuleWildcard) UseField(useField string) *IntervalQueryRuleWildcard {
	r.useField = useField
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRuleWildcard) Source() (interface{}, error) {
	source := make(map[string]interface{})

	source["pattern"] = r.pattern

	if r.analyzer != "" {
		source["analyzer"] = r.analyzer
	}
	if r.useField != "" {
		source["use_field"] = r.useField
	}

	return map[string]interface{}{"wildcard": source}, nil
}

func (r *IntervalQueryRuleWildcard) isIntervalQueryRule() bool {
	return true
}
