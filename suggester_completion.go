// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// For more details, see
// http://www.elasticsearch.org/guide/reference/api/search/completion-suggest/
type CompletionSuggester struct {
	Suggester
	name      string
	text      string
	field     string
	analyzer  string
	size      *int
	shardSize *int
}

// Creates a new completion suggester.
func NewCompletionSuggester(name string) CompletionSuggester {
	return CompletionSuggester{name: name}
}

func (q CompletionSuggester) Name() string {
	return q.name
}

func (q CompletionSuggester) Text(text string) CompletionSuggester {
	q.text = text
	return q
}

func (q CompletionSuggester) Field(field string) CompletionSuggester {
	q.field = field
	return q
}

func (q CompletionSuggester) Analyzer(analyzer string) CompletionSuggester {
	q.analyzer = analyzer
	return q
}

func (q CompletionSuggester) Size(size int) CompletionSuggester {
	q.size = &size
	return q
}

func (q CompletionSuggester) ShardSize(shardSize int) CompletionSuggester {
	q.shardSize = &shardSize
	return q
}

// completionSuggesterRequest is necessary because the order in which
// the JSON elements are routed to Elasticsearch is relevant.
// We got into trouble when using plain maps because the text element
// needs to go before the completion element.
type completionSuggesterRequest struct {
	Text       string      `json:"text"`
	Completion interface{} `json:"completion"`
}

// Creates the source for the completion suggester.
func (q CompletionSuggester) Source(includeName bool) interface{} {
	cs := &completionSuggesterRequest{}

	if q.text != "" {
		cs.Text = q.text
	}

	suggester := make(map[string]interface{})
	cs.Completion = suggester

	if q.analyzer != "" {
		suggester["analyzer"] = q.analyzer
	}

	if q.field != "" {
		suggester["field"] = q.field
	}

	if q.size != nil {
		suggester["size"] = *q.size
	}

	if q.shardSize != nil {
		suggester["shard_size"] = *q.shardSize
	}

	// TODO(oe) Add competion-suggester specific parameters here

	if !includeName {
		return cs
	}

	source := make(map[string]interface{})
	source[q.name] = cs
	return source
}
