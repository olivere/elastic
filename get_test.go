// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestGet(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Count documents
	count, err := client.Count(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected Count = %d; got %d", 3, count)
	}

	// Get document 1
	res, err := client.Get().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if res.Exists != true {
		t.Errorf("expected Exists = true; got %v", res.Exists)
	}
	if res.Source == nil {
		t.Errorf("expected Source != nil; got %v", res.Source)
	}

	// Get non existent document 99
	res, err = client.Get().Index(testIndexName).Type("tweet").Id("99").Do()
	if err != nil {
		t.Fatal(err)
	}
	if res.Exists != false {
		t.Errorf("expected Exists = false; got %v", res.Exists)
	}
	if res.Source != nil {
		t.Errorf("expected Source == nil; got %v", res.Source)
	}
}
