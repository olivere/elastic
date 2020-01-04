package elastic

type IntervalQueryRulePrefix struct {
	prefix   string
	analyzer string
	useField string
}

var _ IntervalQueryRule = &IntervalQueryRulePrefix{}

func NewIntervalQueryRulePrefix(prefix string) *IntervalQueryRulePrefix {
	return &IntervalQueryRulePrefix{prefix: prefix}
}

func (r *IntervalQueryRulePrefix) Analyzer(analyzer string) *IntervalQueryRulePrefix {
	r.analyzer = analyzer
	return r
}
func (r *IntervalQueryRulePrefix) UseField(useField string) *IntervalQueryRulePrefix {
	r.useField = useField
	return r
}

// Source returns JSON for the function score query.
func (r *IntervalQueryRulePrefix) Source() (interface{}, error) {
	source := make(map[string]interface{})

	source["query"] = r.prefix

	if r.analyzer != "" {
		source["analyzer"] = r.analyzer
	}
	if r.useField != "" {
		source["use_field"] = r.useField
	}

	return map[string]interface{}{"prefix": source}, nil
}

func (r *IntervalQueryRulePrefix) isIntervalQueryRule() bool {
	return true
}
