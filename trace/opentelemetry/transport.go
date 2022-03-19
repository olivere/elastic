// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opentelemetry

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Transport for tracing Elastic operations.
type Transport struct {
	rt http.RoundTripper
}

// Option signature for specifying options, e.g. WithRoundTripper.
type Option func(t *Transport)

// WithRoundTripper specifies the http.RoundTripper to call
// next after this transport. If it is nil (default), the
// transport will use http.DefaultTransport.
func WithRoundTripper(rt http.RoundTripper) Option {
	return func(t *Transport) {
		t.rt = rt
	}
}

// NewTransport specifies a transport that will trace Elastic
// and report back via OpenTracing.
func NewTransport(opts ...Option) *Transport {
	t := &Transport{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// RoundTrip captures the request and starts an OpenTracing span
// for Elastic PerformRequest operation.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := otel.Tracer("Elastic").Start(req.Context(), "PerformRequest")
	defer span.End()

	req = req.WithContext(ctx)

	// See General (https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/span-general.md)
	// and HTTP (https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md)
	span.SetAttributes(
		attribute.String("code.namespace", "github.com/olivere/elastic/v7"),
		attribute.String("code.function", "PerformRequest"),
		attribute.String("http.url", req.URL.Redacted()),
		attribute.String("http.method", req.Method),
		attribute.String("http.scheme", req.URL.Scheme),
		attribute.String("http.host", req.URL.Hostname()),
		attribute.String("http.path", req.URL.Path),
		attribute.String("http.user_agent", req.UserAgent()),
	)

	var (
		resp *http.Response
		err  error
	)
	if t.rt != nil {
		resp, err = t.rt.RoundTrip(req)
	} else {
		resp, err = http.DefaultTransport.RoundTrip(req)
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	if resp != nil {
		span.SetAttributes(attribute.Int64("http.status_code", int64(resp.StatusCode)))
	}

	return resp, err
}
