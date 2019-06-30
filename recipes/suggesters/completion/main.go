// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect creates an index with a mapping with different data types.
//
// Example
//
//
//     ./completion -url=http://127.0.0.1:9200 -index=cities
//
// For more details and experimentation, take a look at the official
// documentation at https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-suggesters-completion.html.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic"
)

const (
	mapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"_doc":{
				"properties":{
					"name":{
						"type":"keyword"
					},
					"name_suggest":{
						"type":"completion"
					}
				}
			}
		}
	}
	`
)

// City is used in this example as an index document for suggestions.
type City struct {
	Name        string                `json:"name,omitempty"`
	NameSuggest *elastic.SuggestField `json:"name_suggest,omitempty"`
}

var (
	cities = []string{
		"Amsterdam",
		"Athens",
		"Berlin",
		"Barcelona",
		"Brussels",
		"Dublin",
		"Helsinki",
		"Madrid",
		"Lisbon",
		"London",
		"Oslo",
		"Paris",
		"Stockholm",
		"Rome",
		"Valetta",
		"Vienna",
	}
)

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
		index = flag.String("index", "", "Index name")
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
	client, err := elastic.NewClient(
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
		// elastic.SetTraceLog(log.New(os.Stdout, "", 0)), // uncomment to see the wire protocol
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = client

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

	// Add some cities
	for _, name := range cities {
		fmt.Printf("Add %s...\n", name)
		_, err := client.Index().
			Index(*index).
			Type("_doc").
			BodyJson(&City{
				Name:        name,
				NameSuggest: elastic.NewSuggestField(name),
			}).
			Refresh("true").
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Allow suggestions
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a city name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimRight(name, "\n")
		if name == "" {
			break
		}

		res, err := client.Search().
			Index(*index).
			Type("_doc").
			Suggester(
				elastic.NewCompletionSuggester("name_suggestion").
					Field("name_suggest").
					Text(name),
			).
			Size(0). // we dont' want the hits, just the suggestions
			Pretty(true).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		suggestions, found := res.Suggest["name_suggestion"]
		if !found {
			fmt.Printf("%d matches found:\n", 0)
			continue
		}
		for _, suggestion := range suggestions {
			fmt.Printf("%d suggestions found for text %q:\n", len(suggestion.Options), suggestion.Text)
			for _, opt := range suggestion.Options {
				fmt.Printf("* %s\n", opt.Text)

				// The document's source is in opt.Source
				var city City
				if err = json.Unmarshal(*opt.Source, &city); err != nil {
					log.Fatal(err)
				}
				_ = city
			}
		}
	}
}
