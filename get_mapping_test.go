package elastic

import (
	"fmt"
	"testing"
)

func TestGetMapping(t *testing.T) {
	client := setupTestClient(t)

	// create an empty index
	createIndex, err := client.CreateIndex(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// get the empty index mapping
	mapping, err := client.GetMapping().Get(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	expectEmptyMap := fmt.Sprintf(`{"%s":{"mappings":{}}}`, testIndexName)
	if mapping != expectEmptyMap {
		t.Fatalf("Expected mapping = %s; got: %s", expectEmptyMap, mapping)
	}

	// create an index with a mapping
	newMap := `{"mappings":{"tweet":{"properties":{"location":{"type":"geo_point"},"tags":{"type":"string"}}}}}`
	expectedFullMap := fmt.Sprintf(`{"%s":%s}`, testIndexName2, newMap)

	createIndex, err = client.CreateIndex(testIndexName2).Body(newMap).Do()
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// get the index mapping
	mapping, err = client.GetMapping().Get(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}

	if mapping != expectedFullMap {
		t.Fatalf("Expected mapping = %s; got: %s", expectedFullMap, mapping)
	}
}
