// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Search illustrates how to search using the Point in Time API.
//
// Example
//
// Scroll through an index called "products".
//
//     search -index=products -size=100
//
// If you don't have an index, use the "populate" command to fill one.
//
//     search -index=products populate
//
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		index = flag.String("index", "", "Elasticsearch index name")
		size  = flag.Int("size", 10, "Slice of documents to get per scroll")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		log.Fatal("missing url parameter")
	}
	if *index == "" {
		log.Fatal("missing index parameter")
	}
	if *size <= 0 {
		log.Fatal("size must be greater than zero")
	}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}

	// Use "search ... populate" to populate the index with random data
	if flag.Arg(0) == "populate" {
		if err := populate(client, *index); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Open a Point in Time
	pit, err := client.OpenPointInTime(*index).KeepAlive("2m").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Notice that at this point we could even delete the documents
	// from the index: The Point in Time will make sure that we still
	// get the results when the Point in Time has been created.
	//
	// Notice that you must not pass an index, a routing, or a preference
	// to the Search API: Those values are taken from the Point in Time.
	res, err := client.Search().
		Query(
			// Return random results
			elastic.NewFunctionScoreQuery().AddScoreFunc(elastic.NewRandomFunction()),
		).
		Size(*size).
		PointInTime(
			elastic.NewPointInTime(pit.Id, "2m"),
		).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, hit := range res.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", doc)
	}
}

// populate will fill an example index.
func populate(client *elastic.Client, indexName string) error {
	bulk := client.Bulk().Index(indexName)
	for i := 0; i < 10000; i++ {
		doc := map[string]interface{}{
			"name": fmt.Sprintf("Product %d", i+1),
		}
		bulk = bulk.Add(elastic.NewBulkIndexRequest().
			Id(fmt.Sprint(i)).
			Doc(doc),
		)
		if bulk.NumberOfActions() >= 100 {
			_, err := bulk.Do(context.Background())
			if err != nil {
				return err
			}
			// bulk is reset after Do, so you can reuse it
			// We ignore indexing errors here though!
		}
	}
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
