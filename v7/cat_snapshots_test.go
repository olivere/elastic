// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestCatSnapshots(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	urls, params, err := client.CatSnapshots().Repository("my_repo").Columns("*").buildURL()
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "/_cat/snapshots/my_repo", urls; want != have {
		t.Fatalf("want URL=%q, have %q", want, have)
	}
	if want, have := "format=json&h=%2A", params.Encode(); want != have {
		t.Fatalf("want Params=%q, have %q", want, have)
	}
}
