package elastic

import (
	"encoding/json"
	"testing"
)

func TestWildcardQuery(t *testing.T) {
	q := NewWildcardQuery("user", "*ki*")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":"*ki*"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestWildcardQueryWithOptions(t *testing.T) {
	q := NewWildcardQuery("user", "*ki*")
	q = q.QueryName("my_query_name")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"_name":"my_query_name","wildcard":"*ki*"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
