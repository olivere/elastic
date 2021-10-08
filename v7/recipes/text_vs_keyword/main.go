// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect illustrates exact matches vs fulltext queries with
// "text" and "keyword" types and Get as well as Term- and Match Query.
//
// Example
//
//     go run main.go -url=http://127.0.0.1:9200 -index=testindex
//     go run main.go -url=http://127.0.0.1:9200 -index=test -trace
//
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

const (
	mapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"title":{
					"type":"text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				}
			}
		}
	}
	`
)

// Doc is just an example document.
type Doc struct {
	ID    string `json:"-"` // we use it as an ID
	User  string `json:"user"`
	Title string `json:"title"`
}

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
		index = flag.String("index", "", "Index name")
		trace = flag.Bool("trace", false, "Enable or disable trace output")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}
	if *index == "" {
		log.Fatal("please specify an index name -index")
	}

	// Create an Elasticsearch client
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
	}
	if *trace {
		options = append(options, elastic.SetTraceLog(
			log.New(os.Stdout, "", 0),
		))
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		log.Fatal(err)
	}

	// Check if index already exists. We'll drop it then.
	// Next, we create a fresh index/mapping.
	ctx := context.Background()
	exists, err := client.IndexExists(*index).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		_, err := client.DeleteIndex(*index).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = client.CreateIndex(*index).Body(mapping).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add some documents

	{
		docs := []Doc{
			{
				ID:    "one",
				User:  "olivere",
				Title: "Go and Elasticsearch.",
			},
			{
				ID:    "two",
				User:  "pepper",
				Title: "Amsterdam is nice.",
			},
			{
				ID:    "three",
				User:  "salt",
				Title: "So is Barcelona.",
			},
		}
		for _, doc := range docs {
			_, err := client.Index().
				Index(*index).
				Id(doc.ID).
				BodyJson(doc).
				Refresh("true").
				Do(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Read a document by ID. This isn't searching, but looking
	// up a document by its unique identifier. Same with MultiGet.
	{
		src, err := client.Get().
			Index(*index).
			Id("one").
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		var doc Doc
		if err = json.Unmarshal(src.Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n",
			doc.User,
			doc.Title,
		)
	}

	// Read a document by exact query on keyword field
	// It will find the document because "user" is of type "keyword"
	// and there will be an exact match on the value of "salt".
	{
		resp, err := client.Search().
			Index(*index).
			Query(elastic.NewTermQuery("user", "salt")).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if want, have := int64(1), resp.TotalHits(); want != have {
			log.Fatalf("expected %d hits, got %d", want, have)
		}
		var doc Doc
		if err = json.Unmarshal(resp.Hits.Hits[0].Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n",
			doc.User,
			doc.Title,
		)
	}

	// Search for documents by "fulltext" query on text field.
	// This works fine because the MatchQuery will match on "Amsterdam".
	// Notice that the full title of the found document is "Amsterdam is nice.",
	// so an exact query on the "title" field for "Amsterdam" wouldn't find anything.
	{
		resp, err := client.Search().
			Index(*index).
			Query(elastic.NewMatchQuery("title", "Amsterdam")).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if want, have := int64(1), resp.TotalHits(); want != have {
			log.Fatalf("expected %d hits, got %d", want, have)
		}
		var doc Doc
		if err = json.Unmarshal(resp.Hits.Hits[0].Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n",
			doc.User,
			doc.Title,
		)
	}

	// Search for documents by exact query on text field. This case will fail,
	// because "title" is a field of type "text", and there is no exact match
	// for the term "Amsterdam" on any "title" field.
	{
		resp, err := client.Search().
			Index(*index).
			Query(elastic.NewTermQuery("title", "Amsterdam")).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if want, have := int64(0), resp.TotalHits(); want != have {
			log.Fatalf("expected %d hits, got %d", want, have)
		}
	}

	// Notice how this works again. We now search for an exact match by
	// using "title.keyword", a multi-field. Notice that we can only find
	// the exact match if we use "Amsterdam is nice.". Using the term
	// "Amsterdam" would yield no result, again. Doing a match query for
	// "Amsterdam is nice." will yield two results BTW, because there is
	// a match for the word "is" in "So is Barcelona." as well (albeit with
	// a significantly lower score).
	{
		resp, err := client.Search().
			Index(*index).
			Query(elastic.NewMatchQuery("title.keyword", "Amsterdam is nice.")).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if want, have := int64(1), resp.TotalHits(); want != have {
			log.Fatalf("expected %d hits, got %d", want, have)
		}
		var doc Doc
		if err = json.Unmarshal(resp.Hits.Hits[0].Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n",
			doc.User,
			doc.Title,
		)
	}
}
