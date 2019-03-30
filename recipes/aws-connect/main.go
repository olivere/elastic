// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect simply connects to Elasticsearch Service on AWS.
//
// Example
//
//     aws-connect -url=https://search-xxxxx-yyyyy.eu-central-1.es.amazonaws.com
//
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/olivere/env"
	awsauth "github.com/smartystreets/go-aws-auth"

	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/aws"
)

func main() {
	var (
		accessKey = flag.String("access-key", env.String("", "AWS_ACCESS_KEY"), "Access Key ID")
		secretKey = flag.String("secret-key", env.String("", "AWS_SECRET_KEY"), "Secret access key")
		url       = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff     = flag.Bool("sniff", false, "Enable or disable sniffing")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}
	if *accessKey == "" {
		log.Fatal("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if *secretKey == "" {
		log.Fatal("missing -secret-key or AWS_SECRET_KEY environment variable")
	}

	signingClient := aws.NewV4SigningClient(awsauth.Credentials{
		AccessKeyID:     *accessKey,
		SecretAccessKey: *secretKey,
	})

	// Create an Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(*url),
		elastic.SetSniff(*sniff),
		elastic.SetHealthcheck(*sniff),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	// Just a status message
	fmt.Println("Connection succeeded")
}
