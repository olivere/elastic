// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// For more details, see
// http://www.elasticsearch.org/guide/reference/api/search/phrase-suggest/
type PhraseSuggester struct {
	Suggester
	name      string
	text      string
	field     string
	analyzer  string
	size      *int
	shardSize *int

	maxErrors                *float32
	separator                *string
	realWorldErrorLikelihood *float32
	confidence               *float32
	//generators map[string][]*CandidateGenerator
	gramSize        *int
	smoothing       interface{}
	forceUnigrams   *bool
	tokenLimit      *int
	preTag, postTag *string
}

// Creates a new phrase suggester.
func NewPhraseSuggester(name string) PhraseSuggester {
	return PhraseSuggester{name: name}
}

func (q PhraseSuggester) Name() string {
	return q.name
}

func (q PhraseSuggester) Text(text string) PhraseSuggester {
	q.text = text
	return q
}

func (q PhraseSuggester) Field(field string) PhraseSuggester {
	q.field = field
	return q
}

func (q PhraseSuggester) Analyzer(analyzer string) PhraseSuggester {
	q.analyzer = analyzer
	return q
}

func (q PhraseSuggester) Size(size int) PhraseSuggester {
	q.size = &size
	return q
}

func (q PhraseSuggester) ShardSize(shardSize int) PhraseSuggester {
	q.shardSize = &shardSize
	return q
}

func (q PhraseSuggester) GramSize(gramSize int) PhraseSuggester {
	if gramSize >= 1 {
		q.gramSize = &gramSize
	}
	return q
}

func (q PhraseSuggester) MaxErrors(maxErrors float32) PhraseSuggester {
	q.maxErrors = &maxErrors
	return q
}

func (q PhraseSuggester) Separator(separator string) PhraseSuggester {
	q.separator = &separator
	return q
}

func (q PhraseSuggester) RealWorldErrorLikelihood(realWorldErrorLikelihood float32) PhraseSuggester {
	q.realWorldErrorLikelihood = &realWorldErrorLikelihood
	return q
}

func (q PhraseSuggester) Confidence(confidence float32) PhraseSuggester {
	q.confidence = &confidence
	return q
}

/*
func (q PhraseSuggester) CandidateGenerator(generator *CandidateGenerator) PhraseSuggester {
	q.generators = append(q.generators, generator)
	return q
}
*/

func (q PhraseSuggester) ForceUnigrams(forceUnigrams bool) PhraseSuggester {
	q.forceUnigrams = &forceUnigrams
	return q
}

func (q PhraseSuggester) Smoothing(smoothing interface{}) PhraseSuggester {
	q.smoothing = smoothing
	return q
}

func (q PhraseSuggester) TokenLimit(tokenLimit int) PhraseSuggester {
	q.tokenLimit = &tokenLimit
	return q
}

func (q PhraseSuggester) Highlight(preTag, postTag string) PhraseSuggester {
	q.preTag = &preTag
	q.postTag = &postTag
	return q
}

// simplePhraseSuggesterRequest is necessary because the order in which
// the JSON elements are routed to Elasticsearch is relevant.
// We got into trouble when using plain maps because the text element
// needs to go before the simple_phrase element.
type phraseSuggesterRequest struct {
	Text   string      `json:"text"`
	Phrase interface{} `json:"phrase"`
}

// Creates the source for the phrase suggester.
func (q PhraseSuggester) Source(includeName bool) interface{} {
	ps := &phraseSuggesterRequest{}

	if q.text != "" {
		ps.Text = q.text
	}

	suggester := make(map[string]interface{})
	ps.Phrase = suggester

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

	// Phase-specified parameters
	if q.realWorldErrorLikelihood != nil {
		suggester["real_word_error_likelihood"] = *q.realWorldErrorLikelihood
	}
	if q.confidence != nil {
		suggester["confidence"] = *q.confidence
	}
	if q.separator != nil {
		suggester["separator"] = *q.separator
	}
	if q.maxErrors != nil {
		suggester["max_errors"] = *q.maxErrors
	}
	if q.gramSize != nil {
		suggester["gram_size"] = *q.gramSize
	}
	if q.forceUnigrams != nil {
		suggester["force_unigrams"] = *q.forceUnigrams
	}
	if q.tokenLimit != nil {
		suggester["token_limit"] = *q.tokenLimit
	}
	// TODO(oe) Add CandidateGenerators here
	if q.smoothing != nil {
		suggester["smoothing"] = q.smoothing
	}
	if q.preTag != nil {
		hl := make(map[string]string)
		hl["pre_tag"] = *q.preTag
		if q.postTag != nil {
			hl["post_tag"] = *q.postTag
		}
		suggester["highlight"] = hl
	}

	if !includeName {
		return ps
	}

	source := make(map[string]interface{})
	source[q.name] = ps
	return source
}
