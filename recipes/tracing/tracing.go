// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Tracing is the same as the bulk_insert recipe, but adds
// OpenTracing support.
//
// Example
//
// Bulk index 100.000 documents into the index "warehouse", type "product",
// committing every set of 1.000 documents.
//
//     ./run-tracer.sh
//     ./tracing -index=warehouse -type=product -n=100000 -bulk-size=1000
//
package main

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegerexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	"golang.org/x/sync/errgroup"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/trace/opentracing"
)

func main() {
	var (
		url      = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		index    = flag.String("index", "", "Elasticsearch index name")
		typ      = flag.String("type", "", "Elasticsearch type name")
		sniff    = flag.Bool("sniff", true, "Enable or disable sniffing")
		n        = flag.Int("n", 0, "Number of documents to bulk insert")
		bulkSize = flag.Int("bulk-size", 0, "Number of documents to collect before committing")
	)
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	if *url == "" {
		log.Fatal("missing url parameter")
	}
	if *index == "" {
		log.Fatal("missing index parameter")
	}
	if *typ == "" {
		log.Fatal("missing type parameter")
	}
	if *n <= 0 {
		log.Fatal("n must be a positive number")
	}
	if *bulkSize <= 0 {
		log.Fatal("bulk-size must be a positive number")
	}

	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
	}

	// Initialize Jaeger tracing
	cfg := jaegerconfig.Configuration{}
	closer, err := cfg.InitGlobalTracer(
		"elastic-tracing",
		jaegerconfig.Logger(jaegerlog.StdLogger),
		jaegerconfig.Metrics(jaegerexpvar.NewFactory(10)),
	)
	if err != nil {
		log.Fatalf("unable to initialize jaeger tracer: %v", err)
	}
	defer closer.Close()

	httpClient := &http.Client{
		Transport: opentracing.NewTransport(),
	}
	opts = append(opts, elastic.SetHttpClient(httpClient))

	// Create an Elasticsearch client
	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Setup a group of goroutines from the excellent errgroup package
	g, ctx := errgroup.WithContext(context.TODO())

	// The first goroutine will emit documents and send it to the second goroutine
	// via the docsc channel.
	// The second Goroutine will simply bulk insert the documents.
	type doc struct {
		ID        string    `json:"id"`
		Timestamp time.Time `json:"@timestamp"`
	}
	docsc := make(chan doc)

	begin := time.Now()

	// Goroutine to create documents
	g.Go(func() error {
		defer close(docsc)

		buf := make([]byte, 32)
		for i := 0; i < *n; i++ {
			// Generate a random ID
			_, err := rand.Read(buf)
			if err != nil {
				return err
			}
			id := base64.URLEncoding.EncodeToString(buf)

			// Construct the document
			d := doc{
				ID:        id,
				Timestamp: time.Now(),
			}

			// Send over to 2nd goroutine, or cancel
			select {
			case docsc <- d:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Second goroutine will consume the documents sent from the first and bulk insert into ES
	var total uint64
	g.Go(func() error {
		bulk := client.Bulk().Index(*index).Type(*typ)
		for d := range docsc {
			// Simple progress
			current := atomic.AddUint64(&total, 1)
			dur := time.Since(begin).Seconds()
			sec := int(dur)
			pps := int64(float64(current) / dur)
			fmt.Printf("%10d | %6d req/s | %02d:%02d\r", current, pps, sec/60, sec%60)

			// Enqueue the document
			bulk.Add(elastic.NewBulkIndexRequest().Id(d.ID).Doc(d))
			if bulk.NumberOfActions() >= *bulkSize {
				// Commit
				res, err := bulk.Do(ctx)
				if err != nil {
					return err
				}
				if res.Errors {
					// Look up the failed documents with res.Failed(), and e.g. recommit
					return errors.New("bulk commit failed")
				}
				// "bulk" is reset after Do, so you can reuse it
			}

			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Commit the final batch before exiting
		if bulk.NumberOfActions() > 0 {
			_, err = bulk.Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Wait until all goroutines are finished
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// Final results
	dur := time.Since(begin).Seconds()
	sec := int(dur)
	pps := int64(float64(total) / dur)
	fmt.Printf("%10d | %6d req/s | %02d:%02d\n", total, pps, sec/60, sec%60)
}
