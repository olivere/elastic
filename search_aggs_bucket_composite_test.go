// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestCompositeAggregation(t *testing.T) {
	agg := NewCompositeAggregation().
		Sources(
			NewCompositeAggregationTermsValuesSource("my_terms").Field("a_term").Missing("N/A").Order("asc"),
			NewCompositeAggregationHistogramValuesSource("my_histogram", 5).Field("price").Asc(),
			NewCompositeAggregationDateHistogramValuesSource("my_date_histogram").CalendarInterval("1d").Field("purchase_date").Desc(),
		).
		Size(10).
		AggregateAfter(map[string]interface{}{
			"my_terms":          "1",
			"my_histogram":      2,
			"my_date_histogram": "3",
		})
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"composite":{"after":{"my_date_histogram":"3","my_histogram":2,"my_terms":"1"},"size":10,"sources":[{"my_terms":{"terms":{"field":"a_term","missing":"N/A","order":"asc"}}},{"my_histogram":{"histogram":{"field":"price","interval":5,"order":"asc"}}},{"my_date_histogram":{"date_histogram":{"calendar_interval":"1d","field":"purchase_date","order":"desc"}}}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCompositeAggregationTermsValuesSource(t *testing.T) {
	in := NewCompositeAggregationTermsValuesSource("products").
		Script(NewScript("doc['product'].value").Lang("painless"))
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"products":{"terms":{"script":{"lang":"painless","source":"doc['product'].value"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCompositeAggregationHistogramValuesSource(t *testing.T) {
	in := NewCompositeAggregationHistogramValuesSource("histo", 5).
		Field("price")
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"histo":{"histogram":{"field":"price","interval":5}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCompositeAggregationDateHistogramValuesSourceWithCalendarInterval(t *testing.T) {
	in := NewCompositeAggregationDateHistogramValuesSource("date").CalendarInterval("1d").
		Field("timestamp").
		Format("strict_date_optional_time")
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date":{"date_histogram":{"calendar_interval":"1d","field":"timestamp","format":"strict_date_optional_time"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestCompositeAggregationDateHistogramValuesSourceWithFixedInterval(t *testing.T) {
	in := NewCompositeAggregationDateHistogramValuesSource("date").FixedInterval("1d").
		Field("timestamp").
		Format("strict_date_optional_time")
	src, err := in.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"date":{"date_histogram":{"field":"timestamp","fixed_interval":"1d","format":"strict_date_optional_time"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
