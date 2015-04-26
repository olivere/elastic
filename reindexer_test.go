package elastic

import (
	"testing"
)

func TestReindexer(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	sourceCount, err := client.Count(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	targetCount, err := client.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != 0 {
		t.Fatalf("expected %d documents; got: %d", 0, targetCount)
	}

	r := NewReindexer(client, testIndexName, testIndexName2)
	ret, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatalf("expected result != %v; got: %v", nil, ret)
	}
	if ret.Success != sourceCount {
		t.Errorf("expected success = %d; got: %d", sourceCount, ret.Success)
	}
	if ret.Failed != 0 {
		t.Errorf("expected failed = %d; got: %d", 0, ret.Failed)
	}
	if len(ret.Errors) != 0 {
		t.Errorf("expected to return no errors by default; got: %v", ret.Errors)
	}

	if _, err := client.Flush().Index(testIndexName2).Do(); err != nil {
		t.Fatal(err)
	}

	targetCount, err = client.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != sourceCount {
		t.Fatalf("expected %d documents; got: %d", sourceCount, targetCount)
	}
}

func TestReindexerWithQuery(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	q := NewTermQuery("user", "olivere")

	sourceCount, err := client.Count(testIndexName).Query(q).Do()
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	targetCount, err := client.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != 0 {
		t.Fatalf("expected %d documents; got: %d", 0, targetCount)
	}

	r := NewReindexer(client, testIndexName, testIndexName2)
	r = r.Query(q)
	ret, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatalf("expected result != %v; got: %v", nil, ret)
	}
	if ret.Success != sourceCount {
		t.Errorf("expected success = %d; got: %d", sourceCount, ret.Success)
	}
	if ret.Failed != 0 {
		t.Errorf("expected failed = %d; got: %d", 0, ret.Failed)
	}
	if len(ret.Errors) != 0 {
		t.Errorf("expected to return no errors by default; got: %v", ret.Errors)
	}

	if _, err := client.Flush().Index(testIndexName2).Do(); err != nil {
		t.Fatal(err)
	}

	targetCount, err = client.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != sourceCount {
		t.Fatalf("expected %d documents; got: %d", sourceCount, targetCount)
	}
}

func TestReindexerProgress(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	sourceCount, err := client.Count(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	var calls int64
	totalsOk := true
	progress := func(current, total int64) {
		calls += 1
		totalsOk = totalsOk && total == sourceCount
	}

	r := NewReindexer(client, testIndexName, testIndexName2)
	r = r.Progress(progress)
	ret, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatalf("expected result != %v; got: %v", nil, ret)
	}
	if ret.Success != sourceCount {
		t.Errorf("expected success = %d; got: %d", sourceCount, ret.Success)
	}
	if ret.Failed != 0 {
		t.Errorf("expected failed = %d; got: %d", 0, ret.Failed)
	}
	if len(ret.Errors) != 0 {
		t.Errorf("expected to return no errors by default; got: %v", ret.Errors)
	}

	if calls != sourceCount {
		t.Errorf("expected progress to be called %d times; got: %d", sourceCount, calls)
	}
	if !totalsOk {
		t.Errorf("expected totals in progress to be %d", sourceCount)
	}
}

func TestReindexerWithTargetClient(t *testing.T) {
	sourceClient := setupTestClientAndCreateIndexAndAddDocs(t)
	targetClient, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	sourceCount, err := sourceClient.Count(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if sourceCount <= 0 {
		t.Fatalf("expected more than %d documents; got: %d", 0, sourceCount)
	}

	targetCount, err := targetClient.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != 0 {
		t.Fatalf("expected %d documents; got: %d", 0, targetCount)
	}

	r := NewReindexer(sourceClient, testIndexName, testIndexName2)
	r = r.TargetClient(targetClient)
	ret, err := r.Do()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatalf("expected result != %v; got: %v", nil, ret)
	}
	if ret.Success != sourceCount {
		t.Errorf("expected success = %d; got: %d", sourceCount, ret.Success)
	}
	if ret.Failed != 0 {
		t.Errorf("expected failed = %d; got: %d", 0, ret.Failed)
	}
	if len(ret.Errors) != 0 {
		t.Errorf("expected to return no errors by default; got: %v", ret.Errors)
	}

	if _, err := targetClient.Flush().Index(testIndexName2).Do(); err != nil {
		t.Fatal(err)
	}

	targetCount, err = targetClient.Count(testIndexName2).Do()
	if err != nil {
		t.Fatal(err)
	}
	if targetCount != sourceCount {
		t.Fatalf("expected %d documents; got: %d", sourceCount, targetCount)
	}
}
