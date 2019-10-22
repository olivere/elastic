// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestXPackWatcherGetWatchBuildURL(t *testing.T) {
	client := setupTestClient(t) // , SetURL("http://elastic:elastic@localhost:9210"))

	tests := []struct {
		Id        string
		Expected  string
		ExpectErr bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"my-watch",
			"/_watcher/watch/my-watch",
			false,
		},
	}

	for i, test := range tests {
		builder := client.XPackWatchGet(test.Id)
		err := builder.Validate()
		if err != nil {
			if !test.ExpectErr {
				t.Errorf("case #%d: %v", i+1, err)
				continue
			}
		} else {
			// err == nil
			if test.ExpectErr {
				t.Errorf("case #%d: expected error", i+1)
				continue
			}
			path, _, _ := builder.buildURL()
			if path != test.Expected {
				t.Errorf("case #%d: expected %q; got: %q", i+1, test.Expected, path)
			}
		}
	}
}

func TestXPackWatchActionStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Input     []byte
		ExpectErr bool
	}{
		{
			[]byte(`
			   {
			     "ack" : {
			       "timestamp" : "2019-10-22T15:01:12.163Z",
			       "state" : "ackable"
			     },
			     "last_execution" : {
			       "timestamp" : "2019-10-22T15:01:12.163Z",
			       "successful" : true
			     },
			     "last_successful_execution" : {
			       "timestamp" : "2019-10-22T15:01:12.163Z",
			       "successful" : true
			     }
			   }
			`),
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var status XPackWatchActionStatus
			err := json.Unmarshal(test.Input, &status)
			if err != nil {
				t.Error(err)
			}
			if status.AckStatus == nil {
				t.Error("nil AckStatus")
			}
			if status.LastExecution == nil {
				t.Error("nil LastExecution")
			}
			if status.LastSuccessfulExecution == nil {
				t.Error("nil LastSuccessfulExecution")
			}
		})
	}
}
