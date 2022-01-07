// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opentracing

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/olivere/elastic/v7"
)

func TestTransportIntegration(t *testing.T) {
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
		elastic.SetURL("http://127.0.0.1:9210"),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "elastic"),
		elastic.SetHttpClient(httpClient),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Search("_all").Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		t.Fatal(err)
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
	if want, have := "github.com/olivere/elastic/v7", span.Tag("component"); want != have {
		t.Fatalf("want component tag=%q, have %q", want, have)
	}
	httpURL, ok := span.Tag("http.url").(string)
	if !ok || httpURL == "" {
		t.Fatalf("want http.url tag=%q to be a non-empty string (found type %T)", "http.url", span.Tag("http.url"))
	}
	if want, have := "http://127.0.0.1:9210/_all/_search", httpURL; want != have {
		t.Fatalf("want http.url tag=%q, have %q", want, have)
	}
	if strings.Contains(httpURL, "elastic") {
		t.Fatalf("want http.url tag %q to not contain username and/or password: %s", "URL", span.Tag("http.url"))
	}
	if want, have := "POST", span.Tag("http.method"); want != have {
		t.Fatalf("want http.method tag=%q, have %q", want, have)
	}
	if want, have := uint16(http.StatusOK), span.Tag("http.status_code"); want != have {
		t.Fatalf("want http.status_code tag=%v (%T), have %v (%T)", want, want, have, have)
	}
}
