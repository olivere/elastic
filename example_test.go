// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opensearch_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	opensearch "github.com/disaster37/opensearch/v2"
	"github.com/sirupsen/logrus"
)

type Tweet struct {
	User     string                   `json:"user"`
	Message  string                   `json:"message"`
	Retweets int                      `json:"retweets"`
	Image    string                   `json:"image,omitempty"`
	Created  time.Time                `json:"created,omitempty"`
	Tags     []string                 `json:"tags,omitempty"`
	Location string                   `json:"location,omitempty"`
	Suggest  *opensearch.SuggestField `json:"suggest_field,omitempty"`
}

func Example() {

	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := opensearch.NewClient(opensearch.SetLogger(log), opensearch.SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), opensearch.SetTransport(transport))
	// Trace request and response details like this
	// client, err := opensearch.NewClient(opensearch.SetTraceLog(log.New(os.Stdout, "", 0)))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Opensearch server to get e.g. the version number
	info, code, err := client.Ping("https://127.0.0.1:9200").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Opensearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.OpensearchVersion("https://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Opensearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
                "retweets":{
                    "type":"long"
                },
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}
`
		createIndex, err := client.CreateIndex("twitter").Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("twitter").
		Id("1").
		BodyJson(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Id("2").
		BodyString(tweet2).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Id("1").
		Do(context.Background())
	if err != nil {
		switch {
		case opensearch.IsNotFound(err):
			panic(fmt.Sprintf("Document not found: %v", err))
		case opensearch.IsTimeout(err):
			panic(fmt.Sprintf("Timeout retrieving document: %v", err))
		case opensearch.IsConnErr(err):
			panic(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			panic(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s\n", get1.Id, get1.Version, get1.Index)

	// Refresh to make sure the documents are searchable.
	_, err = client.Refresh().Index("twitter").Do(context.Background())
	if err != nil {
		panic(err)
	}

	// Search with a term query
	termQuery := opensearch.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").        // search in index "twitter"
		Query(termQuery).        // specify the query
		Sort("user", true).      // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Opensearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Update a tweet by the update API of Opensearch.
	// We just increment the number of retweets.
	script := opensearch.NewScript("ctx._source.retweets += params.num").Param("num", 1)
	update, err := client.Update().Index("twitter").Id("1").
		Script(script).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)

	// ...

	// Delete an index.
	deleteIndex, err := client.DeleteIndex("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}

func ExampleNewClient_default() {
	// Obtain a client to the Opensearch instance on https://127.0.0.1:9200.

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client, err := opensearch.NewClient(opensearch.SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), opensearch.SetTransport(transport))
	if err != nil {
		// Handle error
		fmt.Printf("connection failed: %v\n", err)
	} else {
		fmt.Println("connected")
	}
	_ = client
	// Output:
	// connected
}

func ExampleNewClient_cluster() {
	// Obtain a client for an Opensearch cluster of two nodes,
	// running on 10.0.1.1 and 10.0.1.2.
	client, err := opensearch.NewClient(opensearch.SetURL("http://10.0.1.1:9200", "http://10.0.1.2:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}
	_ = client
}

func ExampleNewClient_manyOptions() {
	// Obtain a client for an Opensearch cluster of two nodes,
	// running on 10.0.1.1 and 10.0.1.2. Do not run the sniffer.
	// Set the healthcheck interval to 10s. When requests fail,
	// retry 5 times. Print error messages to os.Stderr and informational
	// messages to os.Stdout.
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	client, err := opensearch.NewClient(
		opensearch.SetURL("http://10.0.1.1:9200", "http://10.0.1.2:9200"),
		opensearch.SetSniff(false),
		opensearch.SetHealthcheckInterval(10*time.Second),
		opensearch.SetMaxRetries(5),
		opensearch.SetLogger(log))
	if err != nil {
		// Handle error
		panic(err)
	}
	_ = client
}

func ExampleIndicesExistsService() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	// Use the IndexExists service to check if the index "twitter" exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if exists {
		// ...
	}
}

func ExampleIndicesCreateService() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	// Create a new index.
	createIndex, err := client.CreateIndex("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}

func ExampleIndicesDeleteService() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	// Delete an index.
	deleteIndex, err := client.DeleteIndex("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}

func ExampleSearchService() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Search with a term query
	termQuery := opensearch.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").        // search in index "twitter"
		Query(termQuery).        // specify the query
		Sort("user", true).      // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Opensearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Number of hits
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
}

func ExampleAggregations() {
	// Get a client to the local Opensearch instance.
	client, err := opensearch.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Create an aggregation for users and a sub-aggregation for a date histogram of tweets (per year).
	timeline := opensearch.NewTermsAggregation().Field("user").Size(10).OrderByCountDesc()
	histogram := opensearch.NewDateHistogramAggregation().Field("created").CalendarInterval("year")
	timeline = timeline.SubAggregation("history", histogram)

	// Search with a term query
	searchResult, err := client.Search().
		Index("twitter").                     // search in index "twitter"
		Query(opensearch.NewMatchAllQuery()). // return all results, but ...
		SearchType("count").                  // ... do not return hits, just the count
		Aggregation("timeline", timeline).    // add our aggregation to the query
		Pretty(true).                         // pretty print request and response JSON
		Do(context.Background())              // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// Access "timeline" aggregate in search result.
	agg, found := searchResult.Aggregations.Terms("timeline")
	if !found {
		log.Fatalf("we should have a terms aggregation called %q", "timeline")
	}
	for _, userBucket := range agg.Buckets {
		// Every bucket should have the user field as key.
		user := userBucket.Key

		// The sub-aggregation history should have the number of tweets per year.
		histogram, found := userBucket.DateHistogram("history")
		if found {
			for _, year := range histogram.Buckets {
				var key string
				if s := year.KeyAsString; s != nil {
					key = *s
				}
				fmt.Printf("user %q has %d tweets in %q\n", user, year.DocCount, key)
			}
		}
	}
}

func ExampleSearchResult() {
	client, err := opensearch.NewClient()
	if err != nil {
		panic(err)
	}

	// Do a search
	searchResult, err := client.Search().Index("twitter").Query(opensearch.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Opensearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a utility function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate hits with full control.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
}

func ExampleClusterHealthService() {
	client, err := opensearch.NewClient()
	if err != nil {
		panic(err)
	}

	// Get cluster health
	res, err := client.ClusterHealth().Index("twitter").Do(context.Background())
	if err != nil {
		panic(err)
	}
	if res == nil {
		panic(err)
	}
	fmt.Printf("Cluster status is %q\n", res.Status)
}

func ExampleClusterHealthService_WaitForStatus() {
	client, err := opensearch.NewClient()
	if err != nil {
		panic(err)
	}

	// Wait for status green
	res, err := client.ClusterHealth().WaitForStatus("green").Timeout("15s").Do(context.Background())
	if err != nil {
		panic(err)
	}
	if res.TimedOut {
		fmt.Printf("time out waiting for cluster status %q\n", "green")
	} else {
		fmt.Printf("cluster status is %q\n", res.Status)
	}
}

func ExampleClusterStateService() {
	client, err := opensearch.NewClient()
	if err != nil {
		panic(err)
	}

	// Get cluster state
	res, err := client.ClusterState().Metric("version").Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Cluster %q has version %d", res.ClusterName, res.Version)
}
