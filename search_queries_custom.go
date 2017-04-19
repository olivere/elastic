// Copyright 2012-2015 Oliver Eilhard, John Stanford. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import "encoding/json"

// CustomQuery can be used to treat a string representation of an ES query
// as a Query.  Example usage:
//    q := CustomQuery("{\"query\":{\"match_all\":{}}}")
//    db.Search().Query(q).From(1).Size(100).Do()
type CustomQuery string

// Source returns the JSON encoded body
func (q CustomQuery) Source() (interface{}, error) {
	var f interface{}
	err := json.Unmarshal([]byte(q), &f)
	return f, err
}
