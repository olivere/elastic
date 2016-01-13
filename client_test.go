// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// -- NewClient --

func TestClientDefaults(t *testing.T) {
	client, err := NewClient(SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}
	if client.basicAuth != false {
		t.Errorf("expected no basic auth; got: %v", client.basicAuth)
	}
	if client.basicAuthUsername != "" {
		t.Errorf("expected no basic auth username; got: %q", client.basicAuthUsername)
	}
	if client.basicAuthPassword != "" {
		t.Errorf("expected no basic auth password; got: %q", client.basicAuthUsername)
	}
	if client.sendGetBodyAs != "GET" {
		t.Errorf("expected sendGetBodyAs to be GET; got: %q", client.sendGetBodyAs)
	}
}

func TestClientWithoutURL(t *testing.T) {
	_, err := NewClient( /* no URL specified */ )
	if err == nil {
		t.Fatal("expected error when no URL specified")
	}
}

func TestClientWithSingleURL(t *testing.T) {
	_, err := NewClient(SetURL("http://localhost:9200"))
	if err != nil {
		t.Fatal(err)
	}
	if !isTravis() {
		//TODO: ping
	}
}

func TestClientWithBasicAuth(t *testing.T) {
	client, err := NewClient(SetBasicAuth("user", "secret"), SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}
	if client.basicAuth != true {
		t.Errorf("expected basic auth; got: %v", client.basicAuth)
	}
	if got, want := client.basicAuthUsername, "user"; got != want {
		t.Errorf("expected basic auth username %q; got: %q", want, got)
	}
	if got, want := client.basicAuthPassword, "secret"; got != want {
		t.Errorf("expected basic auth password %q; got: %q", want, got)
	}
}

// -- ElasticsearchVersion --

func TestElasticsearchVersion(t *testing.T) {
	client, err := NewClient(SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}
	version, err := client.ElasticsearchVersion(DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Errorf("expected a version number, got: %q", version)
	}
}

// -- IndexNames --

func TestIndexNames(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)
	names, err := client.IndexNames()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatalf("expected some index names, got: %d", len(names))
	}
	var found bool
	for _, name := range names {
		if name == testIndexName {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find index %q; got: %v", testIndexName, found)
	}
}

// -- PerformRequest --

func TestPerformRequest(t *testing.T) {
	client, err := NewClient(SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest("GET", "/", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}

	ret := new(PingResult)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		t.Fatalf("expected no error on decode; got: %v", err)
	}
	if ret.Status != 200 {
		t.Errorf("expected HTTP status 200; got: %d", ret.Status)
	}
}

func TestPerformRequestWithLogger(t *testing.T) {
	var w bytes.Buffer
	out := log.New(&w, "LOGGER ", log.LstdFlags)

	client, err := NewClient(SetInfoLog(out), SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.PerformRequest("GET", "/", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}

	ret := new(PingResult)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		t.Fatalf("expected no error on decode; got: %v", err)
	}
	if ret.Status != 200 {
		t.Errorf("expected HTTP status 200; got: %d", ret.Status)
	}

	got := w.String()
	pattern := `^LOGGER \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} GET http://.*/ \[status:200, request:\d+\.\d{3}s\]\n`
	matched, err := regexp.MatchString(pattern, got)
	if err != nil {
		t.Fatalf("expected log line to match %q; got: %v", pattern, err)
	}
	if !matched {
		t.Errorf("expected log line to match %q; got: %v", pattern, got)
	}
}

func TestPerformRequestWithLoggerAndTracer(t *testing.T) {
	var lw bytes.Buffer
	lout := log.New(&lw, "LOGGER ", log.LstdFlags)

	var tw bytes.Buffer
	tout := log.New(&tw, "TRACER ", log.LstdFlags)

	client, err := NewClient(SetInfoLog(lout), SetTraceLog(tout), SetURL(DefaultURL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.PerformRequest("GET", "/", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}

	ret := new(PingResult)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		t.Fatalf("expected no error on decode; got: %v", err)
	}
	if ret.Status != 200 {
		t.Errorf("expected HTTP status 200; got: %d", ret.Status)
	}

	lgot := lw.String()
	if lgot == "" {
		t.Errorf("expected logger output; got: %q", lgot)
	}

	tgot := tw.String()
	if tgot == "" {
		t.Errorf("expected tracer output; got: %q", tgot)
	}
}

// failingTransport will run a fail callback if it sees a given URL path prefix.
type failingTransport struct {
	path string                                      // path prefix to look for
	fail func(*http.Request) (*http.Response, error) // call when path prefix is found
	next http.RoundTripper                           // next round-tripper (use http.DefaultTransport if nil)
}

// RoundTrip implements a failing transport.
func (tr *failingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if strings.HasPrefix(r.URL.Path, tr.path) && tr.fail != nil {
		return tr.fail(r)
	}
	if tr.next != nil {
		return tr.next.RoundTrip(r)
	}
	return http.DefaultTransport.RoundTrip(r)
}

// CancelRequest is required in a http.Transport to support timeouts.
func (tr *failingTransport) CancelRequest(req *http.Request) {
}
