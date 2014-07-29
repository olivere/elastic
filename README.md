# Elastic

Elastic is an
[ElasticSearch](http://www.elasticsearch.org/)
client for
[Google Go](http://www.golang.org/).

## Status

We use Elastic in production for more than two years now.
Although Elastic is quite stable from our experience, we don't have
a stable API yet. The reason for this is that Elasticsearch changes quite
often and at a fast pace. At this moment we focus on features, not on a
stable API. Having said that, there have been no huge changes for the last
12 months that required you to rewrite your application from scratch.
More often than not it's renaming APIs and adding/removing features
so that we are in sync with the Elasticsearch API.

Elastic supports and has been tested in production with
the following Elasticsearch versions: 0.90, 1.0, 1.1, 1.2, and 1.3.

Elasticsearch has quite a few features. A lot of them are
not yet implemented in Elastic (see below for details).
I add features and APIs as required. It's straightforward
to implement missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful. Fork it
as you like.

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

### APIs

* Search APIs
  - Most queries, filters, facets, aggregations etc. implemented (see below)
* Index (ok)
* Get (ok)
* Delete (ok)
* Update (missing)
* Multi Get (ok)
* Bulk (ok)
* Bulk UDP (missing)
* Delete By Query (missing)
* Term vectors (missing)
* Multi term vectors (missing)
* Count (ok)
* Validate (missing)
* Explain (missing)
* Search (ok)
* Search shards (missing)
* Search template (missing)
* Facets (most are implemented -- see below)
* Aggregates (most are implemented -- see below)
* Multi Search (ok)
* Percolate (missing)
* Delete By Query (missing)
* More like this (missing)
* Benchmark (missing)

### Indices

* Create index (ok)
* Delete index (ok)
* Indices exists (ok)
* Open/close index (missing)
* Put mapping (missing)
* Get mapping (missing)
* Get field mapping (missing)
* Types exist (missing)
* Delete mapping (missing)
* Index aliases (ok)
* Update indices settings (missing)
* Get settings (missing)
* Analyze (missing)
* Index templates (missing)
* Warmers (missing)
* Status (missing)
* Indices stats (missing)
* Indices segments (missing)
* Indices recovery (missing)
* Clear cache (missing)
* Flush (ok)
* Refresh (missing)
* Optimize (missing)

### Snapshot and Restore

* Snapshot (missing)
* Restore (missing)
* Snapshot status (missing)
* Monitoring snapshot/restore progress (missing)
* Partial restore (missing)

### Cat APIs

Not implemented. Those are better suited for operating with Elasticsearch
on the command line.

### Cluster

* Health (missing)
* State (missing)
* Stats (missing)
* Pending cluster tasks (missing)
* Cluster reroute (missing)
* Cluster update settings (missing)
* Nodes stats (missing)
* Nodes info (missing)
* Nodes hot_threads (missing)
* Nodes shutdown (missing)

### Query DSL

#### Queries

* `match` (ok)
* `multi_match` (ok)
* `bool` (ok)
* `boosting` (missing)
* `common_terms` (missing)
* `constant_score` (missing)
* `dis_max` (ok)
* `filtered` (ok)
* `flt` (missing)
* `flt_field` (missing)
* `function_score` (ok)
* `fuzzy` (missing)
* `geo_shape` (missing)
* `has_child` (ok)
* `has_parent` (ok)
* `ids` (ok)
* `indices` (missing)
* `match_all` (ok)
* `mlt` (ok)
* `mlt_field` (ok)
* `nested` (ok)
* `prefix` (ok)
* `query_string` (ok)
* `simple_query_string` (ok)
* `range` (ok)
* `regexp` (missing)
* `span_first` (missing)
* `span_multi_term` (missing)
* `span_near` (missing)
* `span_not` (missing)
* `span_or` (missing)
* `span_term` (missing)
* `term` (ok)
* `terms` (ok)
* `top_children` (missing)
* `wildcard` (missing)
* `minimum_should_match` (missing)
* `multi_term_query_rewrite` (missing)
* `template_query` (missing)

#### Filters

* `and` (ok)
* `bool` (ok)
* `exists` (ok)
* `geo_bounding_box` (missing)
* `geo_distance` (missing)
* `geo_distance_range` (missing)
* `geo_polygon` (ok)
* `geoshape` (missing)
* `geohash` (missing)
* `has_child` (ok)
* `has_parent` (ok)
* `ids` (ok)
* `indices` (missing)
* `limit` (ok)
* `match_all` (ok)
* `missing` (ok)
* `nested` (ok)
* `not` (ok)
* `or` (ok)
* `prefix` (ok)
* `query` (missing)
* `range` (ok)
* `regexp` (missing)
* `script` (missing)
* `term` (ok)
* `terms` (ok)
* `type` (ok)

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

### Aggregations

* min (ok)
* max (ok)
* sum (ok)
* avg (ok)
* stats (ok)
* extended stats (ok)
* value count (ok)
* percentiles (ok)
* percentile ranks (ok)
* cardinality (ok)
* geo bounds (ok)
* top hits (ok)
* global (ok)
* filter (ok)
* missing (ok)
* nested (ok)
* reverse nested (missing)
* terms (ok)
* significant terms (ok)
* range (ok)
* date range (ok)
* ipv4 range (missing)
* histogram (ok)
* date histogram (ok)
* geo distance (ok)
* geohash grid (missing)

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

