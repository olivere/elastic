// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

type WatcherBody struct {
	Trigger   map[string]interface{} `json:"trigger"`
	Input     map[string]interface{} `json:"input"`
	Condition map[string]interface{} `json:"condition"`
	Actions   map[string]interface{} `json:"actions"`
}

func TestWatcherPutWatch(t *testing.T) {
	client := setupTestClient(t)
	watcher := NewXpackWatcherPutWatchService(client)
	data, err := json.Marshal(watcher)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
