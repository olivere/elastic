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
	testMapping    = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"location":{
					"type":"geo_point"
				}
			}
		}
	}
}
`
)

type tweet struct {
	User     string    `json:"user"`
	Message  string    `json:"message"`
	Retweets int       `json:"retweets"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Location string    `json:"location,omitempty"`
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
	createIndex, err := client.CreateIndex(testIndexName).Body(testMapping).Do()
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Create second index
	createIndex2, err := client.CreateIndex(testIndexName2).Body(testMapping).Do()
	if err != nil {
		t.Fatal(err)
	}
	if createIndex2 == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex2)
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
	if !createIndex.Acknowledged {
		t.Errorf("expected CreateIndexResult.Acknowledged %q; got %q", true, createIndex.Acknowledged)
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
	if !deleteIndex.Acknowledged {
		t.Errorf("expected DeleteIndexResult.Acknowledged %q; got %q", true, deleteIndex.Acknowledged)
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
	if !createIndex.Acknowledged {
		t.Errorf("expected CreateIndexResult.Ack %q; got %q", true, createIndex.Acknowledged)
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
	if indexResult == nil {
		t.Errorf("expected result to be != nil; got: %v", indexResult)
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
	if deleteResult == nil {
		t.Errorf("expected result to be != nil; got: %v", deleteResult)
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
