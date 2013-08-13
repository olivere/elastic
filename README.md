# Elastic

Elastic is an
[ElasticSearch](http://www.elasticsearch.org/)
client for
[Google Go](http://www.golang.org/).

## Status

This is a work in progress. Although it's not a 1.0 version,
we use it in production.

ElasticSearch has quite a few features. A lot of them are
not yet implemented in Elastic (see below for details).
However, it's should be straightforward to implement
the missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful. Fork it
as you like. Be prepared for structural changes as well.
As I said: This is still a work-in-progress.

Here's a list of the API status.

### Core

* Index (ok)
* Delete (ok)
* Get (ok)
* Multi Get (missing)
* Search (ok)
* Multi Search (missing)
* Percolate (missing)
* Bulk (ok)
* Bulk UDP (missing)
* Count (ok)
* Delete By Query (ok)
* More like this (missing)
* Validate (missing)
* Explain (missing)

### Indices

* Aliases (ok)
* Analyze (missing)
* Create Index (ok)
* Delete Index (ok)
* Open/Close Index (missing)
* Get Settings (missing)
* Get Mapping (missing)
* Put Mapping (missing)
* Delete Mapping (missing)
* Refresh (missing)
* Optimize (missing)
* Flush (ok)
* Snapshot (missing)
* Update Settings (missing)
* Templates (missing)
* Warmers (missing)
* Stats (missing)
* Status (missing)
* Segments (missing)
* Clear Cache (missing)
* Indices Exists (ok)
* Types Exists (missing)

### Cluster

* Health (missing)
* State (missing)
* Update Settings (missing)
* Nodes Info (missing)
* Nodes Stats (missing)
* Nodes Shutdown (missing)
* Nodes Hot Threads (missing)
* Cluster reroute (missing)

### Queries

* `match` (ok)
* `multi_match` (ok)
* `bool` (ok)
* `boosting` (missing)
* `ids` (ok)
* `custom_score` (missing)
* `custom_boost_factor` (missing)
* `constant_score` (missing)
* `dis_max` (ok)
* `field` (missing)
* `filtered` (ok)
* `flt` (missing)
* `flt_field` (missing)
* `fuzzy` (missing)
* `has_child` (missing)
* `has_parent` (missing)
* `match_all` (ok)
* `mlt` (missing)
* `mlt_field` (missing)
* `prefix` (ok)
* `query_string` (ok)
* `range` (missing)
* `regexp` (missing)
* `span_first` (missing)
* `span_multi` (missing)
* `span_near` (missing)
* `span_not` (missing)
* `span_or` (missing)
* `span_term` (missing)
* `term` (ok)
* `terms` (missing)
* `common` (ok)
* `top_children` (missing)
* `wildcard` (missing)
* `nested` (ok)
* `custom_filters_score` (ok)
* `indices` (missing)
* `text` (missing)
* `geo_shape` (missing)

### Filters

* `and` (missing)
* `bool` (missing)
* `exists` (ok)
* `ids` (missing)
* `limit` (missing)
* `type` (ok)
* `geo_bbox` (missing)
* `geo_distance` (missing)
* `geo_distance_range` (missing)
* `geo_polygon` (missing)
* `geo_shape` (missing)
* `has_child` (missing)
* `has_parent` (missing)
* `match_all` (missing)
* `missing` (missing)
* `not` (missing)
* `numeric_range` (missing)
* `or` (missing)
* `prefix` (ok)
* `query` (missing)
* `range` (ok)
* `regexp` (missing)
* `script` (missing)
* `term` (ok)
* `terms` (ok)
* `nested` (missing)

### Facets

* Terms (ok)
* Range (ok)
* Histogram (ok)
* Date Histogram (ok)
* Filter (missing)
* Query (ok)
* Statistical (missing)
* Terms Stats (missing)
* Geo Distance (missing)

### Scan

Scrolling through documents (via `search_type=scan`) is implemented.

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

