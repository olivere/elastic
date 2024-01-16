// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect simply connects to Opensearch.
//
// # Example
//
//	connect -url=http://127.0.0.1:9200 -sniff=false
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/olivere/opensearch"
)

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Opensearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	// Create an Opensearch client
	client, err := opensearch.NewClient(opensearch.SetURL(*url), opensearch.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	// Just a status message
	fmt.Println("Connection succeeded")
}
