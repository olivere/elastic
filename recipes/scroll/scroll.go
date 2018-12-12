// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Scroll illustrates scrolling through a set of documents.
//
// Example
//
// Scroll through an index called "products".
// Use "_uid" as the default field:
//
//     scroll -index=products -size=100
//
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/olivere/elastic"
)

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		index = flag.String("index", "", "Elasticsearch index name")
		typ   = flag.String("type", "", "Elasticsearch type name")
		size  = flag.Int("size", 100, "Slice of documents to get per scroll")
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

	// Setup a group of goroutines from the excellent errgroup package
	g, ctx := errgroup.WithContext(context.TODO())

	// Hits channel will be sent to from the first set of goroutines and consumed by the second
	type hit struct {
		Slice int
		Hit   elastic.SearchHit
	}
	hitsc := make(chan hit)

	begin := time.Now()

	// Start goroutine for this sliced scroll
	g.Go(func() error {
		defer close(hitsc)

		// Prepare the query
		var query elastic.Query
		if *typ == "" {
			query = elastic.NewMatchAllQuery()
		} else {
			query = elastic.NewTypeQuery(*typ)
		}
		svc := client.Scroll(*index).Query(query)
		for {
			res, err := svc.Do(ctx)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			for _, searchHit := range res.Hits.Hits {
				// Pass the hit to the hits channel, which will be consumed below
				select {
				case hitsc <- hit{Hit: *searchHit}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})

	// Second goroutine will consume the hits sent from the workers in first set of goroutines
	var total uint64
	g.Go(func() error {
		for range hitsc {
			// We simply count the hits here.
			current := atomic.AddUint64(&total, 1)
			sec := int(time.Since(begin).Seconds())
			fmt.Printf("%8d | %02d:%02d\r", current, sec/60, sec%60)
			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Wait until all goroutines are finished
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scrolled through a total of %d documents in %v\n", total, time.Since(begin))
}
