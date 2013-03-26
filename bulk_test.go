// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	_ "encoding/json"
	"testing"
)

func TestBulk(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}
	tweet2 := tweet{User: "sandrae", Message: "Dancing all night long. Yeah."}

	index1Req := NewBulkIndexRequest(testIndexName, "tweet", "1", tweet1)
	index2Req := NewBulkIndexRequest(testIndexName, "tweet", "2", tweet2)
	delete1Req := NewBulkDeleteRequest(testIndexName, "tweet", "1")

	bulkRequest := client.Bulk() //.Debug(true)
	bulkRequest = bulkRequest.Add(index1Req)
	bulkRequest = bulkRequest.Add(index2Req)
	bulkRequest = bulkRequest.Add(delete1Req)

	if bulkRequest.NumberOfActions() != 3 {
		t.Errorf("expected bulkRequest.NumberOfActions %q; got %q", 3, bulkRequest.NumberOfActions())
	}

	bulkResponse, err := bulkRequest.Do()
	if err != nil {
		t.Fatal(err)
	}
	if bulkResponse == nil {
		t.Errorf("expected bulkResponse to be != nil; got nil")
	}

	if bulkRequest.NumberOfActions() != 0 {
		t.Errorf("expected bulkRequest.NumberOfActions %q; got %q", 0, bulkRequest.NumberOfActions())
	}

	// Document with Id="1" should not exist
	exists, err := client.Exists().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("expected exists %q; got %q", false, exists)
	}

	// Document with Id="2" should exist
	exists, err = client.Exists().Index(testIndexName).Type("tweet").Id("2").Do()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("expected exists %q; got %q", true, exists)
	}
}

func TestBulkWithIndexSetOnClient(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}
	tweet2 := tweet{User: "sandrae", Message: "Dancing all night long. Yeah."}

	index1Req := &BulkIndexRequest{Id: "1", Data: tweet1}
	index2Req := &BulkIndexRequest{Id: "2", Data: tweet2}
	delete1Req := &BulkDeleteRequest{Id: "1"}

	bulkRequest := client.Bulk().Index(testIndexName).Type("tweet") //.Debug(true)
	bulkRequest = bulkRequest.Add(index1Req)
	bulkRequest = bulkRequest.Add(index2Req)
	bulkRequest = bulkRequest.Add(delete1Req)

	if bulkRequest.NumberOfActions() != 3 {
		t.Errorf("expected bulkRequest.NumberOfActions %q; got %q", 3, bulkRequest.NumberOfActions())
	}

	bulkResponse, err := bulkRequest.Do()
	if err != nil {
		t.Fatal(err)
	}
	if bulkResponse == nil {
		t.Errorf("expected bulkResponse to be != nil; got nil")
	}

	// Document with Id="1" should not exist
	exists, err := client.Exists().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("expected exists %q; got %q", false, exists)
	}

	// Document with Id="2" should exist
	exists, err = client.Exists().Index(testIndexName).Type("tweet").Id("2").Do()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("expected exists %q; got %q", true, exists)
	}
}
