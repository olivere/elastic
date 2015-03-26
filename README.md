# Elastic

Elastic is an [Elasticsearch](http://www.elasticsearch.org/) client for [Go](http://www.golang.org/).

[![Build Status](https://travis-ci.org/olivere/elastic.svg?branch=release-branch.v1)](https://travis-ci.org/olivere/elastic)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/gopkg.in/olivere/elastic.v1)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/olivere/elastic/master/LICENSE)


## Releases

**Notice**: This is version 1.0 of Elastic. There is a newer version
available on [https://github.com/olivere/elastic](https://github.com/olivere/elastic).
I encourage anyone to use the newest version.

However, if you want to continue using the 1.0 version, you need to go-get
a new URL and switch your import path. We're using [gopkg.in](http://gokpg.in/) for that.
Here's how to use Elastic version 1:

```sh
$ go get -u gopkg.in/olivere/elastic.v1
```

In your Go code:

```go
import "gopkg.in/olivere/elastic.v1"
```

If you instead use `github.com/olivere/elastic` in your code base, you are
following master. I try to keep master stable, but things might
break now and then.

## Status

We use Elastic in production for more than two years now.
Although Elastic is quite stable from our experience, we don't have
a stable API yet. The reason for this is that Elasticsearch changes quite
often and at a fast pace. At this moment we focus on features, not on a
stable API. Having said that, there have been no huge changes for the last
12 months that required you to rewrite your application big time.
More often than not it's renaming APIs and adding/removing features
so that we are in sync with the Elasticsearch API.

Elastic supports and has been tested in production with
the following Elasticsearch versions: 0.90, 1.0, 1.1, 1.2, 1.3, and 1.4.

Elasticsearch has quite a few features. A lot of them are
not yet implemented in Elastic (see below for details).
I add features and APIs as required. It's straightforward
to implement missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful. Fork it
as you like.

## Usage

The first thing you do is to create a Client. The client takes a http.Client
and (optionally) a list of URLs to the Elasticsearch servers as arguments.
If the list of URLs is empty, http://localhost:9200 is used by default.
You typically create one client for your app.

```go
client, err := elastic.NewClient(http.DefaultClient)
if err != nil {
    // Handle error
}
```

Notice that you can pass your own http.Client implementation here. You can
also pass more than one URL to a client. Elastic pings the URLs periodically
and takes the first to succeed. By doing this periodically, Elastic provides
automatic failover, e.g. when an Elasticsearch server goes down during
updates.

If no Elasticsearch server is available, services will fail when creating
a new request and will return `ErrNoClient`. While this method is not very
sophisticated and might result in timeouts, it is robust enough for our
use cases. Pull requests are welcome.

```go
client, err := elastic.NewClient(http.DefaultClient, "http://1.2.3.4:9200", "http://1.2.3.5:9200")
if err != nil {
    // Handle error
}
```

A Client provides services. The services usually come with a variety of
methods to prepare the query and a `Do` function to execute it against the
Elasticsearch REST interface and return a response. Here is an example
of the IndexExists service that checks if a given index already exists.

```go
exists, err := client.IndexExists("twitter").Do()
if err != nil {
    // Handle error
}
if !exists {
    // Index does not exist yet.
}
```

Look up the documentation for Client to get an idea of the services provided
and what kinds of responses you get when executing the `Do` function of a service.

Here's a longer example:

```go
// Import Elastic
import (
  "github.com/olivere/elastic"
)

// Obtain a client. You can provide your own HTTP client here.
client, err := elastic.NewClient(http.DefaultClient)
if err != nil {
    // Handle error
    panic(err)
}

// Ping the Elasticsearch server to get e.g. the version number
info, code, err := client.Ping().Do()
if err != nil {
    // Handle error
    panic(err)
}
fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

// Getting the ES version number is quite common, so there's a shortcut
esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
if err != nil {
    // Handle error
    panic(err)
}
fmt.Printf("Elasticsearch version %s", esversion)

// Use the IndexExists service to check if a specified index exists.
exists, err := client.IndexExists("twitter").Do()
if err != nil {
    // Handle error
    panic(err)
}
if !exists {
    // Create a new index.
    createIndex, err := client.CreateIndex("twitter").Do()
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
    Type("tweet").
    Id("1").
    BodyJson(tweet1).
    Do()
if err != nil {
    // Handle error
    panic(err)
}
fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

// Index a second tweet (by string)
tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
put2, err := client.Index().
    Index("twitter").
    Type("tweet").
    Id("2").
    BodyString(tweet2).
    Do()
if err != nil {
    // Handle error
    panic(err)
}
fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

// Get tweet with specified ID
get1, err := client.Get().
    Index("twitter").
    Type("tweet").
    Id("1").
    Do()
if err != nil {
    // Handle error
    panic(err)
}
if get1.Found {
    fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
}

// Flush to make sure the documents got written.
_, err = client.Flush().Index("twitter").Do()
if err != nil {
    panic(err)
}

// Search with a term query
termQuery := elastic.NewTermQuery("user", "olivere")
searchResult, err := client.Search().
    Index("twitter").   // search in index "twitter"
    Query(&termQuery).  // specify the query
    Sort("user", true). // sort by "user" field, ascending
    From(0).Size(10).   // take documents 0-9
    Debug(true).        // print request and response to stdout
    Pretty(true).       // pretty print request and response JSON
    Do()                // execute
if err != nil {
    // Handle error
    panic(err)
}

// searchResult is of type SearchResult and returns hits, suggestions,
// and all kinds of other information from Elasticsearch.
fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

// Number of hits
if searchResult.Hits != nil {
    fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

    // Iterate through results
    for _, hit := range searchResult.Hits.Hits {
        // hit.Index contains the name of the index

        // Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
        var t Tweet
        err := json.Unmarshal(*hit.Source, &t)
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


// Update a tweet by the update API of Elasticsearch.
// We just increment the number of retweets.
update, err := client.Update().Index("twitter").Type("tweet").Id("1").
    Script("ctx._source.retweets += num").
    ScriptParams(map[string]interface{}{"num": 1}).
    Upsert(map[string]interface{}{"retweets": 0}).
    Do()
if err != nil {
    // Handle error
    panic(err)
}
fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)

// ...

// Delete an index.
deleteIndex, err := client.DeleteIndex("twitter").Do()
if err != nil {
    // Handle error
    panic(err)
}
if !deleteIndex.Acknowledged {
    // Not acknowledged
}
```

## Installation

Grab the code with `go get github.com/olivere/elastic`.

## API Status

Here's the current API status.

### APIs

- [x] Search (most queries, filters, facets, aggregations etc. are implemented: see below)
- [x] Index
- [x] Get
- [x] Delete
- [x] Delete By Query
- [x] Update
- [x] Multi Get
- [x] Bulk
- [ ] Bulk UDP
- [ ] Term vectors
- [ ] Multi term vectors
- [x] Count
- [ ] Validate
- [ ] Explain
- [x] Search
- [ ] Search shards
- [x] Search template
- [x] Facets (most are implemented, see below)
- [x] Aggregates (most are implemented, see below)
- [x] Multi Search
- [ ] Percolate
- [ ] More like this
- [ ] Benchmark

### Indices

- [x] Create index
- [x] Delete index
- [x] Indices exists
- [x] Open/close index
- [ ] Put mapping
- [ ] Get mapping
- [ ] Get field mapping
- [ ] Types exist
- [ ] Delete mapping
- [x] Index aliases
- [ ] Update indices settings
- [ ] Get settings
- [ ] Analyze
- [ ] Index templates
- [ ] Warmers
- [ ] Status
- [ ] Indices stats
- [ ] Indices segments
- [ ] Indices recovery
- [ ] Clear cache
- [x] Flush
- [x] Refresh
- [x] Optimize

### Snapshot and Restore

- [ ] Snapshot
- [ ] Restore
- [ ] Snapshot status
- [ ] Monitoring snapshot/restore progress
- [ ] Partial restore

### Cat APIs

Not implemented. Those are better suited for operating with Elasticsearch
on the command line.

### Cluster

- [x] Health
- [x] State
- [ ] Stats
- [ ] Pending cluster tasks
- [ ] Cluster reroute
- [ ] Cluster update settings
- [ ] Nodes stats
- [ ] Nodes info
- [ ] Nodes hot_threads
- [ ] Nodes shutdown

### Query DSL

#### Queries

- [x] `match`
- [x] `multi_match`
- [x] `bool`
- [x] `boosting`
- [ ] `common_terms`
- [ ] `constant_score`
- [x] `dis_max`
- [x] `filtered`
- [x] `fuzzy_like_this_query` (`flt`)
- [x] `fuzzy_like_this_field_query` (`flt_field`)
- [x] `function_score`
- [x] `fuzzy`
- [ ] `geo_shape`
- [x] `has_child`
- [x] `has_parent`
- [x] `ids`
- [ ] `indices`
- [x] `match_all`
- [x] `mlt`
- [x] `mlt_field`
- [x] `nested`
- [x] `prefix`
- [x] `query_string`
- [x] `simple_query_string`
- [x] `range`
- [x] `regexp`
- [ ] `span_first`
- [ ] `span_multi_term`
- [ ] `span_near`
- [ ] `span_not`
- [ ] `span_or`
- [ ] `span_term`
- [x] `term`
- [x] `terms`
- [ ] `top_children`
- [x] `wildcard`
- [ ] `minimum_should_match`
- [ ] `multi_term_query_rewrite`
- [x] `template_query`

#### Filters

- [x] `and`
- [x] `bool`
- [x] `exists`
- [ ] `geo_bounding_box`
- [ ] `geo_distance`
- [ ] `geo_distance_range`
- [x] `geo_polygon`
- [ ] `geoshape`
- [ ] `geohash`
- [x] `has_child`
- [x] `has_parent`
- [x] `ids`
- [ ] `indices`
- [x] `limit`
- [x] `match_all`
- [x] `missing`
- [x] `nested`
- [x] `not`
- [x] `or`
- [x] `prefix`
- [x] `query`
- [x] `range`
- [x] `regexp`
- [ ] `script`
- [x] `term`
- [x] `terms`
- [x] `type`

### Facets

- [x] Terms
- [x] Range
- [x] Histogram
- [x] Date Histogram
- [x] Filter
- [x] Query
- [x] Statistical
- [x] Terms Stats
- [x] Geo Distance

### Aggregations

- [x] min
- [x] max
- [x] sum
- [x] avg
- [x] stats
- [x] extended stats
- [x] value count
- [x] percentiles
- [x] percentile ranks
- [x] cardinality
- [x] geo bounds
- [x] top hits
- [ ] scripted metric
- [x] global
- [x] filter
- [x] filters
- [x] missing
- [x] nested
- [x] reverse nested
- [x] children
- [x] terms
- [x] significant terms
- [x] range
- [x] date range
- [x] ipv4 range
- [x] histogram
- [x] date histogram
- [x] geo distance
- [x] geohash grid

### Sorting

- [x] Sort by score
- [x] Sort by field
- [x] Sort by geo distance
- [x] Sort by script

### Scan

Scrolling through documents (e.g. `search_type=scan`) are implemented via
the `Scroll` and `Scan` services.

## How to contribute

Read [the contribution guidelines](https://github.com/olivere/elastic/blob/master/CONTRIBUTING.md).

## Credits

Thanks a lot for the great folks working hard on
[Elasticsearch](http://www.elasticsearch.org/)
and
[Go](http://www.golang.org/).

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.

