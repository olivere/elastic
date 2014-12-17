// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

func TestClusterState(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	// Get cluster state
	res, err := client.ClusterState().Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatalf("expected res to be != nil; got: %v", res)
	}
	if res.ClusterName == "" {
		t.Fatalf("expected a cluster name; got: %q", res.ClusterName)
	}
}

func TestClusterStateURLs(t *testing.T) {
	tests := []struct {
		Service  *ClusterStateService
		Expected string
	}{
		{
			Service: &ClusterStateService{
				indices: []string{},
				metrics: []string{},
			},
			Expected: "/_cluster/state/_all/_all",
		},
		{
			Service: &ClusterStateService{
				indices: []string{"twitter"},
				metrics: []string{},
			},
			Expected: "/_cluster/state/_all/twitter",
		},
		{
			Service: &ClusterStateService{
				indices: []string{"twitter", "gplus"},
				metrics: []string{},
			},
			Expected: "/_cluster/state/_all/twitter%2Cgplus",
		},
		{
			Service: &ClusterStateService{
				indices: []string{},
				metrics: []string{"nodes"},
			},
			Expected: "/_cluster/state/nodes/_all",
		},
		{
			Service: &ClusterStateService{
				indices: []string{"twitter"},
				metrics: []string{"nodes"},
			},
			Expected: "/_cluster/state/nodes/twitter",
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
