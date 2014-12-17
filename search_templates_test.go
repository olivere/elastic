// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestSearchTemplates(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	// Template
	tmpl := `{"template":{"query":{"match":{"title":"{{query_string}}"}}}}`

	// Create template
	cresp, err := client.PutTemplate().Id("elastic-test").BodyString(tmpl).Do()
	if err != nil {
		t.Fatal(err)
	}
	if cresp == nil {
		t.Fatalf("expected response != nil; got: %v", cresp)
	}
	if !cresp.Created {
		t.Errorf("expected created = %v; got: %v", true, cresp.Created)
	}

	// Get template
	resp, err := client.GetTemplate().Id("elastic-test").Do()
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
	dresp, err := client.DeleteTemplate().Id("elastic-test").Do()
	if err != nil {
		t.Fatal(err)
	}
	if dresp == nil {
		t.Fatalf("expected response != nil; got: %v", dresp)
	}
	if !dresp.Found {
		t.Fatalf("expected found = %v; got: %v", true, dresp.Found)
	}
}
