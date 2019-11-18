package elastic

type WildcardRule struct {
	pattern  string
	analyzer string
	useField string
}

var _ IntervalRule = &WildcardRule{}

func NewWildcardRule(pattern string) *WildcardRule {
	return &WildcardRule{pattern: pattern}
}

func (r *WildcardRule) Analyzer(analyzer string) *WildcardRule {
	r.analyzer = analyzer
	return r
}
func (r *WildcardRule) UseField(useField string) *WildcardRule {
	r.useField = useField
	return r
}

// Source returns JSON for the function score query.
func (r *WildcardRule) Source() (interface{}, error) {
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

func (r *WildcardRule) IsIntervalRule() bool {
	return true
}
