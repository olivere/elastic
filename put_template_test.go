// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"log"
	"os"
	"testing"

	"golang.org/x/net/context"
)

func TestSearchTemplatesLifecycle(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	// Template
	tmpl := `{"template":{"query":{"match":{"title":"{{query_string}}"}}}}`

	// Create template
	cresp, err := client.PutTemplate().Id("elastic-test").BodyString(tmpl).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if cresp == nil {
		t.Fatalf("expected response != nil; got: %v", cresp)
	}
	if !cresp.Acknowledged {
		t.Errorf("expected acknowledged = %v; got: %v", true, cresp.Acknowledged)
	}

	// Get template
	resp, err := client.GetTemplate().Id("elastic-test").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("expected response != nil; got: %v", resp)
	}
	if resp.Template == "" {
		t.Errorf("expected template != %q; got: %q", "", resp.Template)
	}

	// Delete template
	dresp, err := client.DeleteTemplate().Id("elastic-test").Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if dresp == nil {
		t.Fatalf("expected response != nil; got: %v", dresp)
	}
	if !dresp.Acknowledged {
		t.Fatalf("expected acknowledged = %v; got: %v", true, dresp.Acknowledged)
	}
}

func TestSearchTemplatesInlineQuery(t *testing.T) {
	// client := setupTestClientAndCreateIndex(t)
	client := setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Run query with (inline) search template
	// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-template-query.html
	tq := NewTemplateQuery(`{"match_{{template}}": {}}`).Var("template", "all")
	resp, err := client.Search(testIndexName).Query(tq).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("expected response != nil; got: %v", resp)
	}
	if resp.Hits == nil {
		t.Fatalf("expected response hits != nil; got: %v", resp.Hits)
	}
	if resp.Hits.TotalHits != 3 {
		t.Fatalf("expected 3 hits; got: %d", resp.Hits.TotalHits)
	}
}
