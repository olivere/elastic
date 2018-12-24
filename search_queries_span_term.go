package elastic

// TermQuery matches spans containing a term.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/6.2/query-dsl-span-term-query.html

// SpanTermQuery matches spans containing a term.
type SpanTermQuery struct {
	value interface{}
	name  string
	boost *float64
}

// NewSpanTermQuery creates and initializes a new SpanTermQuery.
func NewSpanTermQuery(name string, value interface{}) *SpanTermQuery {
	return &SpanTermQuery{name: name, value: value}
}

// Boost sets the boost for this query.
func (span *SpanTermQuery) Boost(boost float64) *SpanTermQuery {
	span.boost = &boost
	return span
}

// Source returns JSON for the query.
func (span *SpanTermQuery) Source() (interface{}, error) {
	// {"span_term":{"name":"value"}}
	source := make(map[string]interface{})
	spanTermQuery := make(map[string]interface{})
	source["span_term"] = spanTermQuery
	spanTermQuery[span.name] = span.value
	if span.boost != nil {
		spanTermQuery["boost"] = span.boost
	}

	return source, nil
}