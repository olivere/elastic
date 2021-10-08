// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestComponentTemplatesLifecycle(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	const templateName = "template_1"

	// Always make sure the component template is deleted
	defer func() {
		_, _ = client.IndexDeleteComponentTemplate(templateName).Pretty(true).Do(context.Background())
	}()

	// Create an component template
	{
		resp, err := client.IndexPutComponentTemplate(templateName).Pretty(true).BodyString(`{
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
			t.Fatalf("expected to successfully create component template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on creating component template")
		}
		if want, have := true, resp.Acknowledged; want != have {
			t.Errorf("expected Acknowledged=%v, got %v", want, have)
		}
	}

	// Get the component template
	{
		resp, err := client.IndexGetComponentTemplate(templateName).Pretty(true).Do(context.Background())
		if err != nil {
			t.Fatalf("expected to successfully get component template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on getting component template")
		}
	}

	// Delete the component template
	{
		resp, err := client.IndexDeleteComponentTemplate(templateName).Pretty(true).Do(context.Background())
		if err != nil {
			t.Fatalf("expected to successfully delete component template, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response on deleting component template")
		}
	}
}
