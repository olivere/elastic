package elastic

// SuggesterContextQuery is used to define context information within
// a suggestion request.
type SuggesterContextQuery interface {
	Source() interface{}
}
