package elastic

import (
	"context"
	"testing"
)

func TestXPackWatchWorkFlow(t *testing.T) {
	client := setupTestClientAndCreateIndex(t, SetURL("http://elastic:elastic@localhost:9210"))

	watchName := "my-watch"
	watchBody := getWatchBody()
	client.XPackWatchPut().Id(watchName).BodyString(watchBody).Do(context.TODO())

	watch, err := client.XPackWatchGet().Id(watchName).Do(context.TODO())
	if err != nil {
		if IsForbidden(err) {
			t.Skipf("skip due to missing license: %v", err)
		}
		t.Fatal(err)
	}
	if watch.Found == false {
		t.Errorf("expected watch.Found == true; got false")
	}
	if watch.Id != watchName {
		t.Errorf("expected watch.Id == %s; got %s", watchName, watch.Id)
	}

	execution, err := client.XPackWatchExecute().Id(watchName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if execution.WatchRecord.WatchId != watchName {
		t.Errorf("expected execution.WatchId == %s; got %s", watchName, execution.WatchRecord.WatchId)
	}
	if execution.WatchRecord.State != "execution_not_needed" {
		t.Errorf("expected execution.state == %s; got %s", "execution_not_needed", execution.WatchRecord.State)
	}

	ack, _ := client.XPackWatchAck().WatchId(watchName).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if ack.Status.State == nil {
		t.Errorf("expected ack.status != %s; got %s", "nil", ack.Status.State)
	}

	client.XPackWatchActivate().WatchId(watchName).Do(context.TODO())
	watch, err = client.XPackWatchGet().Id(watchName).Do(context.TODO())

	if err != nil {
		t.Fatal(err)
	}
	if watch.Status.State["active"] != true {
		t.Errorf("expected watch.Status.State.Active == %t; got %t", true, watch.Status.State["active"])
	}

	client.XPackWatchDeactivate().WatchId(watchName).Do(context.TODO())
	watch, err = client.XPackWatchGet().Id(watchName).Do(context.TODO())

	if err != nil {
		t.Fatal(err)
	}
	if watch.Status.State["active"] != false {
		t.Errorf("expected watch.Status.State.Active == %t; got %t", false, watch.Status.State["active"])
	}

	client.XPackWatchStop().Do(context.TODO())

	stats, err := client.XPackWatchStats().Do(context.TODO())

	if err != nil {
		t.Fatal(err)
	}
	if stats.Stats[0].WatcherState != "stopping" {
		t.Errorf("expected stats.WatcherState == %s; got %s", "stopping", stats.Stats[0].WatcherState)
	}

	start, err := client.XPackWatchStart().Do(context.TODO())

	stats, err = client.XPackWatchStats().Do(context.TODO())

	if err != nil {
		t.Fatal(err)
	}
	if start.Acknowledged != true {
		t.Errorf("expected start.Acknowledged == %t; got %t", true, start.Acknowledged)
	}

	restart, err := client.XPackWatchStart().Do(context.TODO())

	if err != nil {
		t.Fatal(err)
	}
	if restart.Acknowledged != true {
		t.Errorf("expected stats.WatcherState == %t; got %t", true, restart.Acknowledged)
	}

}

func getWatchBody() string {
	return `
{
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
}
`
}
