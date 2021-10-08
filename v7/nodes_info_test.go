// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestNodesInfoBuildURL(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tests := []struct {
		NodeIDs  []string
		Metrics  []string
		Expected string
	}{
		{
			nil,
			nil,
			"/_nodes/_all/_all",
		},
		{
			[]string{},
			[]string{},
			"/_nodes/_all/_all",
		},
		{
			[]string{"node1"},
			[]string{},
			"/_nodes/node1/_all",
		},
		{
			[]string{"node1", "node2"},
			nil,
			"/_nodes/node1%2Cnode2/_all",
		},
		{
			[]string{"node1", "node2"},
			[]string{"metric1", "metric2"},
			"/_nodes/node1%2Cnode2/metric1%2Cmetric2",
		},
	}

	for i, test := range tests {
		path, _, err := client.NodesInfo().NodeId(test.NodeIDs...).Metric(test.Metrics...).buildURL()
		if err != nil {
			t.Fatalf("case #%d: %v", i+1, err)
		}
		if path != test.Expected {
			t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
		}
	}
}
func TestNodesInfo(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	info, err := client.NodesInfo().Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if info == nil {
		t.Fatal("expected nodes info")
	}

	if info.ClusterName == "" {
		t.Errorf("expected cluster name; got: %q", info.ClusterName)
	}
	if len(info.Nodes) == 0 {
		t.Errorf("expected some nodes; got: %d", len(info.Nodes))
	}
	for id, node := range info.Nodes {
		if id == "" {
			t.Errorf("expected node id; got: %q", id)
		}
		if node == nil {
			t.Fatalf("expected node info; got: %v", node)
		}
		if node.IP == "" {
			t.Errorf("expected node IP; got: %q", node.IP)
		}
	}
}
