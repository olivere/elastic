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
		t.Fatalf("expected res to be != nil; got: %v", res)
	}
	if res.Status != "green" && res.Status != "red" && res.Status != "yellow" {
		t.Fatalf("expected status \"green\", \"red\", or \"yellow\"; got: %q", res.Status)
	}
}

func TestClusterHealthURLs(t *testing.T) {
	tests := []struct {
		Service  *ClusterHealthService
		Expected string
	}{
		{
			Service: &ClusterHealthService{
				indices: []string{},
			},
			Expected: "/_cluster/health/",
		},
		{
			Service: &ClusterHealthService{
				indices: []string{"twitter"},
			},
			Expected: "/_cluster/health/twitter",
		},
		{
			Service: &ClusterHealthService{
				indices: []string{"twitter", "gplus"},
			},
			Expected: "/_cluster/health/twitter%2Cgplus",
		},
		{
			Service: &ClusterHealthService{
				indices:       []string{"twitter"},
				waitForStatus: "yellow",
			},
			Expected: "/_cluster/health/twitter?wait_for_status=yellow",
		},
	}

	for _, test := range tests {
		got, err := test.Service.buildURL()
		if err != nil {
			t.Fatalf("expected no error; got: %v", err)
		}
		if got != test.Expected {
			t.Errorf("expected URL = %q; got: %q", test.Expected, got)
		}
	}
}
