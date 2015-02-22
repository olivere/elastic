// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/olivere/elastic"
)

type Tweet struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

var (
	nodes       = flag.String("nodes", "", "comma-separated list of ES URLs (e.g. 'http://192.168.2.10:9200,http://192.168.2.11:9200')")
	n           = flag.Int("n", 5, "number of goroutines that run searches")
	index       = flag.String("index", "twitter", "name of ES index to use")
	logfile     = flag.String("log", "", "log file")
	tracefile   = flag.String("trace", "", "trace file")
	retries     = flag.Int("retries", 5, "number of retries")
	healthcheck = flag.Int("healthcheck", -1, "healthcheck schedule (in seconds)")
	sniffer     = flag.Int("sniffer", -1, "sniffer schedule (in seconds)")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *nodes == "" {
		log.Fatal("no nodes specified")
	}
	urls := strings.SplitN(*nodes, ",", -1)

	testcase, err := NewTestCase(*index, urls)
	if err != nil {
		log.Fatal(err)
	}

	testcase.SetLogFile(*logfile)
	testcase.SetTraceFile(*tracefile)

	if *retries > 0 {
		testcase.SetMaxRetries(*retries)
	}
	if *healthcheck > 0 {
		testcase.SetHealthcheckSchedule(time.Duration(*healthcheck) * time.Second)
	}
	if *sniffer > 0 {
		testcase.SetSnifferSchedule(time.Duration(*sniffer) * time.Second)
	}

	if err := testcase.Run(*n); err != nil {
		log.Fatal(err)
	}

	select {}
}

type RunInfo struct {
	Success bool
}

type TestCase struct {
	nodes      []string
	client     *elastic.Client
	runs       int64
	failures   int64
	runCh      chan RunInfo
	index      string
	logfile    string
	tracefile  string
	maxRetries int
}

func NewTestCase(index string, nodes []string) (*TestCase, error) {
	if index == "" {
		return nil, errors.New("no index name specified")
	}

	t := &TestCase{index: index, nodes: nodes, runCh: make(chan RunInfo)}

	client, err := elastic.NewClient(http.DefaultClient, t.nodes...)
	if err != nil {
		// Handle error
		return nil, err
	}
	t.client = client

	return t, nil
}

func (t *TestCase) SetIndex(name string) {
	t.index = name
}

func (t *TestCase) SetLogFile(name string) {
	t.logfile = name
}

func (t *TestCase) SetTraceFile(name string) {
	t.tracefile = name
}

func (t *TestCase) SetMaxRetries(n int) {
	t.client.SetMaxRetries(n)
}

func (t *TestCase) SetHealthcheckSchedule(d time.Duration) {
	t.client.SetHealthcheckSchedule(d)
}

func (t *TestCase) SetSnifferSchedule(d time.Duration) {
	t.client.SetSnifferSchedule(d)
}

func (t *TestCase) Run(n int) error {
	if err := t.setup(); err != nil {
		return err
	}

	for i := 1; i < n; i++ {
		go t.search()
	}

	go t.monitor()

	return nil
}

func (t *TestCase) monitor() {
	print := func() {
		fmt.Printf("\033[32m%5d\033[0m; \033[31m%5d\033[0m: %s%s\r", t.runs, t.failures, t.client, "    ")
	}

	for {
		select {
		case run := <-t.runCh:
			atomic.AddInt64(&t.runs, 1)
			if !run.Success {
				atomic.AddInt64(&t.failures, 1)
			}
			print()
		case <-time.After(5 * time.Second):
			// Print stats after some inactivity
			print()
			break
		}
	}
}

func (t *TestCase) setup() error {
	// Log requests to elastic.log
	if t.logfile != "" {
		f, err := os.OpenFile(t.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return err
		}
		t.client.SetLogger(log.New(f, "", log.LstdFlags))
	}

	// Trace request and response details like this
	if t.tracefile != "" {
		f, err := os.OpenFile(t.tracefile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return err
		}
		t.client.SetTracer(log.New(f, "", log.LstdFlags))
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := t.client.IndexExists(t.index).Do()
	if err != nil {
		return err
	}
	if exists {
		deleteIndex, err := t.client.DeleteIndex(t.index).Do()
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			return errors.New("delete index not acknowledged")
		}
	}

	// Create a new index.
	createIndex, err := t.client.CreateIndex(t.index).Do()
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return errors.New("create index not acknowledged")
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	_, err = t.client.Index().
		Index(t.index).
		Type("tweet").
		Id("1").
		BodyJson(tweet1).
		Do()
	if err != nil {
		return err
	}

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	_, err = t.client.Index().
		Index(t.index).
		Type("tweet").
		Id("2").
		BodyString(tweet2).
		Do()
	if err != nil {
		return err
	}

	// Flush to make sure the documents got written.
	_, err = t.client.Flush().Index(t.index).Do()
	if err != nil {
		return err
	}

	return nil
}

func (t *TestCase) search() {
	// Loop forever to check for connection issues
	for {
		// Get tweet with specified ID
		get1, err := t.client.Get().
			Index(t.index).
			Type("tweet").
			Id("1").
			Do()
		if err != nil {
			//failf("Get failed: %v", err)
			t.runCh <- RunInfo{Success: false}
			continue
		}
		if !get1.Found {
			//log.Printf("Document %s not found\n", "1")
			//fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
			t.runCh <- RunInfo{Success: false}
			continue
		}

		// Search with a term query
		termQuery := elastic.NewTermQuery("user", "olivere")
		searchResult, err := t.client.Search().
			Index(t.index).     // search in index t.index
			Query(&termQuery).  // specify the query
			Sort("user", true). // sort by "user" field, ascending
			From(0).Size(10).   // take documents 0-9
			Pretty(true).       // pretty print request and response JSON
			Do()                // execute
		if err != nil {
			//failf("Search failed: %v\n", err)
			t.runCh <- RunInfo{Success: false}
			continue
		}

		// searchResult is of type SearchResult and returns hits, suggestions,
		// and all kinds of other information from Elasticsearch.
		//fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

		// Number of hits
		if searchResult.Hits != nil {
			//fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index

				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var tweet Tweet
				err := json.Unmarshal(*hit.Source, &tweet)
				if err != nil {
					// Deserialization failed
					//failf("Deserialize failed: %v\n", err)
					t.runCh <- RunInfo{Success: false}
					continue
				}

				// Work with tweet
				//fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
			}
		} else {
			// No hits
			//fmt.Print("Found no tweets\n")
		}

		t.runCh <- RunInfo{Success: true}

		// Sleep some time
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
