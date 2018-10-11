package elastic

import (
	"context"
	"testing"
)

const (
	testWatchBody = `{
		"trigger" : {
			"schedule" : { "cron" : "0 0/1 * * * ?" }
		},
		"input" : {
			"search" : {
				"request" : {
					"indices" : [
						"elastic-test"
					],
					"body" : {
						"query" : {
							"bool" : {
								"must" : {
									"match": {
										 "response": 404
									}
								},
								"filter" : {
									"range": {
										"@timestamp": {
											"from": "{{ctx.trigger.scheduled_time}}||-5m",
											"to": "{{ctx.trigger.triggered_time}}"
										}
									}
								}
							}
						}
					}
				}
			}
		},
		"condition" : {
			"compare" : { "ctx.payload.hits.total" : { "gt" : 0 }}
		},
		"actions" : {
			"email_admin" : {
				"email" : {
					"to" : "admin@domain.host.com",
					"subject" : "404 recently encountered"
				}
			}
		}
	}`
)

func TestXPackWatcher(t *testing.T) {
	client := setupTestClientAndCreateIndex(t, SetURL("http://elastic:elastic@localhost:9210"))

	// Add a watch
	watchName := "my-watch"
	_, err := client.XPackWatchPut(watchName).Body(testWatchBody).Do(context.Background())
	if err != nil {
		if IsForbidden(err) {
			t.Skipf("skip due to missing license: %v", err)
		}
		t.Fatal(err)
	}
	defer func() {
		client.XPackWatchDelete(watchName).Do(context.Background())
	}()

	// Get a watch
	watch, err := client.XPackWatchGet(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if watch.Found == false {
		t.Errorf("expected watch.Found == true; got false")
	}
	if want, have := watchName, watch.Id; want != have {
		t.Errorf("expected watch.Id == %q; got %q", want, have)
	}

	// Exec a watch
	execution, err := client.XPackWatchExecute().Id(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := watchName, execution.WatchRecord.WatchId; want != have {
		t.Errorf("expected execution.WatchId == %q; got %q", want, have)
	}
	if want, have := "execution_not_needed", execution.WatchRecord.State; want != have {
		t.Errorf("expected execution.state == %q; got %q", want, have)
	}

	// Ack a watch
	ack, err := client.XPackWatchAck(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ack.Status.State == nil {
		t.Errorf("expected ack.status != nil; got %v", ack.Status.State)
	}

	// Activate a watch
	_, err = client.XPackWatchActivate(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	watch, err = client.XPackWatchGet(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := true, watch.Status.State["active"]; want != have {
		t.Errorf("expected watch.Status.State.Active == %v; got %v", want, have)
	}

	// Deactivate the watch
	_, err = client.XPackWatchDeactivate(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	watch, err = client.XPackWatchGet(watchName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := false, watch.Status.State["active"]; want != have {
		t.Errorf("expected watch.Status.State.Active == %v; got %v", want, have)
	}

	// Stop the watch
	_, err = client.XPackWatchStop().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stats, err := client.XPackWatchStats().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if have := stats.Stats[0].WatcherState; have != "stopping" && have != "stopped" {
		t.Errorf("expected stats.WatcherState == %q (or %q); got %q", "stopping", "stopped", have)
	}

	// Start again
	start, err := client.XPackWatchStart().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stats, err = client.XPackWatchStats().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := true, start.Acknowledged; want != have {
		t.Errorf("expected start.Acknowledged == %v; got %v", want, have)
	}

	// Restart
	restart, err := client.XPackWatchRestart().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := true, restart.Acknowledged; want != have {
		t.Errorf("expected stats.WatcherState == %v; got %v", want, have)
	}
}
