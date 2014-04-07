// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Aggregations can be seen as a unit-of-work that build
// analytic information over a set of documents. It is
// (in many senses) the follow-up of facets in Elasticsearch.
// For more details about aggregations, visit:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-aggregations.html
type Aggregation interface {
	Source() interface{}
}
