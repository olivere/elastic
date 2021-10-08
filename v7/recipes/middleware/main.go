// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Middleware via HTTP RoundTripper.
//
// Example
//
//     middleware "http://127.0.0.1:9200/test-index?sniff=false&healthcheck=false"
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/olivere/elastic/v7"
)

// CountingTransport will count requests.
type CountingTransport struct {
	N    int64             // number of requests passing this transport
	next http.RoundTripper // next round-tripper or http.DefaultTransport if nil
}

// RoundTrip implements a transport that will count requests.
func (tr *CountingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	atomic.AddInt64(&tr.N, 1)
	if tr.next != nil {
		return tr.next.RoundTrip(r)
	}
	return http.DefaultTransport.RoundTrip(r)
}

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

	tr := &CountingTransport{}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
		elastic.SetHttpClient(&http.Client{
			Transport: tr,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get ES version
	indices, err := client.IndexNames()
	if err != nil {
		log.Fatal(err)
	}
	for _, index := range indices {
		fmt.Println(index)
	}

	// Just a status message
	fmt.Printf("%d requests executed.\n", tr.N)
}
