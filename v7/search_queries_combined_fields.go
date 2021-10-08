// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import "fmt"

// CombinedFieldsQuery supports searching multiple text fields as if their
// contents had been indexed into one combined field.
//
// For more details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/7.13/query-dsl-combined-fields-query.html
type CombinedFieldsQuery struct {
	text                            interface{}
	fields                          []string
	fieldBoosts                     map[string]*float64
	autoGenerateSynonymsPhraseQuery *bool
	operator                        string // AND or OR
	minimumShouldMatch              string
	zeroTermsQuery                  string
}

// NewCombinedFieldsQuery creates and initializes a new CombinedFieldsQuery.
func NewCombinedFieldsQuery(text interface{}, fields ...string) *CombinedFieldsQuery {
	q := &CombinedFieldsQuery{
		text:        text,
		fieldBoosts: make(map[string]*float64),
	}
	q.fields = append(q.fields, fields...)
	return q
}

// Field adds a field to run the multi match against.
func (q *CombinedFieldsQuery) Field(field string) *CombinedFieldsQuery {
	q.fields = append(q.fields, field)
	return q
}

// FieldWithBoost adds a field to run the multi match against with a specific boost.
func (q *CombinedFieldsQuery) FieldWithBoost(field string, boost float64) *CombinedFieldsQuery {
	q.fields = append(q.fields, field)
	q.fieldBoosts[field] = &boost
	return q
}

// AutoGenerateSynonymsPhraseQuery indicates whether phrase queries should be
// automatically generated for multi terms synonyms. Defaults to true.
func (q *CombinedFieldsQuery) AutoGenerateSynonymsPhraseQuery(enable bool) *CombinedFieldsQuery {
	q.autoGenerateSynonymsPhraseQuery = &enable
	return q
}

// Operator sets the operator to use when using boolean query.
// It can be either AND or OR (default).
func (q *CombinedFieldsQuery) Operator(operator string) *CombinedFieldsQuery {
	q.operator = operator
	return q
}

// MinimumShouldMatch represents the minimum number of optional should clauses
// to match.
func (q *CombinedFieldsQuery) MinimumShouldMatch(minimumShouldMatch string) *CombinedFieldsQuery {
	q.minimumShouldMatch = minimumShouldMatch
	return q
}

// ZeroTermsQuery can be "all" or "none".
func (q *CombinedFieldsQuery) ZeroTermsQuery(zeroTermsQuery string) *CombinedFieldsQuery {
	q.zeroTermsQuery = zeroTermsQuery
	return q
}

// Source returns JSON for the query.
func (q *CombinedFieldsQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	combinedFields := make(map[string]interface{})
	source["combined_fields"] = combinedFields

	combinedFields["query"] = q.text

	fields := []string{}
	for _, field := range q.fields {
		if boost, found := q.fieldBoosts[field]; found {
			if boost != nil {
				fields = append(fields, fmt.Sprintf("%s^%f", field, *boost))
			} else {
				fields = append(fields, field)
			}
		} else {
			fields = append(fields, field)
		}
	}
	combinedFields["fields"] = fields

	if q.autoGenerateSynonymsPhraseQuery != nil {
		combinedFields["auto_generate_synonyms_phrase_query"] = q.autoGenerateSynonymsPhraseQuery
	}
	if q.operator != "" {
		combinedFields["operator"] = q.operator
	}
	if q.minimumShouldMatch != "" {
		combinedFields["minimum_should_match"] = q.minimumShouldMatch
	}
	if q.zeroTermsQuery != "" {
		combinedFields["zero_terms_query"] = q.zeroTermsQuery
	}
	return source, nil
}
