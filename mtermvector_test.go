package elastic

import (
	"testing"
)

func TestMultiTermVectorsBuildURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		Index    string
		Type     string
		Expected string
	}{
		{
			"twitter",
			"",
			"/twitter/_mtermvectors",
		},
		{
			"twitter",
			"tweet",
			"/twitter/tweet/_mtermvectors",
		},
	}

	for _, test := range tests {
		builder := client.MultiTermVectors(test.Index, test.Type)
		path, _, err := builder.buildURL()
		if err != nil {
			t.Fatal(err)
		}
		if path != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, path)
		}
	}
}

func TestMultiTermVectorsWithIds(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Count documents
	count, err := client.Count(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected Count = %d; got %d", 3, count)
	}

	// MultiTermVectors by specifying ID by 1 and 3
	field := "Message"
	res, err := client.MultiTermVectors(testIndexName, "tweet").
		Add(NewMultiTermvectorItem().Index(testIndexName).Type("tweet").Id("1").Fields(field)).
		Add(NewMultiTermvectorItem().Index(testIndexName).Type("tweet").Id("3").Fields(field)).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected to return information and statistics")
	}
	if res.Docs == nil {
		t.Fatal("expected result docs to be != nil; got nil")
	}
	if len(res.Docs) != 2 {
		t.Fatalf("expected to have 2 docs; got %d", len(res.Docs))
	}
}
