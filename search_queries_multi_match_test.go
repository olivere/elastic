package elastic

import (
	"encoding/json"
	"testing"
)

func TestMultiMatchQuery(t *testing.T) {
	q := NewMultiMatchQuery("this is a test", "subject", "message")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"multi_match":{"fields":["subject","message"],"query":"this is a test"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
