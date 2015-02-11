package elastic

import (
	"reflect"
	"testing"
)

func TestMappings(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	index1 := "elastic_test_mappings1"
	index2 := "elastic_test_mappings2"

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	expected := map[string]Mappings{
		"elastic_test_mappings1": Mappings{
			Mappings: map[string]TypeMappings{
				"tweet": TypeMappings{
					Properties: map[string]FieldProperty{
						"message": FieldProperty{
							Type: "string",
						},
						"retweets": FieldProperty{
							Type: "long",
						},
						"suggest_field": FieldProperty{
							Type:                       "completion",
							Analyzer:                   "simple",
							Payloads:                   true,
							PreserveSeparators:         true,
							PreservePositionIncrements: true,
							MaxInputLength:             50,
						},
						"tags": FieldProperty{
							Type: "string",
						},
						"user": FieldProperty{
							Type: "string",
						},
						"created": FieldProperty{
							Type:   "date",
							Format: "dateOptionalTime",
						},
						"location": FieldProperty{
							Type: "geo_point",
						},
					},
				},
			},
		},

		"elastic_test_mappings2": Mappings{
			Mappings: map[string]TypeMappings{
				"tweet": TypeMappings{
					Properties: map[string]FieldProperty{
						"location": FieldProperty{
							Type: "geo_point",
						},
						"suggest_field": FieldProperty{
							Type:                       "completion",
							Analyzer:                   "simple",
							Payloads:                   true,
							PreserveSeparators:         true,
							PreservePositionIncrements: true,
							MaxInputLength:             50,
						},
						"tags": FieldProperty{
							Type: "string",
						},
					},
				},
			},
		},
	}

	exists, err := client.IndexExists(index1).Do()

	if err != nil {
		// Handle error
		panic(err)
	}

	if exists {
		_, err = client.DeleteIndex(index1).Do()

		if err != nil {
			t.Fatalf("Failed to delete index '%s'...", index1)
		}
	}

	// Create a new index.
	createIndex, err := client.CreateIndex(index1).Do()
	if err != nil {
		// Handle error
		t.Fatalf("Failed to create index %s: %s", index1, err)
	}

	if createIndex.Acknowledged != true {
		t.Error("Index not ack")
	}

	mappingsOk, err := client.CreateMappings().Index(index1).Type("tweet").Body(expected[index1].Mappings).Do()

	if err != nil {
		t.Fatalf("Failed to create mappings for index '%s': %s", index1, err)
	}

	if mappingsOk["acknowledged"] != true {
		t.Fatalf("Failed to create mappings for index '%s'", index1)
	}

	exists, err = client.IndexExists(index2).Do()

	if err != nil {
		// Handle error
		panic(err)
	}

	if exists {
		_, err = client.DeleteIndex(index2).Do()

		if err != nil {
			t.Fatalf("Failed to delete index '%s'...", index2)
		}
	}

	// Create a new index.
	createIndex, err = client.CreateIndex(index2).Do()
	if err != nil {
		// Handle error
		t.Fatalf("Failed to create index %s: %s", index2, err)
	}

	if createIndex.Acknowledged != true {
		t.Error("Index not ack")
	}

	mappingsOk, err = client.CreateMappings().Index(index2).Type("tweet").Body(expected[index2].Mappings).Do()

	if err != nil {
		t.Fatalf("Failed to create mappings for index '%s': %s", index2, err)
	}

	if mappingsOk["acknowledged"] != true {
		t.Fatalf("Failed to create mappings for index '%s'", index2)
	}

	// Add documents
	_, err = client.Index().Index(index1).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(index1).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(index1).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(index1).Do()
	if err != nil {
		t.Fatal(err)
	}

	// get mappings
	mappings, err := client.Mappings().Indices(index1, index2).Do()

	if err != nil {
		t.Fatal(err)
	}
	if mappings == nil {
		t.Errorf("expected mappings != nil")
	}

	eq := reflect.DeepEqual(expected, mappings)

	if !eq {
		t.Error("Mappings are unequal!")
	}

	_, err = client.DeleteIndex(index1).Do()

	if err != nil {
		t.Fatalf("Failed to delete index '%s'...", index1)
	}

	_, err = client.DeleteIndex(index2).Do()

	if err != nil {
		t.Fatalf("Failed to delete index '%s'...", index1)
	}
}
