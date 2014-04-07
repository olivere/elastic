// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	_ "encoding/json"
	_ "net/http"
	"testing"
	"time"
)

func TestSearchAggregates(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{
		User:     "olivere",
		Retweets: 108,
		Message:  "Welcome to Golang and ElasticSearch.",
		Image:    "http://golang.org/doc/gopher/gophercolor.png",
		Created:  time.Date(2012, 12, 12, 17, 38, 34, 0, time.UTC),
	}
	tweet2 := tweet{
		User:     "olivere",
		Retweets: 0,
		Message:  "Another unrelated topic.",
		Created:  time.Date(2012, 10, 10, 8, 12, 03, 0, time.UTC),
	}
	tweet3 := tweet{
		User:     "sandrae",
		Retweets: 12,
		Message:  "Cycling is fun.",
		Created:  time.Date(2011, 11, 11, 10, 58, 12, 0, time.UTC),
	}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	all := NewMatchAllQuery()

	// Terms Aggregate by user name
	globalAgg := NewGlobalAggregation()
	usersAgg := NewTermsAggregation().Field("user").Size(10).OrderByCountDesc()
	avgRetweetsAgg := NewAvgAggregation().Field("retweets")
	minRetweetsAgg := NewMinAggregation().Field("retweets")
	maxRetweetsAgg := NewMaxAggregation().Field("retweets")
	sumRetweetsAgg := NewSumAggregation().Field("retweets")
	statsRetweetsAgg := NewStatsAggregation().Field("retweets")
	extstatsRetweetsAgg := NewExtendedStatsAggregation().Field("retweets")
	valueCountRetweetsAgg := NewValueCountAggregation().Field("retweets")
	percentilesRetweetsAgg := NewPercentilesAggregation().Field("retweets")
	cardinalityAgg := NewCardinalityAggregation().Field("user")
	significantTermsAgg := NewSignificantTermsAggregation().Field("message")
	retweetsRangeAgg := NewRangeAggregation().Field("retweets").
		Lt(10).Between(10, 100).Gt(100)
	dateRangeAgg := NewDateRangeAggregation().Field("created").
		Lt("2012-01-01").Between("2012-01-01", "2013-01-01").Gt("2013-01-01")
	missingImageAgg := NewMissingAggregation().Field("image")
	retweetsHistoAgg := NewHistogramAggregation().Field("retweets").Interval(100)
	dateHistoAgg := NewDateHistogramAggregation().Field("created").Interval("year")
	retweetsFilterAgg := NewFilterAggregation().Filter(
		NewRangeFilter("created").Gte("2012-01-01").Lte("2012-12-31"))

	// Run query
	searchResult, err := client.Search().Index(testIndexName).
		Query(&all).
		Aggregation("global", globalAgg).
		Aggregation("users", usersAgg).
		Aggregation("avgRetweets", avgRetweetsAgg).
		Aggregation("minRetweets", minRetweetsAgg).
		Aggregation("maxRetweets", maxRetweetsAgg).
		Aggregation("sumRetweets", sumRetweetsAgg).
		Aggregation("statsRetweets", statsRetweetsAgg).
		Aggregation("extstatsRetweets", extstatsRetweetsAgg).
		Aggregation("valueCountRetweets", valueCountRetweetsAgg).
		Aggregation("percentilesRetweets", percentilesRetweetsAgg).
		Aggregation("usersCardinality", cardinalityAgg).
		Aggregation("significantTerms", significantTermsAgg).
		Aggregation("retweetsRange", retweetsRangeAgg).
		Aggregation("dateRange", dateRangeAgg).
		Aggregation("missingImage", missingImageAgg).
		Aggregation("retweetsHisto", retweetsHistoAgg).
		Aggregation("dateHisto", dateHistoAgg).
		Aggregation("retweetsFilter", retweetsFilterAgg).
		Pretty(true).Debug(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if searchResult.Hits == nil {
		t.Errorf("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.Hits.TotalHits != 3 {
		t.Errorf("expected SearchResult.Hits.TotalHits = %d; got %d", 3, searchResult.Hits.TotalHits)
	}
	if len(searchResult.Hits.Hits) != 3 {
		t.Errorf("expected len(SearchResult.Hits.Hits) = %d; got %d", 3, len(searchResult.Hits.Hits))
	}
	if searchResult.Aggregations == nil {
		t.Errorf("expected SearchResult.Aggregations != nil; got nil")
	}

	// Search for non-existent aggregate should return (nil, false)
	agg, found := searchResult.Aggregations["no-such-aggregate"]
	if found {
		t.Errorf("expected SearchResult.Aggregations[...] = %v; got %v", false, found)
	}
	if agg != nil {
		t.Errorf("expected SearchResult.Aggregations[...] = nil; got %v", agg)
	}

	// Global
	agg, found = searchResult.Aggregations["global"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"global\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"global\"] != nil; got nil")
	}
	if agg.DocCount != 3 {
		t.Errorf("expected searchResult.Aggregations[\"global\"].DocCount = %v; got %v", 3, agg.DocCount)
	}

	// Search for existent aggregate (by name) should return (aggregate, true)
	agg, found = searchResult.Aggregations["users"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"users\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"users\"] != nil; got nil")
	}
	if len(agg.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"users\"].Buckets) = %v; got %v", 2, len(agg.Buckets))
	}
	if agg.Buckets[0].Key != "olivere" {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[0].Key = %v; got %v", "olivere", agg.Buckets[0].Key)
	}
	if agg.Buckets[0].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[0].DocCount = %v; got %v", 2, agg.Buckets[0].DocCount)
	}
	if agg.Buckets[1].Key != "sandrae" {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[1].Key = %v; got %v", "sandrae", agg.Buckets[1].Key)
	}
	if agg.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[1].DocCount = %v; got %v", 1, agg.Buckets[1].DocCount)
	}

	// avgRetweets
	agg, found = searchResult.Aggregations["avgRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"] != nil; got nil")
	}
	if agg.Value != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"].Value = %v; got %v", 40.0, agg.Value)
	}

	// minRetweets
	agg, found = searchResult.Aggregations["minRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"] != nil; got nil")
	}
	if agg.Value != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"].Value = %v; got %v", 0.0, agg.Value)
	}

	// maxRetweets
	agg, found = searchResult.Aggregations["maxRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"] != nil; got nil")
	}
	if agg.Value != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"].Value = %v; got %v", 108.0, agg.Value)
	}

	// sumRetweets
	agg, found = searchResult.Aggregations["sumRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"] != nil; got nil")
	}
	if agg.Value != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"].Value = %v; got %v", 120.0, agg.Value)
	}

	// statsRetweets
	agg, found = searchResult.Aggregations["statsRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"] != nil; got nil")
	}
	if agg.Count != 3 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Count = %v; got %v", 3, agg.Count)
	}
	if agg.Min != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Min = %v; got %v", 0.0, agg.Min)
	}
	if agg.Max != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Max = %v; got %v", 108.0, agg.Max)
	}
	if agg.Avg != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Avg = %v; got %v", 40.0, agg.Avg)
	}
	if agg.Sum != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Sum = %v; got %v", 120.0, agg.Sum)
	}

	// extstatsRetweets
	agg, found = searchResult.Aggregations["extstatsRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"] != nil; got nil")
	}
	if agg.Count != 3 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Count = %v; got %v", 3, agg.Count)
	}
	if agg.Min != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Min = %v; got %v", 0.0, agg.Min)
	}
	if agg.Max != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Max = %v; got %v", 108.0, agg.Max)
	}
	if agg.Avg != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Avg = %v; got %v", 40.0, agg.Avg)
	}
	if agg.Sum != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Sum = %v; got %v", 120.0, agg.Sum)
	}
	if agg.SumOfSquares != 11808.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].SumOfSquares = %v; got %v", 11808.0, agg.SumOfSquares)
	}
	if agg.Variance != 2336.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Variance = %v; got %v", 2336.0, agg.Variance)
	}
	if agg.StdDeviation != 48.33218389437829 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].StdDeviation = %v; got %v", 48.33218389437829, agg.StdDeviation)
	}

	// valueCountRetweets
	agg, found = searchResult.Aggregations["valueCountRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"] != nil; got nil")
	}
	if agg.Value != 3 {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"].Value = %v; got %v", 3, agg.Value)
	}

	// percentilesRetweets
	agg, found = searchResult.Aggregations["percentilesRetweets"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"] != nil; got nil")
	}
	// TODO(oe) How do we read the results?

	// usersCardinality
	agg, found = searchResult.Aggregations["usersCardinality"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"] != nil; got nil")
	}
	if agg.Value != 2 {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"].Value = %v; got %v", 2, agg.Value)
	}

	// significantTerms
	agg, found = searchResult.Aggregations["significantTerms"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"] != nil; got nil")
	}
	if agg.DocCount != 3 {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"].Value = %v; got %v", 3, agg.DocCount)
	}
	if len(agg.Buckets) != 0 {
		t.Errorf("expected len(searchResult.Aggregations[\"significantTerms\"].Buckets) = %v; got %v", 0, len(agg.Buckets))
	}

	// retweetsRange
	agg, found = searchResult.Aggregations["retweetsRange"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"] != nil; got nil")
	}
	if agg.DocCount != 0 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Value = %v; got %v", 3, agg.DocCount)
	}
	if len(agg.Buckets) != 3 {
		t.Errorf("expected len(searchResult.Aggregations[\"retweetsRange\"].Buckets) = %v; got %v", 0, len(agg.Buckets))
	}
	if agg.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[0].DocCount) = %v; got %v", 1, agg.Buckets[0].DocCount)
	}
	if agg.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[1].DocCount) = %v; got %v", 1, agg.Buckets[1].DocCount)
	}
	if agg.Buckets[2].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[2].DocCount) = %v; got %v", 1, agg.Buckets[2].DocCount)
	}

	// dateRange
	agg, found = searchResult.Aggregations["dateRange"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"] != nil; got nil")
	}
	if agg.DocCount != 0 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Value = %v; got %v", 0, agg.DocCount)
	}
	if len(agg.Buckets) != 3 {
		t.Errorf("expected len(searchResult.Aggregations[\"dateRange\"].Buckets) = %v; got %v", 3, len(agg.Buckets))
	}
	if agg.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].DocCount = %v; got %v", 1, agg.Buckets[0].DocCount)
	}
	if agg.Buckets[0].From != nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].From to be nil")
	}
	if agg.Buckets[0].To == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].To to be != nil")
	}
	if *agg.Buckets[0].To != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].To = %v; got %v", 1.325376e+12, *agg.Buckets[0].To)
	}
	if agg.Buckets[0].ToAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].ToAsStrSing to be != nil")
	}
	if *agg.Buckets[0].ToAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].ToAsString = %v; got %v", "2012-01-01T00:00:00.000Z", *agg.Buckets[0].ToAsString)
	}
	if agg.Buckets[1].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].DocCount = %v; got %v", 2, agg.Buckets[1].DocCount)
	}
	if agg.Buckets[1].From == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].From to be != nil")
	}
	if *agg.Buckets[1].From != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].From = %v; got %v", 1.325376e+12, *agg.Buckets[1].From)
	}
	if agg.Buckets[1].FromAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].FromAsString to be != nil")
	}
	if *agg.Buckets[1].FromAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].FromAsString = %v; got %v", "2012-01-01T00:00:00.000Z", *agg.Buckets[1].FromAsString)
	}
	if agg.Buckets[1].To == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].To to be != nil")
	}
	if *agg.Buckets[1].To != 1.3569984e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].To = %v; got %v", 1.3569984e+12, *agg.Buckets[1].To)
	}
	if agg.Buckets[1].ToAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].ToAsString to be != nil")
	}
	if *agg.Buckets[1].ToAsString != "2013-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].ToAsString = %v; got %v", "2013-01-01T00:00:00.000Z", *agg.Buckets[1].ToAsString)
	}
	if agg.Buckets[2].DocCount != 0 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].DocCount = %v; got %v", 0, agg.Buckets[2].DocCount)
	}
	if agg.Buckets[2].To != nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].To to be nil")
	}
	if agg.Buckets[2].From == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].From to be != nil")
	}
	if *agg.Buckets[2].From != 1.3569984e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].From = %v; got %v", 1.3569984e+12, *agg.Buckets[2].From)
	}
	if agg.Buckets[2].FromAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].FromAsString to be != nil")
	}
	if *agg.Buckets[2].FromAsString != "2013-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].FromAsString = %v; got %v", "2013-01-01T00:00:00.000Z", *agg.Buckets[2].FromAsString)
	}

	// missingImage
	agg, found = searchResult.Aggregations["missingImage"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"] != nil; got nil")
	}
	if agg.DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"].DocCount = %v; got %v", 2, agg.DocCount)
	}

	// retweetsHisto
	agg, found = searchResult.Aggregations["retweetsHisto"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"] != nil; got nil")
	}
	if len(agg.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"retweetsHisto\"].Buckets) = %v; got %v", 2, len(agg.Buckets))
	}
	if agg.Buckets[0].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[0].DocCount) = %v; got %v", 2, agg.Buckets[0].DocCount)
	}
	if agg.Buckets[0].Key != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[0].Key) = %v; got %v", 0.0, agg.Buckets[0].Key)
	}
	if agg.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[1].DocCount) = %v; got %v", 1, agg.Buckets[1].DocCount)
	}
	if agg.Buckets[1].Key != 100.0 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[1].Key) = %v; got %v", 100.0, agg.Buckets[1].Key)
	}

	// dateHisto
	agg, found = searchResult.Aggregations["dateHisto"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"] != nil; got nil")
	}
	if len(agg.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"dateHisto\"].Buckets) = %v; got %v", 2, len(agg.Buckets))
	}
	if agg.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].DocCount = %v; got %v", 1, agg.Buckets[0].DocCount)
	}
	if agg.Buckets[0].Key != 1.29384e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].Key = %v; got %v", 1.29384e+12, agg.Buckets[0].Key)
	}
	if agg.Buckets[0].KeyAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].KeyAsString != nil; got nil")
	}
	if *agg.Buckets[0].KeyAsString != "2011-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].KeyAsString = %v; got %v", "2011-01-01T00:00:00.000Z", agg.Buckets[0].KeyAsString)
	}
	if agg.Buckets[1].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].DocCount = %v; got %v", 2, agg.Buckets[1].DocCount)
	}
	if agg.Buckets[1].Key != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].Key = %v; got %v", 1.325376e+12, agg.Buckets[1].Key)
	}
	if agg.Buckets[1].KeyAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].KeyAsString != nil; got nil")
	}
	if *agg.Buckets[1].KeyAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].KeyAsString = %v; got %v", "2012-01-01T00:00:00.000Z", agg.Buckets[1].KeyAsString)
	}

	// retweetsFilter
	agg, found = searchResult.Aggregations["retweetsFilter"]
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsFilter\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsFilter\"] != nil; got nil")
	}
	if agg.DocCount != 2 {
		t.Fatalf("expected searchResult.Aggregations[\"retweetsFilter\"].DocCount = %v; got %v", 2, agg.DocCount)
	}
}
