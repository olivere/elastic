package elastic

type PrefixRule struct {
	prefix   string
	analyzer string
	useField string
}

var _ IntervalRule = &PrefixRule{}

func NewPrefixRule(prefix string) *PrefixRule {
	return &PrefixRule{prefix: prefix}
}

func (r *PrefixRule) Analyzer(analyzer string) *PrefixRule {
	r.analyzer = analyzer
	return r
}
func (r *PrefixRule) UseField(useField string) *PrefixRule {
	r.useField = useField
	return r
}

// Source returns JSON for the function score query.
func (r *PrefixRule) Source() (interface{}, error) {
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

func (r *PrefixRule) IsIntervalRule() bool {
	return true
}
