package elastic

// SpanMultiQuery allows you to wrap a multi term query (one of wildcard, fuzzy, prefix, range or regexp query)
// as a span query, so it can be nested.
type SpanMultiQuery struct {
	query Query
}

// NewSpanMultiQuery creates and initializes a new SpanMultiQuery.
func NewSpanMultiQuery() *SpanMultiQuery {
	return &SpanMultiQuery{}
}

// Match adds queries to match on span.
func (sm *SpanMultiQuery) Match(query Query) *SpanMultiQuery {
	sm.query = query
	return sm
}

// Source returns JSON for the query.
func (sm *SpanMultiQuery) Source() (interface{}, error) {
	//"span_multi":{
	//	"match":{
	//		"prefix" : { "user" :  { "value" : "ki" } }
	//	}
	//}

	source := make(map[string]interface{})
	spanMultiMatchQuery := make(map[string]interface{})
	spanMatchQuery := make(map[string]interface{})

	source["span_multi"] = spanMultiMatchQuery
	spanMultiMatchQuery["match"] = spanMatchQuery
	if sm.query != nil {

		spanMatchQuery, err := sm.query.Source()
		if err != nil {
			return nil, err
		}
		spanMultiMatchQuery["match"] = spanMatchQuery
	}


	return source, nil
}

