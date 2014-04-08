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
	retweetsFilterAgg := NewFilterAggregation().Filter(
		NewRangeFilter("created").Gte("2012-01-01").Lte("2012-12-31")).
		SubAggregation("avgRetweets", NewAvgAggregation().Field("retweets"))
	significantTermsAgg := NewSignificantTermsAggregation().Field("message")
	retweetsRangeAgg := NewRangeAggregation().Field("retweets").
		Lt(10).Between(10, 100).Gt(100)
	dateRangeAgg := NewDateRangeAggregation().Field("created").
		Lt("2012-01-01").Between("2012-01-01", "2013-01-01").Gt("2013-01-01")
	missingImageAgg := NewMissingAggregation().Field("image")
	retweetsHistoAgg := NewHistogramAggregation().Field("retweets").Interval(100)
	dateHistoAgg := NewDateHistogramAggregation().Field("created").Interval("year")

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
	agg, found := searchResult.GetAggregation("no-such-aggregate")
	if found {
		t.Errorf("expected SearchResult.Aggregations[...] = %v; got %v", false, found)
	}
	if agg != nil {
		t.Errorf("expected SearchResult.Aggregations[...] = nil; got %v", agg)
	}

	// Global
	agg, found = searchResult.GetAggregation("global")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"global\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"global\"] != nil; got nil")
	}
	globalAggRes, found := agg.Global()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"global\"] = %v; got %v", true, found)
	}
	if globalAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"global\"] != nil; got nil")
	}
	if globalAggRes.DocCount != 3 {
		t.Errorf("expected searchResult.Aggregations[\"global\"].DocCount = %v; got %v", 3, globalAggRes.DocCount)
	}

	// Search for existent aggregate (by name) should return (aggregate, true)
	agg, found = searchResult.GetAggregation("users")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"users\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"users\"] != nil; got nil")
	}
	termsAggRes, found := agg.Terms()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"users\"] = %v; got %v", true, found)
	}
	if termsAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"users\"] != nil; got nil")
	}
	if len(termsAggRes.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"users\"].Buckets) = %v; got %v", 2, len(termsAggRes.Buckets))
	}
	if termsAggRes.Buckets[0].Key != "olivere" {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[0].Key = %v; got %v", "olivere", termsAggRes.Buckets[0].Key)
	}
	if termsAggRes.Buckets[0].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[0].DocCount = %v; got %v", 2, termsAggRes.Buckets[0].DocCount)
	}
	if termsAggRes.Buckets[1].Key != "sandrae" {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[1].Key = %v; got %v", "sandrae", termsAggRes.Buckets[1].Key)
	}
	if termsAggRes.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"users\"].Buckets[1].DocCount = %v; got %v", 1, termsAggRes.Buckets[1].DocCount)
	}

	// avgRetweets
	agg, found = searchResult.GetAggregation("avgRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"avgRetweets#\"] != nil; got nil")
	}
	avgAggRes, found := agg.Avg()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"] = %v; got %v", true, found)
	}
	if avgAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"] != nil; got nil")
	}
	if avgAggRes.Value != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"avgRetweets\"].Value = %v; got %v", 40.0, avgAggRes.Value)
	}

	// minRetweets
	agg, found = searchResult.GetAggregation("minRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"minRetweets\"] != nil; got nil")
	}
	minAggRes, found := agg.Min()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"] = %v; got %v", true, found)
	}
	if minAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"] != nil; got nil")
	}
	if minAggRes.Value != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"minRetweets\"].Value = %v; got %v", 0.0, minAggRes.Value)
	}

	// maxRetweets
	agg, found = searchResult.GetAggregation("maxRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"maxRetweets\"] != nil; got nil")
	}
	maxAggRes, found := agg.Max()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"] = %v; got %v", true, found)
	}
	if maxAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"] != nil; got nil")
	}
	if maxAggRes.Value != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"maxRetweets\"].Value = %v; got %v", 108.0, maxAggRes.Value)
	}

	// sumRetweets
	agg, found = searchResult.GetAggregation("sumRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"sumRetweets\"] != nil; got nil")
	}
	sumAggRes, found := agg.Sum()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"] = %v; got %v", true, found)
	}
	if sumAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"] != nil; got nil")
	}
	if sumAggRes.Value != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"sumRetweets\"].Value = %v; got %v", 120.0, sumAggRes.Value)
	}

	// statsRetweets
	agg, found = searchResult.GetAggregation("statsRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"statsRetweets\"] != nil; got nil")
	}
	statsAggRes, found := agg.Stats()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"] = %v; got %v", true, found)
	}
	if statsAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"] != nil; got nil")
	}
	if statsAggRes.Count != 3 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Count = %v; got %v", 3, statsAggRes.Count)
	}
	if statsAggRes.Min != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Min = %v; got %v", 0.0, statsAggRes.Min)
	}
	if statsAggRes.Max != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Max = %v; got %v", 108.0, statsAggRes.Max)
	}
	if statsAggRes.Avg != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Avg = %v; got %v", 40.0, statsAggRes.Avg)
	}
	if statsAggRes.Sum != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"statsRetweets\"].Sum = %v; got %v", 120.0, statsAggRes.Sum)
	}

	// extstatsRetweets
	agg, found = searchResult.GetAggregation("extstatsRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"extstatsRetweets\"] != nil; got nil")
	}
	extStatsAggRes, found := agg.ExtendedStats()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"] = %v; got %v", true, found)
	}
	if extStatsAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"] != nil; got nil")
	}
	if extStatsAggRes.Count != 3 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Count = %v; got %v", 3, extStatsAggRes.Count)
	}
	if extStatsAggRes.Min != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Min = %v; got %v", 0.0, extStatsAggRes.Min)
	}
	if extStatsAggRes.Max != 108.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Max = %v; got %v", 108.0, extStatsAggRes.Max)
	}
	if extStatsAggRes.Avg != 40.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Avg = %v; got %v", 40.0, extStatsAggRes.Avg)
	}
	if extStatsAggRes.Sum != 120.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Sum = %v; got %v", 120.0, extStatsAggRes.Sum)
	}
	if extStatsAggRes.SumOfSquares != 11808.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].SumOfSquares = %v; got %v", 11808.0, extStatsAggRes.SumOfSquares)
	}
	if extStatsAggRes.Variance != 2336.0 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].Variance = %v; got %v", 2336.0, extStatsAggRes.Variance)
	}
	if extStatsAggRes.StdDeviation != 48.33218389437829 {
		t.Errorf("expected searchResult.Aggregations[\"extstatsRetweets\"].StdDeviation = %v; got %v", 48.33218389437829, extStatsAggRes.StdDeviation)
	}

	// valueCountRetweets
	agg, found = searchResult.GetAggregation("valueCountRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"valueCountRetweets\"] != nil; got nil")
	}
	valueCountAggRes, found := agg.ValueCount()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"] = %v; got %v", true, found)
	}
	if valueCountAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"] != nil; got nil")
	}
	if valueCountAggRes.Value != 3 {
		t.Errorf("expected searchResult.Aggregations[\"valueCountRetweets\"].Value = %v; got %v", 3, valueCountAggRes.Value)
	}

	// percentilesRetweets
	agg, found = searchResult.GetAggregation("percentilesRetweets")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"percentilesRetweets\"] != nil; got nil")
	}
	percentilesAggRes, found := agg.Percentiles()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"] = %v; got %v", true, found)
	}
	if percentilesAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"] != nil; got nil")
	}
	if len(percentilesAggRes.Values) != 7 {
		t.Fatalf("expected len(searchResult.Aggregations[\"percentilesRetweets\"].Values) == 7; got %v", len(percentilesAggRes.Values))
	}
	if percentilesAggRes.Values["0.0"] != nil {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"].Values[\"0.0\"] == nil; got %v", percentilesAggRes.Values["0.0"])
	}
	if percentilesAggRes.Values["1.0"] != 0.24 {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"].Values[\"1.0\"] == %v; got %v", 0.24, percentilesAggRes.Values["1.0"])
	}
	if percentilesAggRes.Values["25.0"] != 6.0 {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"].Values[\"1.0\"] == %v; got %v", 6.0, percentilesAggRes.Values["25.0"])
	}
	if percentilesAggRes.Values["99.0"] != 106.08 {
		t.Errorf("expected searchResult.Aggregations[\"percentilesRetweets\"].Values[\"1.0\"] == %v; got %v", 106.08, percentilesAggRes.Values["99.0"])
	}

	// usersCardinality
	agg, found = searchResult.GetAggregation("usersCardinality")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"usersCardinality\"] != nil; got nil")
	}
	cardAggRes, found := agg.Cardinality()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"] = %v; got %v", true, found)
	}
	if cardAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"] != nil; got nil")
	}
	if cardAggRes.Value != 2 {
		t.Errorf("expected searchResult.Aggregations[\"usersCardinality\"].Value = %v; got %v", 2, cardAggRes.Value)
	}

	// retweetsFilter
	agg, found = searchResult.GetAggregation("retweetsFilter")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsFilter\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"retweetsFilter\"] != nil; got nil")
	}
	filterAggRes, found := agg.Filter()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsFilter\"] = %v; got %v", true, found)
	}
	if filterAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsFilter\"] != nil; got nil")
	}
	if filterAggRes.DocCount != 2 {
		t.Fatalf("expected searchResult.Aggregations[\"retweetsFilter\"].DocCount = %v; got %v", 2, filterAggRes.DocCount)
	}
	// Retrieve sub-aggregation
	subAgg, found := agg.GetAggregation("avgRetweets")
	if !found {
		t.Error("expected sub-aggregation \"avgRetweets\" to be found; got false")
	}
	if subAgg == nil {
		t.Fatal("expected sub-aggregation \"avgRetweets\"; got nil")
	}
	avgRetweetsAggRes, found := subAgg.Avg()
	if !found {
		t.Error("expected sub-aggregation \"avgRetweets\" to be found; got false")
	}
	if avgRetweetsAggRes == nil {
		t.Fatal("expected sub-aggregation \"avgRetweets\"; got nil")
	}
	if avgRetweetsAggRes.Value != 54.0 {
		t.Errorf("expected sub-aggregation \"avgRetweets\" to have value = %v; got %v", 54.0, avgRetweetsAggRes.Value)
	}

	// significantTerms
	agg, found = searchResult.GetAggregation("significantTerms")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"significantTerms\"] != nil; got nil")
	}
	stAggRes, found := agg.SignificantTerms()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"] = %v; got %v", true, found)
	}
	if stAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"] != nil; got nil")
	}
	if stAggRes.DocCount != 3 {
		t.Errorf("expected searchResult.Aggregations[\"significantTerms\"].Value = %v; got %v", 3, stAggRes.DocCount)
	}
	if len(stAggRes.Buckets) != 0 {
		t.Errorf("expected len(searchResult.Aggregations[\"significantTerms\"].Buckets) = %v; got %v", 0, len(stAggRes.Buckets))
	}

	// retweetsRange
	agg, found = searchResult.GetAggregation("retweetsRange")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"retweetsRange\"] != nil; got nil")
	}
	rangeAggRes, found := agg.Range()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"] = %v; got %v", true, found)
	}
	if rangeAggRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"] != nil; got nil")
	}
	if len(rangeAggRes.Buckets) != 3 {
		t.Errorf("expected len(searchResult.Aggregations[\"retweetsRange\"].Buckets) = %v; got %v", 3, len(rangeAggRes.Buckets))
	}
	if rangeAggRes.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[0].DocCount) = %v; got %v", 1, rangeAggRes.Buckets[0].DocCount)
	}
	if rangeAggRes.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[1].DocCount) = %v; got %v", 1, rangeAggRes.Buckets[1].DocCount)
	}
	if rangeAggRes.Buckets[2].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsRange\"].Buckets[2].DocCount) = %v; got %v", 1, rangeAggRes.Buckets[2].DocCount)
	}

	// dateRange
	agg, found = searchResult.GetAggregation("dateRange")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"] != nil; got nil")
	}
	dateRangeRes, found := agg.DateRange()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"] = %v; got %v", true, found)
	}
	if dateRangeRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"] != nil; got nil")
	}
	if dateRangeRes.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].DocCount = %v; got %v", 1, dateRangeRes.Buckets[0].DocCount)
	}
	if dateRangeRes.Buckets[0].From != nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].From to be nil")
	}
	if dateRangeRes.Buckets[0].To == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].To to be != nil")
	}
	if *dateRangeRes.Buckets[0].To != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].To = %v; got %v", 1.325376e+12, *dateRangeRes.Buckets[0].To)
	}
	if dateRangeRes.Buckets[0].ToAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].ToAsStrSing to be != nil")
	}
	if *dateRangeRes.Buckets[0].ToAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[0].ToAsString = %v; got %v", "2012-01-01T00:00:00.000Z", *dateRangeRes.Buckets[0].ToAsString)
	}
	if dateRangeRes.Buckets[1].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].DocCount = %v; got %v", 2, dateRangeRes.Buckets[1].DocCount)
	}
	if dateRangeRes.Buckets[1].From == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].From to be != nil")
	}
	if *dateRangeRes.Buckets[1].From != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].From = %v; got %v", 1.325376e+12, *dateRangeRes.Buckets[1].From)
	}
	if dateRangeRes.Buckets[1].FromAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].FromAsString to be != nil")
	}
	if *dateRangeRes.Buckets[1].FromAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].FromAsString = %v; got %v", "2012-01-01T00:00:00.000Z", *dateRangeRes.Buckets[1].FromAsString)
	}
	if dateRangeRes.Buckets[1].To == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].To to be != nil")
	}
	if *dateRangeRes.Buckets[1].To != 1.3569984e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].To = %v; got %v", 1.3569984e+12, *dateRangeRes.Buckets[1].To)
	}
	if dateRangeRes.Buckets[1].ToAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].ToAsString to be != nil")
	}
	if *dateRangeRes.Buckets[1].ToAsString != "2013-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[1].ToAsString = %v; got %v", "2013-01-01T00:00:00.000Z", *dateRangeRes.Buckets[1].ToAsString)
	}
	if dateRangeRes.Buckets[2].DocCount != 0 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].DocCount = %v; got %v", 0, dateRangeRes.Buckets[2].DocCount)
	}
	if dateRangeRes.Buckets[2].To != nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].To to be nil")
	}
	if dateRangeRes.Buckets[2].From == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].From to be != nil")
	}
	if *dateRangeRes.Buckets[2].From != 1.3569984e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].From = %v; got %v", 1.3569984e+12, *dateRangeRes.Buckets[2].From)
	}
	if dateRangeRes.Buckets[2].FromAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].FromAsString to be != nil")
	}
	if *dateRangeRes.Buckets[2].FromAsString != "2013-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateRange\"].Buckets[2].FromAsString = %v; got %v", "2013-01-01T00:00:00.000Z", *dateRangeRes.Buckets[2].FromAsString)
	}

	// missingImage
	agg, found = searchResult.GetAggregation("missingImage")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"missingImage\"] != nil; got nil")
	}
	missingRes, found := agg.Missing()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"] = %v; got %v", true, found)
	}
	if missingRes == nil {
		t.Fatalf("expected searchResult.Aggregations[\"missingImage\"] != nil; got nil")
	}
	if missingRes.DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"missingImage\"].DocCount = %v; got %v", 2, missingRes.DocCount)
	}

	// retweetsHisto
	agg, found = searchResult.GetAggregation("retweetsHisto")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"retweetsHisto\"] != nil; got nil")
	}
	histoRes, found := agg.Histogram()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"] = %v; got %v", true, found)
	}
	if histoRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"] != nil; got nil")
	}
	if len(histoRes.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"retweetsHisto\"].Buckets) = %v; got %v", 2, len(histoRes.Buckets))
	}
	if histoRes.Buckets[0].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[0].DocCount) = %v; got %v", 2, histoRes.Buckets[0].DocCount)
	}
	if histoRes.Buckets[0].Key != 0.0 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[0].Key) = %v; got %v", 0.0, histoRes.Buckets[0].Key)
	}
	if histoRes.Buckets[1].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[1].DocCount) = %v; got %v", 1, histoRes.Buckets[1].DocCount)
	}
	if histoRes.Buckets[1].Key != 100.0 {
		t.Errorf("expected searchResult.Aggregations[\"retweetsHisto\"].Buckets[1].Key) = %v; got %v", 100.0, histoRes.Buckets[1].Key)
	}

	// dateHisto
	agg, found = searchResult.GetAggregation("dateHisto")
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"] = %v; got %v", true, found)
	}
	if agg == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateHisto\"] != nil; got nil")
	}
	dateHistoRes, found := agg.DateHistogram()
	if !found {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"] = %v; got %v", true, found)
	}
	if dateHistoRes == nil {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"] != nil; got nil")
	}
	if len(dateHistoRes.Buckets) != 2 {
		t.Errorf("expected len(searchResult.Aggregations[\"dateHisto\"].Buckets) = %v; got %v", 2, len(dateHistoRes.Buckets))
	}
	if dateHistoRes.Buckets[0].DocCount != 1 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].DocCount = %v; got %v", 1, dateHistoRes.Buckets[0].DocCount)
	}
	if dateHistoRes.Buckets[0].Key != 1.29384e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].Key = %v; got %v", 1.29384e+12, dateHistoRes.Buckets[0].Key)
	}
	if dateHistoRes.Buckets[0].KeyAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].KeyAsString != nil; got nil")
	}
	if *dateHistoRes.Buckets[0].KeyAsString != "2011-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[0].KeyAsString = %v; got %v", "2011-01-01T00:00:00.000Z", dateHistoRes.Buckets[0].KeyAsString)
	}
	if dateHistoRes.Buckets[1].DocCount != 2 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].DocCount = %v; got %v", 2, dateHistoRes.Buckets[1].DocCount)
	}
	if dateHistoRes.Buckets[1].Key != 1.325376e+12 {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].Key = %v; got %v", 1.325376e+12, dateHistoRes.Buckets[1].Key)
	}
	if dateHistoRes.Buckets[1].KeyAsString == nil {
		t.Fatalf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].KeyAsString != nil; got nil")
	}
	if *dateHistoRes.Buckets[1].KeyAsString != "2012-01-01T00:00:00.000Z" {
		t.Errorf("expected searchResult.Aggregations[\"dateHisto\"].Buckets[1].KeyAsString = %v; got %v", "2012-01-01T00:00:00.000Z", dateHistoRes.Buckets[1].KeyAsString)
	}
}
