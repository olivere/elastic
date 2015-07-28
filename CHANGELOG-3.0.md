# Elastic 3.0

**This document is a draft!**

Elasticsearch 2.0 comes with some [breaking changes](https://www.elastic.co/guide/en/elasticsearch/reference/master/breaking-changes-2.0.html). You will probably need to upgrade your application and/or rewrite part of it due to those changes.

We use that window of opportunity to also update Elastic (the Go client) from version 2.0 to 3.0. This will not only introduce changes due to the Elasticsearch 2.0 update but also some changes to make Elastic cleaner by removing some old cruft. When rewriting your application anyway, it is a good chance to upgrade not only Elasticsearch but Elastic as well.

So, to summarize:

1. Elastic 2.0 is compatible with Elasticsearch 1.4+ and is still actively maintained.
2. Elastic 3.0 is compatible with Elasticsearch 2.0+ and will soon become the new master branch.

The rest of the document is a list of all changes in Elastic 3.0.

## Pointer types

All types have changed to be pointer types, not value types. This not only is cleaner but also simplifies the API as illustrated by the following example:

Example for Elastic 2.0 (old):

```go
q := elastic.NewMatchAllQuery()
res, err := elastic.Search("one").Query(&q).Do()  // notice the & here
```

Example for Elastic 3.0 (new):

```go
q := elastic.NewMatchAllQuery()
res, err := elastic.Search("one").Query(q).Do()   // no more &
// ... which can be simplified as:
res, err := elastic.Search("one").Query(elastic.NewMatchAllQuery()).Do()
```

## Query/filter merge

One of the biggest changes in Elasticsearch 2.0 is the [merge of queries and filters](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html#_query_filter_merge). In Elasticsearch 1.x, you had a whole range of queries and filters that were basically identical (e.g. `term_query` and `term_filter`).

The practical aspect of the merge is that you can now basically use queries where once you had to use filters instead. For Elastic 3.0 this means: We could remove a whole bunch of files.

Notice that many methods are still named like e.g. `PostFilter`. However, they accept a `Query` now when they used to accept a `Filter` before.

Example for Elastic 2.0 (old):

```go
q := elastic.NewMatchAllQuery()
f := elastic.NewTermFilter("tag", "important")
res, err := elastic.Search().Index("one").Query(&q).PostFilter(f)
```

Example for Elastic 3.0 (new):

```go
q := elastic.NewMatchAllQuery()
f := elastic.NewTermQuery("tag", "important") // it's a query now!
res, err := elastic.Search().Index("one").Query(q).PostFilter(f)
```

## HTTP Status 404

When Elasticsearch does not find an entity or an index, it generally returns HTTP status code 404. In Elastic 2.0 this was a valid result and didn't raise an error from the `Do` functions. This has now changed in Elastic 3.0.

Starting with Elastic 3.0, there are only two types of responses considered successful. First, responses with HTTP status codes [200..299]. Second, HEAD requests which return HTTP status 404. The latter is used by Elasticsearch to e.g. check for existence of indices or documents. All other responses will return an error.

To check for HTTP Status 404 (with non-HEAD requests), e.g. when trying to get or delete a missing document, you can use the [`IsNotFound`]() helper (see below).

The following example illustrates how to check for a missing document in Elastic 2.0 and what has changed in 3.0.

Example for Elastic 2.0 (old):

```go
res, err = client.Get().Index("one").Type("tweet").Id("no-such-id").Do()
if err != nil {
  // Something else went wrong (but 404 is NOT an error in Elastic 2.0)
}
if !res.Found {
	// Document has not been found
}
```

Example for Elastic 3.0 (new):

```go
res, err = client.Get().Index("one").Type("tweet").Id("no-such-id").Do()
if err != nil {
  if elastic.IsNotFound(err) {
    // Document has not been found
  } else {
    // Something else went wrong
  }
}
```

## Errors

Elasticsearch 2.0 returns more information when an error occurs. Elastic 3.0 now reads all this information and makes it accessible by the consumer.

Errors and all its details are now returned in [`Error`](https://github.com/olivere/elastic/blob/master/errors.go#L49).

### Bulk Errors

The error response of a bulk operation used to be a simple string in Elasticsearch 1.x.
In Elasticsearch 2.0, it returns a structured JSON object with a lot more details about the error.
These errors are now captured in an object of type [`ErrorDetails`](https://github.com/olivere/elastic/blob/master/errors.go#L57) which is used in [BulkResponseItem](https://github.com/olivere/elastic/blob/master/bulk.go#L207).

### Removed ErrMissingIndex, ErrMissingType, and ErrMissingId

The specific error types `ErrMissingIndex`, `ErrMissingType`, and `ErrMissingId` have been removed. They were only used by `DeleteService` and are replaced by a generic error message.

## Numeric types

Elastic 3.0 has settled to use `float64` everywhere. It used to be a mix of `float32` and `float64` in Elastic 2.0. E.g. all boostable queries in Elastic 3.0 now have a boost type of `float64` where it used to be `float32`.

## Pluralization

Some services accept zero, one or more indices or types to operate on.
E.g. in the `SearchService` accepts a list of zero, one, or more indices to
search and therefor had a func called `Index(index string)` and a func
called `Indices(indices ...string)`.

Elastic 3.0 now only uses the singular form that, when applicable, accepts a
variadic type. E.g. in the case of the `SearchService`, you now only have
one func with the following signature: `Index(indices ...string)`.

Notice this is only limited to `Index(...)` and `Type(...)`. There are other
services with variadic functions. These have not been changed.

TODO Add example here

## Multiple calls to variadic functions

Some services with variadic functions have cleared the underlying slice when
called while other services just add to the existing slice. This has now been
normalized to always add to the underlying slice.

Example for Elastic 2.0 (old):

```go
// Would only cleared scroll id "two"
// because ScrollId cleared the values when called multiple times
client.ClearScroll().ScrollId("one").ScrollId("two").Do()
```

Example for Elastic 3.0 (new):

```go
// Now (correctly) clears noth scroll id "one" and "two"
// because ScrollId no longer clears the values when called multiple times
client.ClearScroll().ScrollId("one").ScrollId("two").Do()
```

## Ping service requires URL

The `Ping` service raised some issues because it is different from all
other services. If not explicitly given a URL, it always pings 127.0.0.1:9200.

Users expected to ping the cluster, but that is not possible as the cluster
can be a set of many machines: So which machine do we ping then?

To make it more clear, the `Ping` function on the client now requires users
to explicitly set a URL.

## Meta fields

Many of the meta fields e.g. `_parent` or `_routing` are now
[part of the top-level of a document](https://www.elastic.co/guide/en/elasticsearch/reference/master/_meta_fields_returned_under_the_top_level_json_object.html)
and are no longer returned as parts of the `fields` object. We had to change
larger parts of e.g. the `Reindexer` to get it to work seamlessly with Elasticsearch 2.0.

## HasParentQuery / HasChildQuery

`NewHasParentQuery` and `NewHasChildQuery` must now include both parent/child type and query. It is now in line with the Java API.

Example for Elastic 2.0 (old):

```go
allQ := elastic.NewMatchAllQuery()
q := elastic.NewHasChildFilter("tweet").Query(&allQ)
```

Example for Elastic 3.0 (new):

```go
q := elastic.NewHasChildQuery("tweet", elastic.NewMatchAllQuery())
```

## HasPlugin helper

Some of the core functionality of Elasticsearch has now been moved into plugins. E.g. the Delete-by-Query API is [a plugin now]().

You need to make sure to add these plugins to your Elasticsearch installation to still be able to use the `DeleteByQueryService`. You can test this now with the `HasPlugin(name string)` helper in the client.

TODO Implement this first

Example for Elastic 3.0 (new):

```go
err, found := client.HasPlugin("delete-by-query")
if err == nil && found {
	// ... Delete By Query API is available
}
```

## Delete-by-Query API

The Delete-by-Query API is [a plugin now](). It is no longer core part of Elasticsearch.

Elastic 3.0 still contains the `DeleteByQueryService` but it will fail with `ErrPluginNotFound` when the plugin is not installed.

TODO Find reference to Delete-by-Query Plugin (https://github.com/elastic/elasticsearch/pull/11584/files)
TODO Check the example

Example for Elastic 3.0 (new):

```go
_, err := client.DeleteByQuery().Query(elastic.NewTermQuery("client", "1")).Do()
if err == elastic.ErrPluginNotFound {
	// Delete By Query API is not available
}
```

## Common Query -> Common Terms Query

The `CommonQuery` has been renamed to `CommonTermsQuery` to be in line with the [Java API](https://www.elastic.co/guide/en/elasticsearch/reference/master/_java_api.html).

TODO Double-check

## Remove `MoreLikeThis` and `MoreLikeThisField`

The More Like This API and the More Like This Field query have been removed and replaced with the `MoreLikeThisQuery`. This is a result of [this change in Elasticsearch 2.0](https://www.elastic.co/guide/en/elasticsearch/reference/master/_more_like_this.html).

TODO Double-check

## Remove Filtered Query

With the merge of queries and filters, the [filtered query became deprecated](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html). While it is only deprecated and therefore still available in Elasticsearch 2.0, we have decided to remove it from Elastic 3.0. Why? Because we think that when you're already required to rewrite many of your application code, it might be a good chance to get rid of things that are deprecated as well. So you might simply change your filtered query with a boolean query as [described here](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html).

TODO Really remove FilteredQuery?

## Remove FuzzyLikeThis and FuzzyLikeThisField

Both have been [removed from Elasticsearch 2.0 as well](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html).

## Remove LimitFilter

The `limit` filter is [deprecated in Elasticsearch 2.0](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html) and becomes a no-op. Now is a good chance to remove it from your application as well. Use the `terminate_after` parameter in your search [as described here](https://www.elastic.co/guide/en/elasticsearch/reference/master/search-request-body.html) to achieve similar effects.

## Remove `_cache` and `_cache_key` from filters

Both have been [removed from Elasticsearch 2.0 as well](https://www.elastic.co/guide/en/elasticsearch/reference/master/_query_dsl.html).

## Partial fields are gone

Partial fields are [removed in Elasticsearch 2.0](https://www.elastic.co/guide/en/elasticsearch/reference/master/_partial_fields.html) in favor of [source filtering](https://www.elastic.co/guide/en/elasticsearch/reference/master/search-request-source-filtering.html).

## Scripting

A Script type has been added. In Elastic 2.0, there were various places (e.g. aggregations) where you could just add the script as a string, specify the scripting language, add parameters etc. With Elastic 3.0, you should always use the Script type.

Example for Elastic 2.0 (old):

```go
update, err := client.Update().Index("twitter").Type("tweet").Id("1").
	Script("ctx._source.retweets += num").
	ScriptParams(map[string]interface{}{"num": 1}).
	Upsert(map[string]interface{}{"retweets": 0}).
	Do()
```

Example for Elastic 3.0 (new):

```go
update, err := client.Update().Index("twitter").Type("tweet").Id("1").
	Script(elastic.NewScript("ctx._source.retweets += num").Param("num", 1)).
	Upsert(map[string]interface{}{"retweets": 0}).
	Do()
```

## SetBasicAuth helper

You can now Elastic to pass HTTP Basic Auth credentials with each request. In previous versions of Elastic you had to set up your own `http.Transport` to do this. This should make it more convenient to use Elastic in combination with Shield in its [basic setup](https://www.elastic.co/guide/en/shield/current/enable-basic-auth.html).

Example:

```go
client, err := elastic.NewClient(elastic.SetBasicAuth("user", "secret"))
if err != nil {
  t.Fatal(err)
}
```


## Services

### REST API specification

As you might know, Elasticsearch comes with a REST API specification. The specification describes the endpoints in a JSON structure.

Most services in Elastic predated the REST API specification. Well, now they are in line with the specification. All services can be generated by go generate (not 100% automatic though). All services have been reviewed for being up-to-date with the 2.0 specification.

For you, this probably doesn't mean a lot. However, you can now be more confident that Elastic supports all features the REST API specification specifies.

At the same time, the file names of the services are renamed to match the REST API specification naming.

### REST API Test Suite

The REST API specification of Elasticsearch comes along with a test suite that official clients typically use to test for conformance. Up until now, Elastic didn't run this test suite. However, we are in the process of setting up infrastructure and tests to match this suite as well.

This process in not completed though.


