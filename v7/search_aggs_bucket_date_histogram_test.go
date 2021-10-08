// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestDateHistogramAggregationLegacyInterval(t *testing.T) {
	agg := NewDateHistogramAggregation().
		Field("date").
		Interval("week").
		Format("yyyy-MM").
		TimeZone("UTC").
		Offset("+6h")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"field":"date","format":"yyyy-MM","interval":"week","offset":"+6h","time_zone":"UTC"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDateHistogramAggregationFixed(t *testing.T) {
	agg := NewDateHistogramAggregation().
		Field("date").
		FixedInterval("month").
		Format("yyyy-MM").
		TimeZone("UTC").
		Offset("+6h")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"field":"date","fixed_interval":"month","format":"yyyy-MM","offset":"+6h","time_zone":"UTC"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDateHistogramAggregationCalendar(t *testing.T) {
	agg := NewDateHistogramAggregation().
		Field("date").
		CalendarInterval("1d").
		Format("yyyy-MM").
		TimeZone("UTC").
		Offset("+6h")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"calendar_interval":"1d","field":"date","format":"yyyy-MM","offset":"+6h","time_zone":"UTC"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDateHistogramAggregationWithKeyedResponse(t *testing.T) {
	agg := NewDateHistogramAggregation().Field("date").CalendarInterval("year").Missing("1900").Keyed(true)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"calendar_interval":"year","field":"date","keyed":true,"missing":"1900"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestDateHistogramAggregationWithMissing(t *testing.T) {
	agg := NewDateHistogramAggregation().Field("date").CalendarInterval("year").Missing("1900")
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date_histogram":{"calendar_interval":"year","field":"date","missing":"1900"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
