// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opensearch_test

import (
	"context"

	opensearch "github.com/disaster37/opensearch/v2"
)

func ExamplePrefixQuery() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Define wildcard query
	q := opensearch.NewPrefixQuery("user", "oli")
	q = q.QueryName("my_query_name")

	searchResult, err := client.Search().
		Index("twitter").
		Query(q).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	_ = searchResult
}
