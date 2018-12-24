package elastic

import (
	"encoding/json"
	"testing"
)


//"span_multi":{
//	"match":{
//		"prefix" : { "user" :  { "value" : "ki" } }
//	}
//}expected := `{"span_multi":{"match":{"prefix": {"user: {"value": "ki""}}}`
func TestSpanMultiQuery(t *testing.T) {
	q := NewSpanMultiQuery()
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_multi":{"match":{}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestSpanMultiQueryWithQuery(t *testing.T) {
	q1 := NewPrefixQuery("user", "ki")
	q := NewSpanMultiQuery()
	q.Match(q1)

	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"span_multi":{"match":{"prefix":{"user":"ki"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
