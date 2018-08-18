// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestWatcherActivateWatch(t *testing.T) {
	client := setupTestClient(t)
	watcher := NewXpackWatcherActivateWatchService(client)
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
