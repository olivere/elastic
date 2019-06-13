// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIndicesLifecycle(t *testing.T) {
	client := setupTestClient(t)

	// Create index
	createIndex, err := client.CreateIndex(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !createIndex.Acknowledged {
		t.Errorf("expected IndicesCreateResult.Acknowledged %v; got %v", true, createIndex.Acknowledged)
	}

	// Check if index exists
	indexExists, err := client.IndexExists(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !indexExists {
		t.Fatalf("index %s should exist, but doesn't\n", testIndexName)
	}

	// Delete index
	deleteIndex, err := client.DeleteIndex(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if !deleteIndex.Acknowledged {
		t.Errorf("expected DeleteIndexResult.Acknowledged %v; got %v", true, deleteIndex.Acknowledged)
	}

	// Check if index exists
	indexExists, err = client.IndexExists(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if indexExists {
		t.Fatalf("index %s should not exist, but does\n", testIndexName)
	}
}

func TestIndicesCreateValidate(t *testing.T) {
	client := setupTestClient(t)

	// No index name -> fail with error
	res, err := NewIndicesCreateService(client).Body(testMapping).Do(context.TODO())
	if err == nil {
		t.Fatalf("expected IndicesCreate to fail without index name")
	}
	if res != nil {
		t.Fatalf("expected result to be == nil; got: %v", res)
	}
}

func TestIndicesCreateService_buildParams(t *testing.T) {
	tests := []struct {
		pretty          bool
		masterTimeout   string
		timeout         string
		includeTypeName bool
		expectedParams  url.Values
	}{
		{
			pretty:          true,
			masterTimeout:   "3s",
			timeout:         "5s",
			includeTypeName: true,
			expectedParams: url.Values{
				"pretty":            []string{"true"},
				"master_timeout":    []string{"3s"},
				"timeout":           []string{"5s"},
				"include_type_name": []string{"true"},
			},
		},
	}

	for _, tt := range tests {
		params := NewIndicesCreateService(nil).
			Pretty(tt.pretty).
			MasterTimeout(tt.masterTimeout).
			Timeout(tt.timeout).
			IncludeTypeName(tt.includeTypeName).
			buildParams()
		if want, have := tt.expectedParams, params; !cmp.Equal(want, have) {
			t.Errorf("expected params=%#v; got: %#v\ndiff: %s", want, have, cmp.Diff(want, have))
		}
	}
}
