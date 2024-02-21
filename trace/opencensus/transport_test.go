// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opencensus

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/disaster37/opensearch/v2"
	"go.opencensus.io/trace"
)

func init() {
	// Always sample
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

type testExporter struct {
	spans []*trace.SpanData
}

func (t *testExporter) ExportSpan(s *trace.SpanData) {
	t.spans = append(t.spans, s)
}

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "alice" || password != "secret" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

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

	// Register test exporter
	var te testExporter
	trace.RegisterExporter(&te)

	// Setup a simple transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	tr := NewTransport(
		WithDefaultAttributes(
			trace.StringAttribute("Opaque-Id", "12345"),
		),
		WithRoundTripper(transport),
	)
	httpClient := &http.Client{
		Transport: tr,
	}

	// Create a simple Ping request via Elastic
	client, err := opensearch.NewClient(
		opensearch.SetURL(ts.URL),
		opensearch.SetHttpClient(httpClient),
		opensearch.SetHealthcheck(false),
		opensearch.SetSniff(false),
		opensearch.SetBasicAuth("alice", "secret"),
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
	trace.UnregisterExporter(&te)

	// Check the data written into tracer
	spans := te.spans
	if want, have := 1, len(spans); want != have {
		t.Fatalf("want %d finished spans, have %d", want, have)
	}
	span := spans[0]

	if want, have := "opensearch:PerformRequest", span.Name; want != have {
		t.Fatalf("want Span.Name=%q, have %q", want, have)
	}
	if attr, ok := span.Attributes["Component"].(string); !ok {
		t.Fatalf("attribute %q not found", "Component")
	} else if want, have := "github.com/disaster37/opensearch/v2", attr; want != have {
		t.Fatalf("want attribute=%q, have %q", want, have)
	}
	if attr, ok := span.Attributes["Method"].(string); !ok {
		t.Fatalf("attribute %q not found", "Method")
	} else if want, have := "GET", attr; want != have {
		t.Fatalf("want attribute=%q, have %q", want, have)
	}
	if attr, ok := span.Attributes["URL"].(string); !ok || attr == "" {
		t.Fatalf("attribute %q not found", "URL")
	} else if strings.Contains(attr, "alice") || strings.Contains(attr, "password") {
		t.Fatalf("attribute %q contains username and/or password: %s", "URL", attr)
	}
	if attr, ok := span.Attributes["Hostname"].(string); !ok || attr == "" {
		t.Fatalf("attribute %q not found", "Hostname")
	}
	if port, ok := span.Attributes["Port"].(int64); !ok || port <= 0 {
		t.Fatalf("attribute %q not found", "Port")
	}
	if attr, ok := span.Attributes["Opaque-Id"].(string); !ok {
		t.Fatalf("attribute %q not found", "Opaque-Id")
	} else if want, have := "12345", attr; want != have {
		t.Fatalf("want attribute=%q, have %q", want, have)
	}
}
