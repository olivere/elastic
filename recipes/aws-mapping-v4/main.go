// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect creates an index with a mapping with different data types.
//
// Example
//
//
//     aws-mapping-v4 -url=https://search-xxxxx-yyyyy.eu-central-1.es.amazonaws.com -index=twitter -type=tweet -sniff=false
//
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/olivere/env"

	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
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
					"user":{
						"type":"keyword"
					},
					"message":{
						"type":"text"
					},
					"retweets":{
						"type":"integer"
					},
					"created":{
						"type":"date"
					},
					"attributes":{
						"type":"object"
					}
				}
			}
		}
	}
	`
)

// Tweet is just an example document.
type Tweet struct {
	User     string                 `json:"user"`
	Message  string                 `json:"message"`
	Retweets int                    `json:"retweets"`
	Created  time.Time              `json:"created"`
	Attrs    map[string]interface{} `json:"attributes,omitempty"`
}

func main() {
	var (
		accessKey = flag.String("access-key", env.String("", "AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"), "Access Key ID")
		secretKey = flag.String("secret-key", env.String("", "AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"), "Secret access key")
		url       = flag.String("url", "", "Elasticsearch URL")
		sniff     = flag.Bool("sniff", false, "Enable or disable sniffing")
		trace     = flag.Bool("trace", false, "Enable or disable tracing")
		index     = flag.String("index", "", "Index name")
		region    = flag.String("region", "eu-west-1", "AWS Region name")
	)
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *url == "" {
		log.Fatal("please specify a URL with -url")
	}
	if *index == "" {
		log.Fatal("please specify an index name with -index")
	}
	if *region == "" {
		log.Fatal("please specify an AWS region with -regiom")
	}

	// Create an Elasticsearch client
	signingClient := aws.NewV4SigningClient(credentials.NewStaticCredentials(
		*accessKey,
		*secretKey,
		"",
	), *region)

	// Create an Elasticsearch client
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
		elastic.SetHealthcheck(*sniff),
		elastic.SetHttpClient(signingClient),
	}
	if *trace {
		opts = append(opts, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Check if index already exists. We'll drop it then.
	// Next, we create a fresh index/mapping.
	ctx := context.Background()
	exists, err := client.IndexExists(*index).Pretty(true).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		_, err := client.DeleteIndex(*index).Pretty(true).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = client.CreateIndex(*index).Body(mapping).Pretty(true).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add a tweet
	{
		tweet := Tweet{
			User:     "olivere",
			Message:  "Welcome to Go and Elasticsearch.",
			Retweets: 0,
			Created:  time.Now(),
			Attrs: map[string]interface{}{
				"views": 17,
				"vip":   true,
			},
		}
		_, err := client.Index().
			Index(*index).
			Type("_doc").
			Id("1").
			BodyJson(&tweet).
			Refresh("true").
			Pretty(true).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read the tweet
	{
		doc, err := client.Get().
			Index(*index).
			Type("_doc").
			Id("1").
			Pretty(true).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		var tweet Tweet
		if err = json.Unmarshal(*doc.Source, &tweet); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s at %s: %s (%d retweets)\n",
			tweet.User,
			tweet.Created,
			tweet.Message,
			tweet.Retweets,
		)
		fmt.Printf("  %v\n", tweet.Attrs)
	}
}
