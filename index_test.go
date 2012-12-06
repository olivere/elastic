// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
	"testing"
)

const (
	testIndexName = "elastic-test"
)

func setupTestClient(t *testing.T) *Client {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	client.DeleteIndex(testIndexName).Do()

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
