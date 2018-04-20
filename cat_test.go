// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestCatModels(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	c := client.Cat()

	tests := []struct {
		f    func() *CatService
		name string
	}{
		{f: c.Aliases, name: "aliases"},
		{f: c.Allocation, name: "allocation"},
		{f: c.Count, name: "count"},
		{f: c.FieldData, name: "fielddata"},
		{f: c.Health, name: "health"},
		{f: c.Indices, name: "indices"},
		{f: c.Master, name: "master"},
		{f: c.NodeAttrs, name: "nodeattrs"},
		{f: c.Nodes, name: "nodes"},
		{f: c.PendingTasks, name: "pendingtasks"},
		{f: c.Plugins, name: "plugins"},
		{f: c.Recovery, name: "recovery"},
		{f: c.Repositories, name: "repositories"},
		{f: c.Shards, name: "shards"},
		{f: c.Segments, name: "segments"},
		// TODO - it is not possible to cat snapshots where no repo exists.
		// Creating one for the purpose of this test fails without setting the path.repo or repositories.url.allowed_url setting .
		// Putting this off to a later issue.
		{f: c.Templates, name: "templates"},
		{f: c.ThreadPool, name: "threadpool"},
	}

	for _, tt := range tests {
		model := tt.f
		res, err := model().Do(context.TODO())

		if err != nil {
			t.Errorf("model %s - expected to not get an error, got %v", tt.name, err)
		}

		if res == "" {
			t.Errorf("model %s - expected res to not be an empty string", tt.name)
		}
	}
}

func TestCatUnknownModel(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	c := client.Cat()
	c.model = "boogers"

	_, err := c.Do(context.TODO())

	if err == nil {
		t.Fatalf("expected to get an error, but did not")
	}
}

func TestCatNoModel(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	_, err := client.Cat().Do(context.TODO())

	if err == nil {
		t.Fatalf("expected to get an error, but did not")
	}
}

func TestCatHelpCommand(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	res, err := client.Cat().Health().Help(true).Do(context.TODO())

	if err != nil {
		t.Errorf("expected to not get an error with help command, got %v", err)
	}

	if res == "" {
		t.Errorf("expected response from help command to not be an empty string")
	}
}
