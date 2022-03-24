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
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/olivere/elastic/v7"
	elasticawsv4 "github.com/olivere/elastic/v7/aws/v4"
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
		url   = flag.String("url", "", "AWS ES Endpoint URL")
		index = flag.String("index", "", "Elasticsearch index name")
		loop  = flag.Bool("loop", false, "Run in an endless loop")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		log.Fatal("please specify an AWS ES Endpoint URL with -url")
	}
	if *index == "" {
		log.Fatal("please specify an index name with -index")
	}

	// Create a pre-configured client to connect to AWS by the given endpoint
	client, err := ConnectToAWS(context.Background(), *url)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Check if index already exists. We'll drop it then.
		// Next, we create a fresh index/mapping.
		ctx := context.Background()
		exists, err := client.IndexExists(*index).Pretty(true).Do(ctx)
		if err != nil {
			if !*loop {
				log.Fatal(err)
			}
			log.Print(err)
			continue
		}
		if exists {
			_, err := client.DeleteIndex(*index).Pretty(true).Do(ctx)
			if err != nil {
				if !*loop {
					log.Fatal(err)
				}
				log.Print(err)
				continue
			}
		}
		_, err = client.CreateIndex(*index).Body(mapping).Pretty(true).Do(ctx)
		if err != nil {
			if !*loop {
				log.Fatal(err)
			}
			log.Print(err)
			continue
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
				Id("1").
				BodyJson(&tweet).
				Refresh("true").
				Pretty(true).
				Do(context.TODO())
			if err != nil {
				if !*loop {
					log.Fatal(err)
				}
				log.Print(err)
				continue
			}
		}

		// Read the tweet
		{
			doc, err := client.Get().
				Index(*index).
				Id("1").
				Pretty(true).
				Do(context.TODO())
			if err != nil {
				if !*loop {
					log.Fatal(err)
				}
				log.Print(err)
				continue
			}
			var tweet Tweet
			if err = json.Unmarshal(doc.Source, &tweet); err != nil {
				if !*loop {
					log.Fatal(err)
				}
				log.Print(err)
			}
			fmt.Printf("%s at %s: %s (%d retweets)\n",
				tweet.User,
				tweet.Created,
				tweet.Message,
				tweet.Retweets,
			)
			fmt.Printf("  %v\n", tweet.Attrs)
		}

		if !*loop {
			break
		}
	}
}

// ConnectToAWS creates an elastic.Client that connects to the ES cluster
// specified by given URL endpoint.
//
// ConnectToAWS ensures we configure all settings to properly use AWS ES with
// this client, e.g.:
// * Disable sniffing
// * Disable health checks
// * Close idle connections when a dead node is found
// * Use a HTTP transport to automatically sign HTTP requests
func ConnectToAWS(ctx context.Context, url string) (*elastic.Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	// We need to sign HTTP requests with AWS
	httpClient := elasticawsv4.NewV4SigningClientWithOptions(
		elasticawsv4.WithCredentials(sess.Config.Credentials),
		elasticawsv4.WithSigner(v4.NewSigner(sess.Config.Credentials, func(s *v4.Signer) {
			s.DisableURIPathEscaping = true
		})),
		elasticawsv4.WithRegion(*sess.Config.Region), // use the AWS region from the session
	)
	options := []elastic.ClientOptionFunc{
		elastic.SetSniff(false),               // do not sniff with AWS ES
		elastic.SetHealthcheck(false),         // do not perform healthchecks with AWS ES
		elastic.SetCloseIdleConnections(true), // close idle connections when dead nodes are found
		elastic.SetURL(url),
		elastic.SetHttpClient(httpClient), // use a HTTP client that does the signing
	}

	// Create a client configured for using with AWS ES
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
