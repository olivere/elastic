// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Seamlessly connect to an Elasticsearch Service on AWS.
//
// Example
//
//     aws-es-client -domain-name=escluster1 -index=tweets -trace=false
//
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"

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
		domainName = flag.String("domain-name", "", "AWS Elasticsearch Service Domain Name")
		index      = flag.String("index", "", "Index name")
		trace      = flag.Bool("trace", false, "Enable trace logging")
	)
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *domainName == "" {
		log.Fatal("please specify an AWS Elasticsearch Service Domain Name with -domain-name")
	}
	if *index == "" {
		log.Fatal("please specify an index name with -index")
	}

	client, err := ConnectToAWS(context.Background(), *domainName, *trace)
	if err != nil {
		log.Fatal(err)
	}

	// Just a status message
	fmt.Println("Connection succeeded")

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
			Id("1").
			Pretty(true).
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		var tweet Tweet
		if err = json.Unmarshal(doc.Source, &tweet); err != nil {
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

// ConnectToAWS creates an elastic.Client that connects to the ES cluster
// specified by esDomainName.
//
// It creates an AWS session first, by using different approaches, as
// documented by the AWS SDK for Go. Notice that for AWS Elasticsearch,
// we also use the region configured in the session.
//
// Next, it creates an ElasticsearchService instance to access the
// configuration settings of the cluster specified by esDomainName.
// For ConnectToAWS, we only use it to lookup the URL endpoint from the
// configuration.
//
// Finally, we configure all settings to be used with AWS ES. That is:
// * Disable sniffing
// * Disable health checks
// * Close idle connections when a dead node is found
// * Use a HTTP transport that automatically signs HTTP requests
// * (optionally) Trace output to stdout
func ConnectToAWS(ctx context.Context, esDomainName string, trace bool) (*elastic.Client, error) {
	// Create a new AWS session
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	// Create a new ElasticsearchService instance to dynamically retrieve
	// the AWS ES endpoint URL
	svc := elasticsearchservice.New(sess)

	// See https://docs.aws.amazon.com/sdk-for-go/api/service/elasticsearchservice/#ElasticsearchDomainStatus
	out, err := svc.DescribeElasticsearchDomain(&elasticsearchservice.DescribeElasticsearchDomainInput{
		DomainName: &esDomainName,
	})
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", out.DomainStatus)
	// fmt.Printf("AWS Endpoint: %s\n", *out.DomainStatus.Endpoint)

	// Configure the AWS ES Endpoint URL from the ES Domain settings
	awsESEndpoint := &url.URL{
		Host: *out.DomainStatus.Endpoint, // e.g. search-<random-string>.<region>.es.amazonaws.com
	}
	if *out.DomainStatus.DomainEndpointOptions.EnforceHTTPS {
		awsESEndpoint.Scheme = "https"
	} else {
		awsESEndpoint.Scheme = "http"
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
		elastic.SetSniff(false),                // do not sniff with AWS ES
		elastic.SetHealthcheck(false),          // do not perform healthchecks with AWS ES
		elastic.SetCloseIdleConnections(true),  // close idle connections when dead nodes are found
		elastic.SetURL(awsESEndpoint.String()), // use the dynamically retrieved endpoint URL
		elastic.SetHttpClient(httpClient),      // use a HTTP client that does the signing
	}
	if trace {
		// Optional: Trace output
		options = append(options, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}

	// Create a client configured for using with AWS ES
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
