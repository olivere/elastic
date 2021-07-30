// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
	"time"
)

func TestCatSnapshotsIntegration(t *testing.T) {
	if isCI() {
		t.Skip("this test requires local directories")
	}

	client := setupTestClientAndCreateIndexAndAddDocs(t, SetDecoder(&strictDecoder{})) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	{
		// Create a repository for this test
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, err := client.SnapshotCreateRepository("my_backup").
			Type("fs").
			Settings(map[string]interface{}{
				// Notice the path is configured as path.repo in docker-compose.yml
				"location": "/usr/share/elasticsearch/backup",
			}).
			Do(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// Make a snapshot
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, err = client.SnapshotCreate("my_backup", "snapshot_1").
			WaitForCompletion(true).
			Do(ctx)
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			// Remove snapshot
			_, _ = client.SnapshotDelete("my_backup", "snapshot_1").Do(context.Background())
			// Remove repository
			_, _ = client.SnapshotDeleteRepository("my_backup").Do(context.Background())
		}()
	}

	// List snapshots of repository
	ctx := context.Background()
	res, err := client.CatSnapshots().Repository("my_backup").Columns("*").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("want response, have nil")
	}
	if want, have := 1, len(res); want != have {
		t.Fatalf("want %d snapshot, have %d", want, have)
	}
	if want, have := "snapshot_1", res[0].ID; want != have {
		t.Fatalf("want ID=%q, have %q", want, have)
	}
}
