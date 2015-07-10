// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// ScriptQuery allows to define scripts as filters.
//
// For details, see
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-script-query.html
type ScriptQuery struct {
	script    string
	queryName string
}

// NewScriptQuery creates and initializes a new ScriptQuery.
func NewScriptQuery(script string) *ScriptQuery {
	return &ScriptQuery{
		script: script,
	}
}

// QueryName sets the query name for the filter that can be used
// when searching for matched_filters per hit
func (q *ScriptQuery) QueryName(queryName string) *ScriptQuery {
	q.queryName = queryName
	return q
}

// Source returns JSON for the query.
func (q *ScriptQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source["script"] = params

	params["script"] = q.script

	if q.queryName != "" {
		params["_name"] = q.queryName
	}
	return source, nil
}
