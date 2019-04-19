// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// Connect simply connects to Elasticsearch Service on AWS.
//
// Example
//
//     aws-connect-v4 -url=https://search-xxxxx-yyyyy.eu-central-1.es.amazonaws.com
//
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/olivere/env"

	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
)

func main() {
	var (
		accessKey = flag.String("access-key", env.String("", "AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"), "Access Key ID")
		secretKey = flag.String("secret-key", env.String("", "AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"), "Secret access key")
		url       = flag.String("url", "", "Elasticsearch URL")
		sniff     = flag.Bool("sniff", false, "Enable or disable sniffing")
		region    = flag.String("region", "eu-west-1", "AWS Region name")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		log.Fatal("please specify a URL with -url")
	}
	if *accessKey == "" {
		log.Fatal("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if *secretKey == "" {
		log.Fatal("missing -secret-key or AWS_SECRET_KEY environment variable")
	}
	if *region == "" {
		log.Fatal("please specify an AWS region with -region")
	}

	signingClient := aws.NewV4SigningClient(credentials.NewStaticCredentials(
		*accessKey,
		*secretKey,
		"",
	), *region)

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
