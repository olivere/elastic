// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestClusterStats(t *testing.T) {
	client := setupTestClientAndCreateIndexAndAddDocs(t)

	// Get cluster stats
	res, err := client.ClusterStats().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatalf("expected res to be != nil; got: %v", res)
	}
	if res.ClusterName == "" {
		t.Fatalf("expected a cluster name; got: %q", res.ClusterName)
	}
	if res.Nodes == nil {
		t.Fatalf("expected nodes; got: %v", res.Nodes)
	}
	if res.Nodes.Count == nil {
		t.Fatalf("expected nodes count; got: %v", res.Nodes.Count)
	}
}
