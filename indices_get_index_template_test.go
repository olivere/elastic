// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndexGetIndexTemplateURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		Name     string
		Expected string
	}{
		{
			"template_1",
			"/_index_template/template_1",
		},
	}

	for _, test := range tests {
		path, _, err := client.IndexGetIndexTemplate(test.Name).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}

func TestIndexGetIndexTemplateService(t *testing.T) {
	// client := setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	create := true
	body := `
{
	"index_patterns": ["te*"],
	"priority": 1,
	"template": {
		"settings": {
			"index": {
		  		"number_of_shards": 1
			}
	  	},
	  	"mappings": {
			"_source": {
		  		"enabled": false
			},
			"properties": {
		  		"host_name": {
					"type": "keyword"
		  		},
		  		"created_at": {
					"type": "date",
					"format": "yyyy MM dd HH:mm:ss Z"
		  		}
			}
	  	}
  	}
}
`
	_, err := client.IndexPutIndexTemplate("template_1").BodyString(body).Create(create).Pretty(true).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer client.IndexDeleteIndexTemplate("template_1").Pretty(true).Do(context.TODO())

	res, err := client.IndexGetIndexTemplate("template_1").Pretty(true).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatalf("expected result; got: %v", res)
	}
	template, found := res.IndexTemplates.ByName("template_1")
	if !found {
		t.Fatalf("expected template %q to be found; got: %v", "template_1", found)
	}
	if template == nil {
		t.Fatalf("expected template %q to be != nil; got: %v", "template_1", template)
	}
	if template.IndexTemplate == nil {
		t.Fatalf("expected index template of template %q to be != nil; got: %v", "template_1", template.IndexTemplate)
	}
	if len(template.IndexTemplate.IndexPatterns) != 1 || template.IndexTemplate.IndexPatterns[0] != "te*" {
		t.Fatalf("expected index settings of %q to be [\"index1\"]; got: %v", testIndexName, template.IndexTemplate.IndexPatterns)
	}
}
