// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

// UpdateByQuery illustrates how to update documents that match a specified query.
//
// See
// https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-update-by-query.html
// for details on the Update By Query API in Elasticsearch.
//
// Example
//
//     go run main.go -url=http://127.0.0.1:9200 -index=testindex
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	mapping = `
	{
		"settings":{
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text"
				},
				"retweets":{
					"type":"integer"
				}
			}
		}
	}
	`
)

// Tweet is just an example document.
type Tweet struct {
	ID       string `json:"id"`
	User     string `json:"user"`
	Message  string `json:"message"`
	Retweets int    `json:"retweets"`
}

func main() {
	var (
		url   = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff = flag.Bool("sniff", true, "Enable or disable sniffing")
		index = flag.String("index", "", "Index name")
	)

	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	if *index == "" {
		log.Fatal("please specify an index name -index")
	}

	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff))
	if err != nil {
		log.Fatal(err)
	}

	// Check if index already exists. We'll drop it then.
	// Next, we create a fresh index/mapping.
	ctx := context.Background()
	exists, err := client.IndexExists(*index).Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	if exists {
		_, err := client.DeleteIndex(*index).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = client.CreateIndex(*index).Body(mapping).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add some tweets
	{
		tweets := []Tweet{
			{
				ID:       "1",
				User:     "olivere",
				Retweets: 1,
				Message:  "Welcome to Golang and Elasticsearch.",
			},
			{
				ID:       "2",
				User:     "olivere",
				Retweets: 2,
				Message:  "Another unrelated topic.",
			},
			{
				ID:       "3",
				User:     "someone",
				Retweets: 3,
				Message:  "Another another unrelated topic.",
			},
		}

		for _, tweet := range tweets {
			_, err := client.Index().Index(*index).Id(tweet.ID).BodyJson(&tweet).Refresh("true").Do(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Build query that match all tweets with unrelated topics.
	q := elastic.NewMatchQuery("message", "unrelated")

	// Build update script.
	s := elastic.NewScript("ctx._source.retweets = 0")

	// Execute the script in matched documents
	res, err := client.
		UpdateByQuery().
		Index(*index).
		Query(q).
		Script(s).Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated tweets: %d with %v failures and %d version conflicts in %dms \n",
		res.Updated,
		res.Failures,
		res.VersionConflicts,
		res.Took)

	listTweets(client, *index)

	// It's possible to add params to a script
	params := map[string]interface{}{"retweets": 0, "user": "blank-space"}
	s = elastic.NewScript("ctx._source.user = params.user; ctx._source.retweets = params.retweets;").Params(params)
	res, err = client.
		UpdateByQuery().
		Index(*index).
		Query(q).
		Script(s).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("updated tweets: %d with %v failures and %d version conflicts in %dms \n",
		res.Updated,
		res.Failures,
		res.VersionConflicts,
		res.Took)

	listTweets(client, *index)

	// Using slices config the updating process can be automatic parallelized, also it's possible
	// to define the amount of slices
	//
	// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/docs-update-by-query.html#docs-update-by-query-slice
	// for details.
	params = map[string]interface{}{"message": "unrelated topic"}
	s = elastic.NewScript("ctx._source.message = params.message").Params(params)
	res, err = client.
		UpdateByQuery().
		Index(*index).
		Query(q).
		Script(s).
		Slices("auto").
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("updated tweets: %d with %v failures and %d version conflicts in %dms \n",
		res.Updated,
		res.Failures,
		res.VersionConflicts,
		res.Took)

	listTweets(client, *index)
}

func listTweets(es *elastic.Client, name string) {
	// Add some interval between updates to prevent conflicts
	time.Sleep(2 * time.Second)

	res, err := es.Search().Index(name).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var tweet Tweet

	for _, t := range res.Hits.Hits {
		if err = json.Unmarshal(t.Source, &tweet); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\t%s: %s (%d retweets)\n",
			tweet.User,
			tweet.Message,
			tweet.Retweets,
		)
	}
}
