// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestRuntimeMappingsSource(t *testing.T) {
	rm := RuntimeMappings{
		"day_of_week": map[string]interface{}{
			"type": "keyword",
		},
	}
	src, err := rm.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"day_of_week":{"type":"keyword"}}`
	if want, have := expected, string(data); want != have {
		t.Fatalf("want %s, have %s", want, have)
	}
}

func TestRuntimeMappings(t *testing.T) {
	client := setupTestClient(t) //, SetTraceLog(log.New(os.Stdout, "", 0)))

	ctx := context.Background()
	indexName := testIndexName

	// Create index
	createIndex, err := client.CreateIndex(indexName).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	mapping := `{
		"dynamic": "runtime",
		"properties": {
			"@timestamp": {
				"type":"date"
			}
		},
		"runtime": {
			"day_of_week": {
				"type": "keyword",
				"script": {
					"source": "emit(doc['@timestamp'].value.dayOfWeekEnum.getDisplayName(TextStyle.FULL, Locale.ROOT))"
				}
			}
		}
	}`
	type Doc struct {
		Timestamp time.Time `json:"@timestamp"`
	}
	type DynamicDoc struct {
		Timestamp time.Time `json:"@timestamp"`
		DayOfWeek string    `json:"day_of_week"`
	}

	// Create mapping
	putResp, err := client.PutMapping().
		Index(indexName).
		BodyString(mapping).
		Do(ctx)
	if err != nil {
		t.Fatalf("expected put mapping to succeed; got: %v", err)
	}
	if putResp == nil {
		t.Fatalf("expected put mapping response; got: %v", putResp)
	}
	if !putResp.Acknowledged {
		t.Fatalf("expected put mapping ack; got: %v", putResp.Acknowledged)
	}

	// Add a document
	timestamp := time.Date(2021, 1, 17, 23, 24, 25, 26, time.UTC)
	indexResult, err := client.Index().
		Index(indexName).
		Id("1").
		BodyJson(&Doc{
			Timestamp: timestamp,
		}).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if indexResult == nil {
		t.Errorf("expected result to be != nil; got: %v", indexResult)
	}

	// Execute a search to check for runtime fields
	searchResp, err := client.Search(indexName).
		Query(NewMatchAllQuery()).
		DocvalueFields("@timestamp", "day_of_week").
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if searchResp == nil {
		t.Errorf("expected result to be != nil; got: %v", searchResp)
	}
	if want, have := int64(1), searchResp.TotalHits(); want != have {
		t.Fatalf("expected %d search hits, got %d", want, have)
	}

	// The hit should not have the "day_of_week"
	hit := searchResp.Hits.Hits[0]
	var doc DynamicDoc
	if err := json.Unmarshal(hit.Source, &doc); err != nil {
		t.Fatalf("unable to deserialize hit: %v", err)
	}
	if want, have := timestamp, doc.Timestamp; want != have {
		t.Fatalf("expected timestamp=%v, got %v", want, have)
	}
	if want, have := "", doc.DayOfWeek; want != have {
		t.Fatalf("expected day_of_week=%q, got %q", want, have)
	}

	// The fields should include a "day_of_week" of ["Sunday"]
	dayOfWeekIntfSlice, ok := hit.Fields["day_of_week"].([]interface{})
	if !ok {
		t.Fatalf("expected a slice of strings, got %T", hit.Fields["day_of_week"])
	}
	if want, have := 1, len(dayOfWeekIntfSlice); want != have {
		t.Fatalf("expected a slice of size %d, have %d", want, have)
	}
	dayOfWeek, ok := dayOfWeekIntfSlice[0].(string)
	if !ok {
		t.Fatalf("expected an element of string, got %T", dayOfWeekIntfSlice[0])
	}
	if want, have := "Sunday", dayOfWeek; want != have {
		t.Fatalf("expected day_of_week=%q, have %q", want, have)
	}
}
