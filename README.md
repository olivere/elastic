# Elastic

Elastic is an
[ElasticSearch](http://www.elasticsearch.org/)
client for
[Google Go](http://www.golang.org/).

## Status

We use Elastic in production for more than a year.
The reason it doesn't have the 1.0 version tag is
that it's incomplete.

ElasticSearch has quite a few features. A lot of them are
not yet implemented in Elastic (see below for details).
I add features and APIs as required. It's straightforward
to implement missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful. Fork it
as you like. There might be some structural changes as well.
As I said: It's not 1.0 yet.

## Usage

Show, don't tell:

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

If you have more than one Elasticsearch server, you can specify
all their URLs when creating a client:

    // Obtain a client that uses two servers.
    client, err := elastic.NewClient(http.DefaultClient,
        "http://127.0.0.1:9200", "http://127.0.0.2:9200")

The client will ping the Elasticsearch servers periodically and
check if they are available. The first that is available will be
used for subsequent requests. In case no Elasticsearch server is
available, creating a new request (via `NewRequest`) will fail
with error `ErrNoClient`. While this method is not very sophisticated
and might result in timeouts, it is robust enough for our use cases.
Pull requests are welcome.

## Installation

Grab the code with `go get github.com/olivere/elastic`.

## API Status

Here's the current API status.

### Core

* Index (ok)
* Get (ok)
* Delete (ok)
* Update (missing)
* Multi Get (missing)
* Search (ok)
* Multi Search (missing)
* Percolate (missing)
* Bulk (ok)
* Bulk UDP (missing)
* Count (ok)
* Delete By Query (missing)
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
* Cluster Reroute (missing)

### Search

* Request Body (ok)
* URI Request (missing)
* Query (ok)
* Filter (ok)
* From/Size (ok)
* Indices/Types (ok)
* Sort (ok)
* Rescore (missing)
* Term Suggest (ok)
* Phrase Suggest (ok)
* Completion Suggest (incomplete)
* Highlighting (missing)
* Fields (missing)
* Script Fields (missing)
* Preference (ok)
* Facets (ok)
* Named Filters (ok)
* Search Type (ok)
* Index Boost (missing)
* Scroll (ok)
* Explain (ok)
* Version (ok)
* Min Score (ok)

### Queries

* `match` (ok)
* `multi_match` (ok)
* `bool` (ok)
* `boosting` (missing)
* `ids` (ok)
* `custom_score` (ok)
* `custom_filters_score` (ok)
* `custom_boost_factor` (missing)
* `constant_score` (missing)
* `dis_max` (ok)
* `field` (missing)
* `filtered` (ok)
* `flt` (missing)
* `flt_field` (missing)
* `fuzzy` (missing)
* `has_child` (ok)
* `has_parent` (ok)
* `match_all` (ok)
* `mlt` (ok)
* `mlt_field` (ok)
* `prefix` (ok)
* `query_string` (ok)
* `simple_query_string` (ok)
* `range` (ok)
* `regexp` (missing)
* `span_first` (missing)
* `span_multi` (missing)
* `span_near` (missing)
* `span_not` (missing)
* `span_or` (missing)
* `span_term` (missing)
* `term` (ok)
* `terms` (ok)
* `common` (ok)
* `top_children` (missing)
* `wildcard` (missing)
* `nested` (ok)
* `custom_filters_score` (ok)
* `indices` (missing)
* `text` (missing)
* `geo_shape` (missing)

### Filters

* `and` (ok)
* `bool` (ok)
* `exists` (ok)
* `ids` (ok)
* `limit` (ok)
* `type` (ok)
* `geo_bbox` (missing)
* `geo_distance` (missing)
* `geo_distance_range` (missing)
* `geo_polygon` (missing)
* `geo_shape` (missing)
* `has_child` (ok)
* `has_parent` (ok)
* `match_all` (ok)
* `missing` (missing)
* `not` (ok)
* `numeric_range` (missing)
* `or` (ok)
* `prefix` (ok)
* `query` (missing)
* `range` (ok)
* `regexp` (missing)
* `script` (missing)
* `term` (ok)
* `terms` (ok)
* `nested` (ok)

### Facets

* Terms (ok)
* Range (ok)
* Histogram (ok)
* Date Histogram (ok)
* Filter (ok)
* Query (ok)
* Statistical (ok)
* Terms Stats (ok)
* Geo Distance (ok)

### Scan

Scrolling through documents (via `search_type=scan`) is implemented.

## Credits

Thanks a lot for the great folks working hard on
[ElasticSearch](http://www.elasticsearch.org/)
and
[Google Go](http://www.golang.org/).

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.

