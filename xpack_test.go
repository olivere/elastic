package elastic

import (
	"context"
	"testing"
)

const (
	testRoleBody = `{
		"cluster" : [ "all" ],
		"indices" : [
			{
				"names" : [ "index1", "index2" ],
				"privileges" : [ "all" ],
				"field_security" : {
					"grant" : [ "title", "body" ]
				}
			}
		],
		"applications" : [ ],
		"run_as" : [ "other_user" ],
		"global" : {
			"application": {
			  "manage": {    
				  "applications": [ "my-test-app" ] 
			  }
			}
		  },
		"metadata" : {
			"version" : 1
		},
		"transient_metadata": {
			"enabled": true
		}
	  }`

	testRoleMappingBody = `{
		"enabled": false,
		"roles": [
			"user"
		],
		"rules": {
			"all": [
				{
					"field": {
					"username": "esadmin"
					}
				},
				{
					"field": {
					"groups": "cn=admins,dc=example,dc=com"
					}
				}
			]
		},
		"metadata": {
			"version": 1
		}
	  }`

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

func TestXpackInfo(t *testing.T) {
	client := setupTestClientForXpackSecurity(t)
	tagline := "You know, for X"

	// Get xpack info
	info, err := client.XPackInfo().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info == &(XPackInfoServiceResponse{}) {
		t.Errorf("expected data from response; got empty response")
	}
	if info.Tagline != tagline {
		t.Errorf("expected %s as a tagline; received %s", tagline, info.Tagline)
	}
}

func TestXPackSecurityRole(t *testing.T) {
	client := setupTestClientForXpackSecurity(t)

	xpack_info, err := client.XPackInfo().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !xpack_info.Features.Security.Enabled {
		t.Skip("skip due to deactivated xpack security")
	}

	roleName := "my-role"

	// Add a role
	_, err = client.XPackSecurityPutRole(roleName).Body(testRoleBody).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		client.XPackSecurityDeleteRole(roleName).Do(context.Background())
	}()

	// Get a role
	role, err := client.XPackSecurityGetRole(roleName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(*role) == 0 {
		t.Errorf("expected len(Mappings) > 0; got empty")
	}
	if _, ok := (*role)[roleName]; !ok {
		t.Errorf("expected role mapping %s; key did not exist", roleName)
	}
	if role == &(XPackSecurityGetRoleResponse{}) {
		t.Errorf("expected data from response; got empty response")
	}

	// Delete a role
	deletedRole, err := client.XPackSecurityDeleteRole(roleName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !deletedRole.Found {
		t.Error("expected test role to be found; was not found")
	}

}

func TestXPackSecurityRoleMapping(t *testing.T) {
	client := setupTestClientForXpackSecurity(t)

	xpack_info, err := client.XPackInfo().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !xpack_info.Features.Security.Enabled {
		t.Skip("skip due to deactivated xpack security")
	}

	roleMappingName := "my-role-mapping"

	// Add a role mapping
	_, err = client.XPackSecurityPutRoleMapping(roleMappingName).Body(testRoleMappingBody).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		client.XPackSecurityDeleteRoleMapping(roleMappingName).Do(context.Background())
	}()

	// Get a role mapping
	roleMappings, err := client.XPackSecurityGetRoleMapping(roleMappingName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(*roleMappings) == 0 {
		t.Errorf("expected len(Mappings) > 0; got empty")
	}
	if _, ok := (*roleMappings)[roleMappingName]; !ok {
		t.Errorf("expected role mapping %s; key did not exist", roleMappingName)
	}
	if roleMappings == &(XPackSecurityGetRoleMappingResponse{}) {
		t.Errorf("expected data from response; got empty response")
	}

	// Delete a role mapping
	_, err = client.XPackSecurityDeleteRoleMapping(roleMappingName).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

}

func TestXPackWatcher(t *testing.T) {
	client := setupTestClientAndCreateIndex(t, SetURL("http://elastic:elastic@localhost:9210"))

	xpack_info, err := client.XPackInfo().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !xpack_info.Features.Watcher.Enabled {
		t.Skip("skip due to deactivated xpack watcher")
	}

	// Add a watch
	watchName := "my-watch"
	_, err = client.XPackWatchPut(watchName).Body(testWatchBody).Do(context.Background())
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
	if want, have := true, watch.Status.State.Active; want != have {
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
	if want, have := false, watch.Status.State.Active; want != have {
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
	_, err = client.XPackWatchStats().Do(context.Background())
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
