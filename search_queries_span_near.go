package elastic

// SpanNearQuery Matches spans which are near one another.
type SpanNearQuery struct {
	clauses []Query
	slop    int64
	inOrder bool
}

// NewSpanNearQuery creates and initializes a new SpanNearQuery.
func NewSpanNearQuery() *SpanNearQuery {
	return &SpanNearQuery{}
}

// InOrder sets the ordering for this query.
func (sn *SpanNearQuery) InOrder(inOrder bool) *SpanNearQuery {
	sn.inOrder = inOrder
	return sn
}

// Slop sets the slop for this query.
func (sn *SpanNearQuery) Slop(slop int64) *SpanNearQuery {
	sn.slop = slop
	return sn
}

// Clauses adds span clauses for the SpanNearQuery.
func (sn *SpanNearQuery) Clauses(clauses ...Query) *SpanNearQuery {
	sn.clauses = append(sn.clauses, clauses...)
	return sn
}

// Source returns JSON for the query.
func (sn *SpanNearQuery) Source() (interface{}, error) {
	//"span_near" : {
	//	"clauses" : [
	//	{ "span_term" : { "field" : "value1" } },
	//	{ "span_term" : { "field" : "value2" } },
	//	{ "span_term" : { "field" : "value3" } }
	//],
	//	"slop" : 12,
	//		"in_order" : false
	//}

	source := make(map[string]interface{})
	spanNearQuery := make(map[string]interface{})
	source["span_near"] = spanNearQuery
	spanNearQuery["slop"] = sn.slop
	spanNearQuery["in_order"] = sn.inOrder

	if len(sn.clauses) > 0 {
		var clauses []interface{}
		for _, spanClause := range sn.clauses {
			src, err := spanClause.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		spanNearQuery["clauses"] = clauses
	}

	return source, nil
}
