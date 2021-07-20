package elastic

import (
	"context"
	"testing"
)

func TestCatSnapshots(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))
	ctx := context.Background()
	res, err := client.CatSnapshots().Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].ID; have == "" {
		t.Fatalf("ID[0]: want != %q, have %q", "", have)
	}
}
