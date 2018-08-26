// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestIndexGetTemplateURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		Names    []string
		Expected string
	}{
		{
			[]string{},
			"/_template",
		},
		{
			[]string{"index1"},
			"/_template/index1",
		},
		{
			[]string{"index1", "index2"},
			"/_template/index1%2Cindex2",
		},
	}

	for _, test := range tests {
		path, _, err := client.IndexGetTemplate().Name(test.Names...).buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}

func TestIndexGetTemplateService(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)
	create := true
	body := `
{
  "index_patterns": ["te*"],
  "settings": {
    "index": {
      "number_of_shards": 1
    }
  },
  "mappings": {
    "type1": {
      "_source": {
        "enabled": false
      },
      "properties": {
        "host_name": {
          "type": "keyword"
        },
        "created_at": {
          "type": "date",
          "format": "EEE MMM dd HH:mm:ss Z YYYY"
        }
      }
    }
  }
}
`
	_, err := client.IndexPutTemplate(testIndexName).BodyString(body).Create(create).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.IndexGetTemplate(testIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatalf("expected result; got: %v", res)
	}
	template := res[testIndexName]
	if template == nil {
		t.Fatalf("expected template %q to be != nil; got: %v", testIndexName, template)
	}
	if len(template.IndexPatterns) != 1 || template.IndexPatterns[0] != "te*" {
		t.Fatalf("expected index settings of %q to be [\"index1\"]; got: %v", testIndexName, template.IndexPatterns)
	}
}
