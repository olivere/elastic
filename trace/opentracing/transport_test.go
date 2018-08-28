// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opentracing

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/olivere/elastic"
)

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{
				"name" : "Qg28M36",
				"cluster_name" : "docker-cluster",
				"cluster_uuid" : "rwHa7BBnRC2h8KoDfCbmuQ",
				"version" : {
					"number" : "6.3.2",
					"build_flavor" : "oss",
					"build_type" : "tar",
					"build_hash" : "053779d",
					"build_date" : "2018-07-20T05:20:23.451332Z",
					"build_snapshot" : false,
					"lucene_version" : "7.3.1",
					"minimum_wire_compatibility_version" : "5.6.0",
					"minimum_index_compatibility_version" : "5.0.0"
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

	// Mock tracer
	tracer := mocktracer.New()
	opentracing.InitGlobalTracer(tracer)

	// Setup a simple transport
	tr := NewTransport()
	httpClient := &http.Client{
		Transport: tr,
	}

	// Create a simple Ping request via Elastic
	client, err := elastic.NewClient(
		elastic.SetURL(ts.URL),
		elastic.SetHttpClient(httpClient),
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

	// Check the data written into tracer
	spans := tracer.FinishedSpans()
	if want, have := 1, len(spans); want != have {
		t.Fatalf("want %d finished spans, have %d", want, have)
	}
	span := spans[0]
	if want, have := "PerformRequest", span.OperationName; want != have {
		t.Fatalf("want Span.OperationName=%q, have %q", want, have)
	}
	if want, have := "github.com/olivere/elastic/v6", span.Tag("component"); want != have {
		t.Fatalf("want component tag=%q, have %q", want, have)
	}
	if want, have := ts.URL+"/", span.Tag("http.url"); want != have {
		t.Fatalf("want http.url tag=%q, have %q", want, have)
	}
	if want, have := "GET", span.Tag("http.method"); want != have {
		t.Fatalf("want http.method tag=%q, have %q", want, have)
	}
	if want, have := uint16(http.StatusOK), span.Tag("http.status_code"); want != have {
		t.Fatalf("want http.status_code tag=%v (%T), have %v (%T)", want, want, have, have)
	}
}
