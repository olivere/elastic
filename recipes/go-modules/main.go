// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect simply connects to Elasticsearch, but uses Go modules
// as introduced with Go 1.11.
//
// Example
//
//
//     GO111MODULE=on go run main.go -url=http://127.0.0.1:9200 -sniff=false
//
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/olivere/elastic" // <- should end with /v6, but missing due to compatibility reasons
)

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}

	// Just a status message
	fmt.Println("Connection succeeded")

	version, err := client.ElasticsearchVersion(*url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Elasticsearch version %s\n", version)
}
