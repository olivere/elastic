// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	testIndexName      = "elastic-test"
	testIndexName2     = "elastic-test2"
	testIndexName3     = "elastic-test3"
	testIndexName4     = "elastic-test4"
	testIndexName5     = "elastic-test5"
	testIndexNameEmpty = "elastic-test-empty"
	testMapping        = `
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
				"type":"text",
				"store": true,
				"fielddata": true
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
`
	testMappingWithContext = `
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
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion",
				"contexts":[
					{
						"name":"user_name",
						"type":"category"
					}
				]
			}
		}
	}
}
`

	testNoSourceIndexName = "elastic-nosource-test"
	testNoSourceMapping   = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"_source": {
			"enabled": false
		},
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion",
				"contexts":[
					{
						"name":"user_name",
						"type":"category"
					}
				]
			}
		}
	}
}
`

	testJoinIndex   = "elastic-joins"
	testJoinMapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"message":{
					"type":"text"
				},
				"my_join_field": {
					"type": "join",
					"relations": {
						"question": "answer"
					}
				}
			}
		}
	}
`

	testOrderIndex   = "elastic-orders"
	testOrderMapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"article":{
				"type":"text"
			},
			"manufacturer":{
				"type":"keyword"
			},
			"price":{
				"type":"float"
			},
			"time":{
				"type":"date",
				"format": "yyyy-MM-dd"
			}
		}
	}
}
`

	/*
		   	testDoctypeIndex   = "elastic-doctypes"
		   	testDoctypeMapping = `
		   {
		   	"settings":{
		   		"number_of_shards":1,
		   		"number_of_replicas":0
		   	},
		   	"mappings":{
				"properties":{
					"message":{
						"type":"text",
						"store": true,
						"fielddata": true
					}
				}
		   	}
		   }
		   `
	*/

	testQueryIndex   = "elastic-queries"
	testQueryMapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"query": {
				"type":	"percolator"
			}
		}
	}
}
`
)

type tweet struct {
	User     string        `json:"user"`
	Message  string        `json:"message"`
	Retweets int           `json:"retweets"`
	Image    string        `json:"image,omitempty"`
	Created  time.Time     `json:"created,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Location string        `json:"location,omitempty"`
	Suggest  *SuggestField `json:"suggest_field,omitempty"`
}

func (t tweet) String() string {
	return fmt.Sprintf("tweet{User:%q,Message:%q,Retweets:%d}", t.User, t.Message, t.Retweets)
}

type tweetWithID struct {
	User      string        `json:"user"`
	Message   string        `json:"message"`
	Retweets  int           `json:"retweets"`
	Image     string        `json:"image,omitempty"`
	Created   time.Time     `json:"created,omitempty"`
	Tags      []string      `json:"tags,omitempty"`
	Location  string        `json:"location,omitempty"`
	Suggest   *SuggestField `json:"suggest_field,omitempty"`
	ElasticID string
}

func (t tweetWithID) String() string {
	return fmt.Sprintf("tweet{User:%q,Message:%q,Retweets:%d}", t.User, t.Message, t.Retweets)
}

func (t tweetWithID) SetID(ID string) {
	t.ElasticID = ID
}

type joinDoc struct {
	Message   string      `json:"message"`
	JoinField interface{} `json:"my_join_field,omitempty"`
}

type joinField struct {
	Name   string `json:"name"`
	Parent string `json:"parent,omitempty"`
}

type order struct {
	Article      string  `json:"article"`
	Manufacturer string  `json:"manufacturer"`
	Price        float64 `json:"price"`
	Time         string  `json:"time,omitempty"`
}

func (o order) String() string {
	return fmt.Sprintf("order{Article:%q,Manufacturer:%q,Price:%v,Time:%v}", o.Article, o.Manufacturer, o.Price, o.Time)
}

// doctype is required for Percolate tests.
type doctype struct {
	Message string `json:"message"`
}

func isCI() bool {
	return os.Getenv("TRAVIS") != "" || os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

type logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fail()
	FailNow()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func boolPtr(b bool) *bool { return &b }

// strictDecoder returns an error if any JSON fields aren't decoded.
type strictDecoder struct{}

func (d *strictDecoder) Decode(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

var (
	logDeprecations = flag.String("deprecations", "off", "log or fail on deprecation warnings")
	logTypesRemoval = flag.Bool("types-removal", false, "log deprecation warnings regarding types removal")
	strict          = flag.Bool("strict-decoder", false, "treat missing unknown fields in response as errors")
	noSniff         = flag.Bool("no-sniff", false, "allows to disable sniffing globally")
	noHealthcheck   = flag.Bool("no-healthcheck", false, "allows to disable healthchecks globally")
)

func setupTestClient(t logger, options ...ClientOptionFunc) (client *Client) {
	var err error

	if *noSniff {
		options = append(options, SetSniff(false))
	}
	if *noHealthcheck {
		options = append(options, SetHealthcheck(false))
	}

	client, err = NewClient(options...)
	if err != nil {
		t.Fatal(err)
	}

	// Use strict JSON decoder (unless a specific decoder has been specified already)
	if *strict {
		if client.decoder == nil {
			client.decoder = &strictDecoder{}
		} else if _, ok := client.decoder.(*DefaultDecoder); ok {
			client.decoder = &strictDecoder{}
		}
	}

	// Log deprecations during tests
	if loglevel := *logDeprecations; loglevel != "off" {
		client.deprecationlog = func(req *http.Request, res *http.Response) {
			for _, warning := range res.Header["Warning"] {
				if !*logTypesRemoval && strings.Contains(warning, "[types removal]") {
					continue
				}
				switch loglevel {
				default:
					t.Logf("[%s] Deprecation warning: %s", req.URL, warning)
				case "fail", "error":
					t.Errorf("[%s] Deprecation warning: %s", req.URL, warning)
				}
			}
		}
	}

	client.DeleteIndex(testIndexName).Do(context.TODO())
	client.DeleteIndex(testIndexName2).Do(context.TODO())
	client.DeleteIndex(testIndexName3).Do(context.TODO())
	client.DeleteIndex(testIndexName4).Do(context.TODO())
	client.DeleteIndex(testIndexName5).Do(context.TODO())
	client.DeleteIndex(testIndexNameEmpty).Do(context.TODO())
	client.DeleteIndex(testOrderIndex).Do(context.TODO())
	client.DeleteIndex(testNoSourceIndexName).Do(context.TODO())
	//client.DeleteIndex(testDoctypeIndex).Do(context.TODO())
	client.DeleteIndex(testQueryIndex).Do(context.TODO())
	client.DeleteIndex(testJoinIndex).Do(context.TODO())

	return client
}

func setupTestClientAndCreateIndex(t logger, options ...ClientOptionFunc) *Client {
	client := setupTestClient(t, options...)

	// Create index
	createIndex, err := client.CreateIndex(testIndexName).Body(testMapping).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex)
	}

	// Create second index
	createIndex2, err := client.CreateIndex(testIndexName2).Body(testMappingWithContext).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createIndex2 == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndex2)
	}

	// Create no source index
	createNoSourceIndex, err := client.CreateIndex(testNoSourceIndexName).Body(testNoSourceMapping).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createNoSourceIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createNoSourceIndex)
	}

	// Create order index
	createOrderIndex, err := client.CreateIndex(testOrderIndex).Body(testOrderMapping).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createOrderIndex == nil {
		t.Errorf("expected result to be != nil; got: %v", createOrderIndex)
	}

	// Create empty index
	createIndexEmpty, err := client.CreateIndex(testIndexNameEmpty).Body(testMapping).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if createIndexEmpty == nil {
		t.Errorf("expected result to be != nil; got: %v", createIndexEmpty)
	}

	return client
}

func setupTestClientAndCreateIndexAndLog(t logger, options ...ClientOptionFunc) *Client {
	return setupTestClientAndCreateIndex(t, SetTraceLog(log.New(os.Stdout, "", 0)))
}

var _ = setupTestClientAndCreateIndexAndLog // remove unused warning in staticcheck

func setupTestClientAndCreateIndexAndAddDocs(t logger, options ...ClientOptionFunc) *Client {
	client := setupTestClientAndCreateIndex(t, options...)

	// Add tweets
	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 108, Tags: []string{"golang", "elasticsearch"}}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 0, Tags: []string{"golang"}}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 12, Tags: []string{"sports", "cycling"}}

	_, err := client.Index().Index(testIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Id("3").Routing("someroutingkey").BodyJson(&tweet3).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	// Add orders
	var orders []order
	orders = append(orders, order{Article: "Apple MacBook", Manufacturer: "Apple", Price: 1290, Time: "2015-01-18"})
	orders = append(orders, order{Article: "Paper", Manufacturer: "Canon", Price: 100, Time: "2015-03-01"})
	orders = append(orders, order{Article: "Apple iPad", Manufacturer: "Apple", Price: 499, Time: "2015-04-12"})
	orders = append(orders, order{Article: "Dell XPS 13", Manufacturer: "Dell", Price: 1600, Time: "2015-04-18"})
	orders = append(orders, order{Article: "Apple Watch", Manufacturer: "Apple", Price: 349, Time: "2015-04-29"})
	orders = append(orders, order{Article: "Samsung TV", Manufacturer: "Samsung", Price: 790, Time: "2015-05-03"})
	orders = append(orders, order{Article: "Hoodie", Manufacturer: "h&m", Price: 49, Time: "2015-06-03"})
	orders = append(orders, order{Article: "T-Shirt", Manufacturer: "h&m", Price: 19, Time: "2015-06-18"})
	for i, o := range orders {
		id := fmt.Sprintf("%d", i)
		_, err = client.Index().Index(testOrderIndex).Id(id).BodyJson(&o).Do(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
	}

	// Refresh
	_, err = client.Refresh().Index(testIndexName, testOrderIndex).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupTestClientAndCreateIndexAndAddDocsNoSource(t logger, options ...ClientOptionFunc) *Client {
	client := setupTestClientAndCreateIndex(t, options...)

	// Add tweets
	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}

	_, err := client.Index().Index(testNoSourceIndexName).Id("1").BodyJson(&tweet1).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testNoSourceIndexName).Id("2").BodyJson(&tweet2).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	// Refresh
	_, err = client.Refresh().Index(testNoSourceIndexName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func setupTestClientForXpackSecurity(t logger) (client *Client) {
	var err error
	// Set URL and Auth to use the platinum ES cluster
	options := []ClientOptionFunc{SetURL("http://127.0.0.1:9210"), SetBasicAuth("elastic", "elastic")}

	client, err = NewClient(options...)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type lexicographically struct {
	strings []string
}

func (l lexicographically) Len() int {
	return len(l.strings)
}

func (l lexicographically) Less(i, j int) bool {
	return l.strings[i] < l.strings[j]
}

func (l lexicographically) Swap(i, j int) {
	l.strings[i], l.strings[j] = l.strings[j], l.strings[i]
}
