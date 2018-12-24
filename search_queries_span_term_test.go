package elastic

import (
	"encoding/json"
	"testing"
)

func TestSpanTermQuery(t *testing.T) {
	q := NewSpanTermQuery("user", "ki")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_term":{"user":"ki"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSpanTermQueryWithOptions(t *testing.T) {
	q := NewSpanTermQuery("user", "ki")
	q = q.Boost(2.79)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_term":{"boost":2.79,"user":"ki"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
