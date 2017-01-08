package elastic

import (
	"fmt"
	"strings"
	"testing"
)

func TestIndicesAnalyze(t *testing.T) {
	var client = setupTestClient(t)
	analyzer := client.IndexAnalyze()
	res, err := analyzer.Text("hello hi guy").Do()

	if err != nil {
		t.Fatalf("analyzer response error", err)
	}

	if len(res.Tokens) != 3 {
		t.Errorf("index api analyze: expected %d, got %d", 3, len(res.Tokens))
	}
}

func TestIndicesAnalyzeDetail(t *testing.T) {
	var client = setupTestClient(t)
	analyzer := client.IndexAnalyze()
	res, err := analyzer.Text("hello hi guy").Explain(true).Do()

	if err != nil {
		t.Fatalf("receive unexpected error")
	}

	if len(res.Detail.Analyzer.Tokens) != 3 {
		t.Errorf("index api analyze: expected %d, got %d", 3, len(res.Detail.Tokenizer.Tokens))
	}
}

func TestIndicesAnalyzeWithIndex(t *testing.T) {
	var client = setupTestClient(t)
	analyzer := client.IndexAnalyze()
	_, err := analyzer.Index("foo").Text("hello hi guy").Do()

	if err == nil || !strings.Contains(err.Error(), "no such index") {
		t.Fatalf("Expect error: no such index")
	}
}

func TestIndicesAnalyzeValidate(t *testing.T) {
	var client = setupTestClient(t)
	analyzer := client.IndexAnalyze()
	_, err := analyzer.Do()

	fmt.Println(err)
	if err == nil || !strings.Contains(err.Error(), "missing required fields") {
		t.Fatalf("Expect error: missing required fields")
	}
}
