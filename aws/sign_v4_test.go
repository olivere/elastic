// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package aws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/go-aws-auth"

	"gopkg.in/olivere/elastic.v5"
)

func TestSigningClient(t *testing.T) {
	var req *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			req = r // capture the HTTP request
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{
				"name" : "n5VFSHr",
				"cluster_name" : "elasticsearch",
				"cluster_uuid" : "8nwVxw7bTRuzfBT45QIRWA",
				"version" : {
				  "number" : "5.6.10",
				  "build_hash" : "b727a60",
				  "build_date" : "2018-06-06T15:48:34.860Z",
				  "build_snapshot" : false,
				  "lucene_version" : "6.6.1"
				},
				"tagline" : "You Know, for Search"
			  }`)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()

	cred := awsauth.Credentials{
		AccessKeyID:     "dev",
		SecretAccessKey: "secret",
	}
	signingClient := NewV4SigningClient(cred)

	// Create a simple Ping request via Elastic
	client, err := elastic.NewClient(
		elastic.SetURL(ts.URL),
		elastic.SetHttpClient(signingClient),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		t.Fatal(err)
	}
	res, code, err := client.Ping(ts.URL).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want, have := http.StatusOK, code; want != have {
		t.Fatalf("want Status=%d, have %d", want, have)
	}
	if want, have := "You Know, for Search", res.TagLine; want != have {
		t.Fatalf("want TagLine=%q, have %q", want, have)
	}

	// Check the request recorded in the HTTP test server (see above)
	if req == nil {
		t.Fatal("expected to capture HTTP request")
	}
	if have := req.Header.Get("Authorization"); have == "" {
		t.Fatal("expected Authorization header")
	}
	if have := req.Header.Get("X-Amz-Date"); have == "" {
		t.Fatal("expected X-Amz-Date header")
	}
	if want, have := `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`, req.Header.Get("X-Amz-Content-Sha256"); want != have {
		t.Fatalf("want header of X-Amz-Content-Sha256=%q, have %q", want, have)
	}
}
