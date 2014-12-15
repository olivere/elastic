// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestClusterHealth(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	// Get cluster health
	res, err := client.ClusterHealth().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected res to be != nil; got: %v", res)
	}
	if res.Status != "green" {
		t.Fatalf("expected status %q; got: %q", "green", res.Status)
	}
}
