# Elastic

Elastic is an 
[ElasticSearch](http://www.elasticsearch.org/) 
client for
[Google Go](http://www.golang.org/).

## Status

This is a work in progress, not production ready.

ElasticSearch has quite a few features. A lot of them are
not yet implemented in Elastic (see below for details). 
However, it's should be straightforward to implement 
the missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful. Fork it
as you like. Be prepared for structural changes as well.
As I said: This is still a work-in-progress.

Here's a list of what's currently implemented (but
not thoroughly tested).

### Core

* Index
* Delete
* Get
* Search
* Count (not completed)
* Delete By Query
* Bulk

### Indices

* Aliases
* Create Index
* Delete Index
* Indices Exists
* Flush

### Queries

* Bool
* Filtered
* Ids
* Match
* MatchAll
* MultiMatch
* Prefix
* QueryString
* Term

### Filters

* Exists
* Prefix
* Range
* Term
* Type

### Facets

* Terms
* Range
* Histogram
* Date histogram

## Installation

Grab the code with `go get github.com/olivere/elastic`.

## Example code

Find some typical usage scenarios below:

    // Import Elastic
    import (
      "github.com/olivere/elastic"
    )

    // Obtain a client. 
    // You can provide your own HTTP client here.
    client, err := elastic.NewClient(http.DefaultClient)
    
    // Check if a specified index exists.
    exists, err := client.IndexExists("twitter").Do()
    if exists {
        // Index does exist
    }

    // Create a new index.
    createIndex, err := client.CreateIndex("twitter").Do()

    // Index a tweet (using JSON serialization)
    tweet1 := Tweet{User: "olivere", Message: "Take Five"}
    put1, err := client.Index().
        Index("twitter").
        Type("tweet").
        Id("1").
        BodyJson(tweet1).
        Do()
    
    // Index a second tweet (by string)
    tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
    put2, err := client.Index().
        Index("twitter").
        Type("tweet").
        Id("2").
        BodyString(tweet2).
        Do()

    // Get tweet with specified ID
    get1, err := client.Get().
        Index("twitter").
        Type("tweet").
        Id("1").
        Do()
    
    // Search with a term query
    termQuery := elastic.NewTermQuery("user", "olivere")
    termQueryResult, err := client.Search().
        Index("twitter").
        Query(&termQuery).
        Sort("user", true).
        From(0).Size(10).
        Do()

    // TODO examples of other queries/filters

    // Delete an index.
    deleteIndex, err := client.DeleteIndex("twitter").Do()


## Credits

Thanks a lot for the great folks working hard on
[ElasticSearch](http://www.elasticsearch.org/) 
and
[Google Go](http://www.golang.org/).

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.

