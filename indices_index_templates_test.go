// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndexTemplatesLifecycle(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	const templateName = "template_1"

	// Always make sure the index template is deleted
	defer func() {
		_, _ = client.IndexDeleteIndexTemplate(templateName).Pretty(true).Do(context.Background())
	}()

	// Create an index template
	{
		resp, err := client.IndexPutIndexTemplate(templateName).Pretty(true).BodyString(`{
			"index_patterns": ["elastic-index-templates-*"],
			"priority": 1,
			"template": {
				"settings": {
					"number_of_shards": 2,
					"number_of_replicas": 0
				},
				"mappings": {
					"_source": { "enabled": true }
				}
			}
		}`).Do(context.Background())
		if err != nil {
			t.Fatalf("expected to successfully create index template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on creating index template")
		}
		if want, have := true, resp.Acknowledged; want != have {
			t.Errorf("expected Acknowledged=%v, got %v", want, have)
		}
	}

	// Get the index template
	{
		resp, err := client.IndexGetIndexTemplate(templateName).Pretty(true).Do(context.Background())
		if err != nil {
			t.Fatalf("expected to successfully get index template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on getting index template")
		}
	}

	// Delete the index template
	{
		resp, err := client.IndexDeleteIndexTemplate(templateName).Pretty(true).Do(context.Background())
		if err != nil {
			t.Fatalf("expected to successfully delete index template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on deleting index template")
		}
	}
}
