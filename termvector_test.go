// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import "testing"

func TestTermVector(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}

	// Add a document
	indexResult, err := client.Index().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		BodyJson(&tweet1).
		Refresh(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if indexResult == nil {
		t.Errorf("expected result to be != nil; got: %v", indexResult)
	}

	// TermVectors
	field := "Message"
	result, err := client.TermVector(testIndexName, "tweet", "1").
		Fields([]string{field}).
		FieldStatistics(true).
		TermStatistics(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected to return information and statistics")
	}
	if !result.Found {
		t.Errorf("expected found to be %v; got: %v", true, result.Found)
	}
}
