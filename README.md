[![Go Report Card](https://goreportcard.com/badge/github.com/disaster37/opensearch/v2)](https://goreportcard.com/report/github.com/disaster37/opensearch/v2)
[![Build Status](https://github.com/disaster37/opensearch/workflows/Test/badge.svg)](https://github.com/disaster37/opensearch/actions)
https://goreportcard.com/badge/github.com/disaster37/opensearch/v2
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/disaster37/opensearch?tab=doc)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/disaster37/opensearch/master/LICENSE)
[![codecov](https://codecov.io/gh/disaster37/opensearch/graph/badge.svg?token=4MNPOU84EK)](https://codecov.io/gh/disaster37/opensearch)

# Opensearch

**This is a development branch that is actively being worked on. DO NOT USE IN PRODUCTION! If you want to use stable versions of Opensearch, please use Go modules for the 2.x release (or later) or a dependency manager like [dep](https://github.com/golang/dep) for earlier releases.**

Opensearch is an [Opensearch](http://www.opensearch.org/) client for the
[Go](http://www.golang.org/) programming language.


## Releases

**The release branches (e.g. [`release-branch.v2`](https://github.com/disaster37/opensearch/tree/release-branch.v2))
are actively being worked on and can break at any time.
If you want to use stable versions of Opensearch, please use Go modules.**

Here's the version matrix:

Opensearch version | Opensearch version  | Package URL | Remarks |
----------------------|------------------|-------------|---------|
2.x                   | 2.12.0              | [`gopkg.in/disaster37/opensearch.v2`](https://gopkg.in/disaster37/opensearch.v2) ([source](https://github.com/disaster37/opensearch/tree/release-branch.v2) [doc](http://godoc.org/gopkg.in/disaster37/opensearch.v2)) | Last version

**Example:**

You have installed Opensearchsearch 2.12.0 and want to use Opensearch.
As listed above, you should use Opensearch v2 (code is in `release-branch.v2`).

To use the required version of Opensearch in your application, you
should use [Go modules](https://github.com/golang/go/wiki/Modules)
to manage dependencies. Make sure to use a version such as `2.0.0` or later.

To use Opensearch, import:

```go
import "github.com/disaster37/opensearch/v2"
```

### Opensearch 2.0

Opensearch 2.0 targets Opensearch 2.x.


## Status

We use Opensearch in production since 2024. Opensearch is stable but the API changes
now and then. We strive for API compatibility.
However, Opensearchsearch sometimes introduces
and we sometimes have to adapt.

Having said that, there have been no big API changes that required you
to rewrite your application big time. More often than not it's renaming APIs
and adding/removing features so that Opensearch is in sync with Opensearch cluster.


Opensearch has quite a few features. Most of them are implemented
by Opensearch. I add features and APIs as required. It's straightforward
to implement missing pieces. I'm accepting pull requests :-)

Having said that, I hope you find the project useful.


## Getting Started

The first thing you do is to create a [Client](https://github.com/disaster37/opensearch/blob/master/client.go).
The client connects to Opensearchsearch on `http://127.0.0.1:9200` by default.

You typically create one client for your app. Here's a complete example of
creating a client, creating an index, adding a document, executing a search etc.

An example is available [here](https://disaster37.github.io/opensearch/).

Here's a [link to a complete working example for v2](@todo).

## API Status

### Document APIs

- [x] Index API
- [x] Get API
- [x] Delete API
- [x] Delete By Query API
- [x] Update API
- [x] Update By Query API
- [x] Multi Get API
- [x] Bulk API
- [x] Reindex API
- [x] Term Vectors
- [x] Multi termvectors API

### Search APIs

- [x] Search
- [x] Search Template
- [ ] Multi Search Template
- [x] Search Shards API
- [x] Suggesters
  - [x] Term Suggester
  - [x] Phrase Suggester
  - [x] Completion Suggester
  - [x] Context Suggester
- [x] Multi Search API
- [x] Count API
- [x] Validate API
- [x] Explain API
- [x] Profile API
- [x] Field Capabilities API

### Aggregations

- Metrics Aggregations
  - [x] Avg
  - [x] Cardinality
  - [x] Extended Stats
  - [x] Geo Bounds
  - [x] Geo Centroid
  - [x] Matrix stats
  - [x] Max
  - [x] Median absolute deviation
  - [x] Min
  - [x] Percentile Ranks
  - [x] Percentiles
  - [ ] Scripted Metric
  - [x] Stats
  - [x] Sum
  - [x] Top Hits
  - [x] Value Count
  - [x] Weighted avg
- Bucket Aggregations
  - [x] Adjacency Matrix
  - [x] Auto-interval Date Histogram
  - [x] Children
  - [x] Composite
  - [x] Date Histogram
  - [x] Date Range
  - [x] Diversified Sampler
  - [x] Filter
  - [x] Filters
  - [x] Geo Distance
  - [x] Geohash Grid
  - [x] Geotile grid
  - [x] Global
  - [x] Histogram
  - [x] IP Range
  - [x] Missing
  - [x] Nested
  - [ ] Parent
  - [x] Range
  - [ ] Rare terms
  - [x] Reverse Nested
  - [x] Sampler
  - [x] Significant Terms
  - [x] Significant Text
  - [x] Terms
  - [ ] Variable width histogram
- Pipeline Aggregations
  - [x] Avg Bucket
  - [x] Bucket Script
  - [x] Bucket Selector
  - [x] Bucket Sort
  - [x] Cumulative Sum
  - [x] Derivative
  - [ ] Extended Stats Bucket
  - [x] Max Bucket
  - [x] Min Bucket
  - [x] Moving Average
  - [x] Moving function
  - [x] Percentiles Bucket
  - [x] Serial Differencing
  - [x] Stats Bucket
  - [x] Sum Bucket
- [x] Aggregation Metadata

### Indices APIs

- [x] Create Index
- [x] Delete Index
- [x] Get Index
- [x] Indices Exists
- [x] Open / Close Index
- [x] Shrink Index
- [x] Rollover Index
- [x] Put Mapping
- [x] Get Mapping
- [x] Get Field Mapping
- [x] Types Exists
- [x] Index Aliases
- [x] Update Indices Settings
- [x] Get Settings
- [x] Analyze
  - [x] Explain Analyze
- [x] Index Templates
- [x] Indices Stats
- [x] Indices Segments
- [ ] Indices Recovery
- [ ] Indices Shard Stores
- [x] Clear Cache
- [x] Flush
  - [x] Synced Flush
- [x] Refresh
- [x] Force Merge

### cat APIs

- [X] cat aliases
- [X] cat allocation
- [X] cat count
- [X] cat fielddata
- [X] cat health
- [X] cat indices
- [x] cat master
- [ ] cat nodeattrs
- [ ] cat nodes
- [ ] cat pending tasks
- [ ] cat plugins
- [ ] cat recovery
- [ ] cat repositories
- [ ] cat thread pool
- [ ] cat shards
- [ ] cat segments
- [X] cat snapshots
- [ ] cat templates

### Cluster APIs

- [x] Cluster Health
- [x] Cluster State
- [x] Cluster Stats
- [ ] Pending Cluster Tasks
- [x] Cluster Reroute
- [ ] Cluster Update Settings
- [x] Nodes Stats
- [x] Nodes Info
- [ ] Nodes Feature Usage
- [ ] Remote Cluster Info
- [x] Task Management API
- [ ] Nodes hot_threads
- [ ] Cluster Allocation Explain API

### Rollup APIs (XPack)
- [x] Create Job
- [x] Delete Job
- [x] Get Job
- [x] Start Job
- [x] Stop Job

### Query DSL

- [x] Match All Query
- [x] Inner hits
- Full text queries
  - [x] Match Query
  - [x] Match Boolean Prefix Query
  - [x] Match Phrase Query
  - [x] Match Phrase Prefix Query
  - [x] Multi Match Query
  - [x] Common Terms Query
  - [x] Query String Query
  - [x] Simple Query String Query
  - [x] Combined Fields Query
  - [x] Intervals Query
- Term level queries
  - [x] Term Query
  - [x] Terms Query
  - [x] Terms Set Query
  - [x] Range Query
  - [x] Exists Query
  - [x] Prefix Query
  - [x] Wildcard Query
  - [x] Regexp Query
  - [x] Fuzzy Query
  - [x] Type Query
  - [x] Ids Query
- Compound queries
  - [x] Constant Score Query
  - [x] Bool Query
  - [x] Dis Max Query
  - [x] Function Score Query
  - [x] Boosting Query
- Joining queries
  - [x] Nested Query
  - [x] Has Child Query
  - [x] Has Parent Query
  - [x] Parent Id Query
- Geo queries
  - [ ] GeoShape Query
  - [x] Geo Bounding Box Query
  - [x] Geo Distance Query
  - [x] Geo Polygon Query
- Specialized queries
  - [x] Distance Feature Query
  - [x] More Like This Query
  - [x] Script Query
  - [x] Script Score Query
  - [x] Percolate Query
- Span queries
  - [x] Span Term Query
  - [ ] Span Multi Term Query
  - [x] Span First Query
  - [x] Span Near Query
  - [ ] Span Or Query
  - [ ] Span Not Query
  - [ ] Span Containing Query
  - [ ] Span Within Query
  - [ ] Span Field Masking Query
- [ ] Minimum Should Match
- [ ] Multi Term Query Rewrite

### Modules

- Snapshot and Restore
  - [x] Repositories
  - [x] Snapshot get
  - [x] Snapshot create
  - [x] Snapshot delete
  - [ ] Restore
  - [ ] Snapshot status
  - [ ] Monitoring snapshot/restore status
  - [ ] Stopping currently running snapshot and restore
- Scripting
  - [x] GetScript
  - [x] PutScript
  - [x] DeleteScript

### Sorting

- [x] Sort by score
- [x] Sort by field
- [x] Sort by geo distance
- [x] Sort by script
- [x] Sort by doc

### Security

- Security plugin
  - [x] Internal user
  - [x] Role
  - [x] Role mapping
  - [x] Action group
  - [x] Tenant
  - [x] Distinguished name (DN)
  - [x] Flush cache
  - [x] Security config
  - [x] Security audit

### Index Management State
- ISM plugin
  - [x] Index Management State

### Snapshot Management
- SM plugin
  - [x] Snapshot Management

### Alerting Management
- Alerting plugin
  - [x] Monitor


### Scrolling

Scrolling is supported via a  `ScrollService`. It supports an iterator-like interface.
The `ClearScroll` API is implemented as well.

A pattern for [efficiently scrolling in parallel](https://github.com/disaster37/opensearch/wiki/ScrollParallel)
is described in the [Wiki](https://github.com/disaster37/opensearch/wiki).

## How to contribute

Read [the contribution guidelines](https://github.com/disaster37/opensearch/blob/master/CONTRIBUTING.md).

## Credits

Thanks a lot for the great folks working hard on
[Opensearch](https://www.opensearch.co/products/opensearch)
and
[Go](https://golang.org/).

Opensearch uses portions of the
[uritemplates](https://github.com/jtacoma/uritemplates) library
by Joshua Tacoma,
[backoff](https://github.com/cenkalti/backoff) by Cenk Altı and
[leaktest](https://github.com/fortytw2/leaktest) by Ian Chiles.

## LICENSE

MIT-LICENSE.
