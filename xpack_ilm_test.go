// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"testing"
)

func TestXPackIlmPolicyLifecycle(t *testing.T) {
	client := setupTestClient(t, SetURL("http://elastic:elastic@localhost:9210"))

	testPolicyName := "test-policy"

	body := `{
		"policy": {
			"phases": {
				"delete": {
					"min_age": "20m",
					"actions": {
						"delete": {}
					}
				}
			}
		}
	}`

	// Create the policy
	putilm, err := client.XPackIlmPutLifecycle().Policy(testPolicyName).BodyString(body).Do(context.TODO())
	if err != nil {
		t.Fatalf("expected put lifecycle to succeed; got: %v", err)
	}
	if putilm == nil {
		t.Fatalf("expected put lifecycle response; got: %v", putilm)
	}
	if !putilm.Acknowledged {
		t.Fatalf("expected put lifecycle ack; got: %v", putilm.Acknowledged)
	}

	// Get the policy
	getilm, err := client.XPackIlmGetLifecycle().Policy(testPolicyName).Do(context.TODO())
	if err != nil {
		t.Fatalf("expected get lifecycle to succeed; got: %v", err)
	}
	if getilm == nil {
		t.Fatalf("expected get lifecycle response; got: %v", getilm)
	}

	// Check the policy exists
	_, found := getilm[testPolicyName]
	if !found {
		t.Fatalf("expected to get policy for %q", testPolicyName)
	}

	// Delete the policy
	delilm, err := client.XPackIlmDeleteLifecycle().Policy(testPolicyName).Do(context.TODO())
	if err != nil {
		t.Fatalf("expected deletelifecycle to succeed; got: %v", err)
	}
	if delilm == nil {
		t.Fatalf("expected delete lifecycle response; got: %v", delilm)
	}
	if !delilm.Acknowledged {
		t.Fatalf("expected delete lifecycle ack; got: %v", delilm.Acknowledged)
	}

	// Get the policy
	getilm, err = client.XPackIlmGetLifecycle().Policy(testPolicyName).Do(context.TODO())
	if err == nil {
		t.Fatalf("expected lifecycle to be deleted; got: %v", getilm)
	}
}
