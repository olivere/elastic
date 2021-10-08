// Copyright 2012-2018 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
		var status XPackWatchActionStatus
		err := json.Unmarshal(test.Input, &status)
		if err != nil {
			t.Errorf("#%d: expected no error, got %v", i+1, err)
		}
		if status.AckStatus == nil {
			t.Errorf("#%d: expected AckStatus!=nil", i+1)
		}
		if status.LastExecution == nil {
			t.Errorf("#%d: expected LastExecution!=nil", i+1)
		}
		if status.LastSuccessfulExecution == nil {
			t.Errorf("#%d: expected LastSuccessfulExecution!=nil", i+1)
		}
	}
}

func TestXPackWatchResponseParser(t *testing.T) {
	mustParseTime := func(layout, value string) time.Time {
		dt, err := time.Parse(layout, value)
		if err != nil {
			t.Fatal(err)
		}
		return dt
	}

	tests := []struct {
		Input     []byte
		ExpectErr bool
		Response  *XPackWatcherGetWatchResponse
	}{
		{
			[]byte(`
				{
					"found": true,
					"_id": "my_watch",
					"_seq_no": 0,
					"_primary_term": 1,
					"_version": 1,
					"status": { 
						"version": 1,
						"state": {
							"active": true,
							"timestamp": "2015-05-26T18:21:08.630Z"
						},
						"actions": {
							"test_index": {
								"ack": {
									"timestamp": "2015-05-26T18:21:08.630Z",
									"state": "awaits_successful_execution"
								}
							}
						}
					},
					"watch": {
						"input": {
							"simple": {
								"payload": {
									"send": "yes"
								}
							}
						},
						"condition": {
							"always": {}
						},
						"trigger": {
							"schedule": {
								"hourly": {
									"minute": [0, 5]
								}
							}
						},
						"actions": {
							"test_index": {
								"index": {
									"index": "test"
								}
							}
						}
					}
				}
			`),
			false,
			&XPackWatcherGetWatchResponse{
				Found:   true,
				Id:      "my_watch",
				Version: 1,
				Status: &XPackWatchStatus{
					State: &XPackWatchExecutionState{
						Active:    true,
						Timestamp: mustParseTime(time.RFC3339, "2015-05-26T18:21:08.630Z"),
					},
					Actions: map[string]*XPackWatchActionStatus{
						"test_index": {
							AckStatus: &XPackWatchActionAckStatus{
								State:     "awaits_successful_execution",
								Timestamp: mustParseTime(time.RFC3339, "2015-05-26T18:21:08.630Z"),
							},
						},
					},
					Version: 1,
				},
				Watch: &XPackWatch{
					Trigger: map[string]map[string]interface{}{
						"schedule": {
							"hourly": map[string]interface{}{"minute": []interface{}{float64(0), float64(5)}},
						},
					},
					Input:     map[string]map[string]interface{}{"simple": {"payload": map[string]interface{}{"send": string("yes")}}},
					Condition: map[string]map[string]interface{}{"always": {}},
					Actions:   map[string]map[string]interface{}{"test_index": {"index": map[string]interface{}{"index": string("test")}}},
				},
			},
		},
	}

	for i, test := range tests {
		var resp XPackWatcherGetWatchResponse
		err := json.Unmarshal(test.Input, &resp)
		if err != nil {
			if !test.ExpectErr {
				t.Errorf("#%d: expected no error, got %v", i+1, err)
			}
			continue
		}
		if test.ExpectErr {
			t.Errorf("#%d: expected error, got none", i+1)
		} else {
			if want, have := *test.Response, resp; !cmp.Equal(want, have) {
				t.Errorf("#%d: diff: %s\n", i+1, cmp.Diff(want, have))
			}
		}
	}
}
