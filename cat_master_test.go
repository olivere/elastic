package elastic

import (
	"context"
	"testing"
)

func TestCatMaster(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{}))
	ctx := context.Background()
	res, err := client.CatMaster().Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if len(res) == 0 {
		t.Fatalf("want response, have: %v", res)
	}
	if have := res[0].IP; have == "" {
		t.Fatalf("IP[0]: want != %q, have %q", "", have)
	}
}
