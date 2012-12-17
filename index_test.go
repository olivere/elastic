// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const (
	testIndexName  = "elastic-test"
	testIndexName2 = "elastic-test2"
)

type tweet struct {
	User     string    `json:"user"`
	Message  string    `json:"message"`
	Retweets int       `json:"retweets"`
	Created  time.Time `json:"created,omitempty"`
}

func setupTestClient(t *testing.T) *Client {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	client.DeleteIndex(testIndexName).Do()
	client.DeleteIndex(testIndexName2).Do()

	return client
}

func setupTestClientAndCreateIndex(t *testing.T) *Client {
	client := setupTestClient(t)

	// Create index
	createIndex, err := client.CreateIndex(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !createIndex.Ok {
		t.Errorf("expected CreateIndexResult.Ok %q; got %q", true, createIndex.Ok)
	}

	// Create second index
	createIndex2, err := client.CreateIndex(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !createIndex2.Ok {
		t.Errorf("expected CreateIndexResult.Ok %q; got %q", true, createIndex2.Ok)
	}

	return client
}

func TestIndexLifecycle(t *testing.T) {
	client := setupTestClient(t)

	// Create index
	createIndex, err := client.CreateIndex(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !createIndex.Ok {
		t.Errorf("expected CreateIndexResult.Ok %q; got %q", true, createIndex.Ok)
	}
	if !createIndex.Ack {
		t.Errorf("expected CreateIndexResult.Ack %q; got %q", true, createIndex.Ack)
	}

	// Check if index exists
	indexExists, err := client.IndexExists(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !indexExists {
		t.Fatalf("index %s should exist, but doesn't\n", testIndexName)
	}

	// Delete index
	deleteIndex, err := client.DeleteIndex(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !deleteIndex.Ok {
		t.Errorf("expected DeleteIndexResult.Ok %q; got %q", true, deleteIndex.Ok)
	}
	if !deleteIndex.Ack {
		t.Errorf("expected DeleteIndexResult.Ack %q; got %q", true, deleteIndex.Ack)
	}

	// Check if index exists
	indexExists, err = client.IndexExists(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if indexExists {
		t.Fatalf("index %s should not exist, but does\n", testIndexName)
	}
}

func TestIndexExistScenarios(t *testing.T) {
	client := setupTestClient(t)

	// Should return false if index does not exist
	indexExists, err := client.IndexExists(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if indexExists {
		t.Fatalf("expected index exists to return %q, got %q\n", false, indexExists)
	}

	// Create index
	createIndex, err := client.CreateIndex(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !createIndex.Ok {
		t.Errorf("expected CreateIndexResult.Ok %q; got %q", true, createIndex.Ok)
	}
	if !createIndex.Ack {
		t.Errorf("expected CreateIndexResult.Ack %q; got %q", true, createIndex.Ack)
	}

	// Should return true if index does not exist
	indexExists, err = client.IndexExists(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if !indexExists {
		t.Fatalf("expected index exists to return %q, got %q\n", true, indexExists)
	}
}

func TestDocumentLifecycle(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and ElasticSearch."}

	// Add a document
	indexResult, err := client.Index().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		BodyJson(&tweet1).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if !indexResult.Ok {
		t.Errorf("expected IndexResult.Ok %q; got %q", true, indexResult.Ok)
	}

	// Exists
	exists, err := client.Exists().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("expected exists %q; got %q", true, exists)
	}

	// Get document
	getResult, err := client.Get().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if getResult.Index != testIndexName {
		t.Errorf("expected GetResult.Index %q; got %q", testIndexName, getResult.Index)
	}
	if getResult.Type != "tweet" {
		t.Errorf("expected GetResult.Type %q; got %q", "tweet", getResult.Type)
	}
	if getResult.Id != "1" {
		t.Errorf("expected GetResult.Id %q; got %q", "1", getResult.Id)
	}
	if getResult.Source == nil {
		t.Errorf("expected GetResult.Source to be != nil; got nil")
	}

	// Decode the Source field
	var tweetGot tweet
	err = json.Unmarshal(*getResult.Source, &tweetGot)
	if err != nil {
		t.Fatal(err)
	}
	if tweetGot.User != tweet1.User {
		t.Errorf("expected Tweet.User to be %q; got %q", tweet1.User, tweetGot.User)
	}
	if tweetGot.Message != tweet1.Message {
		t.Errorf("expected Tweet.Message to be %q; got %q", tweet1.Message, tweetGot.Message)
	}

	// Delete document again
	deleteResult, err := client.Delete().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if !deleteResult.Ok {
		t.Errorf("expected DeleteResult.Ok %q; got %q", true, deleteResult.Ok)
	}

	// Exists
	exists, err = client.Exists().Index(testIndexName).Type("tweet").Id("1").Do()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("expected exists %q; got %q", false, exists)
	}
}
