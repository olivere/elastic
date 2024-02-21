// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package opensearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/sirupsen/logrus"

	"github.com/disaster37/opensearch/v2/config"
)

func findConn(s string, slice ...*conn) (int, bool) {
	for i, t := range slice {
		if s == t.URL() {
			return i, true
		}
	}
	return -1, false
}

// -- NewClient --

func TestClientDefaults(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	if client.healthcheckEnabled != true {
		t.Errorf("expected health checks to be enabled, got: %v", client.healthcheckEnabled)
	}
	if client.healthcheckTimeoutStartup != DefaultHealthcheckTimeoutStartup {
		t.Errorf("expected health checks timeout on startup = %v, got: %v", DefaultHealthcheckTimeoutStartup, client.healthcheckTimeoutStartup)
	}
	if client.healthcheckTimeout != DefaultHealthcheckTimeout {
		t.Errorf("expected health checks timeout = %v, got: %v", DefaultHealthcheckTimeout, client.healthcheckTimeout)
	}
	if client.healthcheckInterval != DefaultHealthcheckInterval {
		t.Errorf("expected health checks interval = %v, got: %v", DefaultHealthcheckInterval, client.healthcheckInterval)
	}
	if client.snifferEnabled != true {
		t.Errorf("expected sniffing to be enabled, got: %v", client.snifferEnabled)
	}
	if client.snifferTimeoutStartup != DefaultSnifferTimeoutStartup {
		t.Errorf("expected sniffer timeout on startup = %v, got: %v", DefaultSnifferTimeoutStartup, client.snifferTimeoutStartup)
	}
	if client.snifferTimeout != DefaultSnifferTimeout {
		t.Errorf("expected sniffer timeout = %v, got: %v", DefaultSnifferTimeout, client.snifferTimeout)
	}
	if client.snifferInterval != DefaultSnifferInterval {
		t.Errorf("expected sniffer interval = %v, got: %v", DefaultSnifferInterval, client.snifferInterval)
	}
	if client.sendGetBodyAs != "GET" {
		t.Errorf("expected sendGetBodyAs to be GET; got: %q", client.sendGetBodyAs)
	}
}

func TestClientWithoutURL(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	// Two things should happen here:
	// 1. The client starts sniffing the cluster on DefaultURL
	// 2. The sniffing process should find (at least) one node in the cluster, i.e. the DefaultURL
	if len(client.conns) == 0 {
		t.Fatalf("expected at least 1 node in the cluster, got: %d (%v)", len(client.conns), client.conns)
	}
	if !isCI() {
		if _, found := findConn(DefaultURL, client.conns...); !found {
			t.Errorf("expected to find node with default URL of %s in %v", DefaultURL, client.conns)
		}
	}
}

func TestClientWithSingleURL(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport), SetURL("https://127.0.0.1:9200"))
	if err != nil {
		t.Fatal(err)
	}
	// Two things should happen here:
	// 1. The client starts sniffing the cluster on DefaultURL
	// 2. The sniffing process should find (at least) one node in the cluster, i.e. the DefaultURL
	if len(client.conns) == 0 {
		t.Fatalf("expected at least 1 node in the cluster, got: %d (%v)", len(client.conns), client.conns)
	}
	if !isCI() {
		if _, found := findConn(DefaultURL, client.conns...); !found {
			t.Errorf("expected to find node with default URL of %s in %v", DefaultURL, client.conns)
		}
	}
}

func TestClientWithMultipleURLs(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport), SetURL("https://127.0.0.1:9200", "https://127.0.0.1:9201"))
	if err != nil {
		t.Fatal(err)
	}
	// The client should sniff both URLs, but only 127.0.0.1:9200 should return nodes.
	if len(client.conns) != 1 {
		t.Fatalf("expected exactly 1 node in the local cluster, got: %d (%v)", len(client.conns), client.conns)
	}
	if !isCI() {
		if client.conns[0].URL() != DefaultURL {
			t.Errorf("expected to find node with default URL of %s in %v", DefaultURL, client.conns)
		}
	}
}

func TestClientWithInvalidURLs(t *testing.T) {
	client, err := NewClient(SetURL(" http://foo.com", "http://[fe80::%31%25en0]:8080/"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if want, have := `first path segment in URL cannot contain colon`, err.Error(); !strings.Contains(have, want) {
		t.Fatalf("expected error to contain %q, have %q", want, have)
	}
	if client != nil {
		t.Fatal("expected client == nil")
	}
}

func TestClientWithBasicAuth(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := client.basicAuthUsername, "admin"; got != want {
		t.Errorf("expected basic auth username %q; got: %q", want, got)
	}
	if got, want := client.basicAuthPassword, "vLPeJYa8.3RqtZCcAK6jNz"; got != want {
		t.Errorf("expected basic auth password %q; got: %q", want, got)
	}
}

func TestClientWithBasicAuthInUserInfo(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetTransport(transport), SetURL("https://admin:vLPeJYa8.3RqtZCcAK6jNz@127.0.0.1:9200"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := client.basicAuthUsername, "admin"; got != want {
		t.Errorf("expected basic auth username %q; got: %q", want, got)
	}
	if got, want := client.basicAuthPassword, "vLPeJYa8.3RqtZCcAK6jNz"; got != want {
		t.Errorf("expected basic auth password %q; got: %q", want, got)
	}
}

func TestClientWithBasicAuthDuringHealthcheck(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" || r.URL.String() != "/" {
			t.Fatalf("expected HEAD / request, got %s %s", r.Method, r.URL)
			http.Error(w, fmt.Sprintf("expected HEAD / request, got %s %s", r.Method, r.URL), http.StatusBadRequest)
			return
		}
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Fatal("expected HEAD basic auth")
			http.Error(w, "expected HTTP basic auth", http.StatusBadRequest)
			return
		}
		if username != "admin" && password != "vLPeJYa8.3RqtZCcAK6jNz" {
			t.Fatalf("invalid HTTP basic auth username %q and password %q", username, password)
			http.Error(w, fmt.Sprintf("invalid HTTP basic auth username %q and password %q", username, password), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport), SetURL(ts.URL), SetSniff(false))
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected a client")
	}
}

func TestClientWithoutBasicAuthButAuthEnabledInElasticDuringHealthcheck(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" || r.URL.String() != "/" {
			t.Fatalf("expected HEAD / request, got %s %s", r.Method, r.URL)
			http.Error(w, fmt.Sprintf("expected HEAD / request, got %s %s", r.Method, r.URL), http.StatusBadRequest)
			return
		}
		_, _, ok := r.BasicAuth()
		if ok {
			t.Fatal("unexpected HEAD basic auth")
			http.Error(w, "unexpected HTTP basic auth", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	_, err := NewClient(SetURL(ts.URL), SetSniff(false))
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	if !IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized = %v, got err=%+v", true, err)
	}
}

func TestClientFromConfig(t *testing.T) {
	cfg, err := config.Parse("https://127.0.0.1:9200")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	cfg.Username = "admin"
	cfg.Password = "vLPeJYa8.3RqtZCcAK6jNz"
	client, err := NewClientFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Two things should happen here:
	// 1. The client starts sniffing the cluster on DefaultURL
	// 2. The sniffing process should find (at least) one node in the cluster, i.e. the DefaultURL
	if len(client.conns) == 0 {
		t.Fatalf("expected at least 1 node in the cluster, got: %d (%v)", len(client.conns), client.conns)
	}
	if !isCI() {
		if _, found := findConn(DefaultURL, client.conns...); !found {
			t.Errorf("expected to find node with default URL of %s in %v", DefaultURL, client.conns)
		}
	}
}

func TestClientDialFromConfig(t *testing.T) {
	cfg, err := config.Parse("https://127.0.0.1:9200")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	cfg.Username = "admin"
	cfg.Password = "vLPeJYa8.3RqtZCcAK6jNz"
	client, err := DialWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Two things should happen here:
	// 1. The client starts sniffing the cluster on DefaultURL
	// 2. The sniffing process should find (at least) one node in the cluster, i.e. the DefaultURL
	if len(client.conns) == 0 {
		t.Fatalf("expected at least 1 node in the cluster, got: %d (%v)", len(client.conns), client.conns)
	}
	if !isCI() {
		if _, found := findConn(DefaultURL, client.conns...); !found {
			t.Errorf("expected to find node with default URL of %s in %v", DefaultURL, client.conns)
		}
	}
}

func TestClientDialContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := DialContext(ctx, SetURL("https://localhost:9200"), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatalf("expected successful connection, got %v", err)
	}
	client.Stop()
}

func TestClientDialContextTimeoutFromHealthcheck(t *testing.T) {
	start := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	_, err := DialContext(ctx, SetURL("https://localhost:9299"), SetHealthcheckTimeoutStartup(5*time.Second), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if !IsContextErr(err) {
		t.Fatal(err)
	}
	if time.Since(start) < 3*time.Second {
		t.Fatalf("early timeout")
	}
	if time.Since(start) >= 5*time.Second {
		t.Fatalf("timeout probably due to healthcheck, not context cancellation")
	}
}

func TestClientDialContextTimeoutFromSniffer(t *testing.T) {
	start := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	_, err := DialContext(ctx, SetURL("https://localhost:9299"), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if !IsContextErr(err) {
		t.Fatal(err)
	}
	if time.Since(start) < 3*time.Second {
		t.Fatalf("early timeout")
	}
	if time.Since(start) >= 5*time.Second {
		t.Fatalf("timeout probably not caused by context cancellation")
	}
}

func TestClientSniffSuccess(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetURL("https://127.0.0.1:19200", "https://127.0.0.1:9200"), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	// The client should sniff both URLs, but only 127.0.0.1:9200 should return nodes.
	if len(client.conns) != 1 {
		t.Fatalf("expected exactly 1 node in the local cluster, got: %d (%v)", len(client.conns), client.conns)
	}
}

func TestClientSniffFailure(t *testing.T) {
	_, err := NewClient(SetURL("http://127.0.0.1:19200", "http://127.0.0.1:19201"))
	if err == nil {
		t.Fatalf("expected cluster to fail with no nodes found")
	}
}

func TestClientSnifferCallback(t *testing.T) {
	var calls int
	cb := func(node *NodesInfoNode) bool {
		calls++
		return false
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	_, err := NewClient(
		SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"),
		SetTransport(transport),
		SetURL("https://127.0.0.1:19200", "https://127.0.0.1:9200"),
		SetSnifferCallback(cb))
	if err == nil {
		t.Fatalf("expected cluster to fail with no nodes found")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call to the sniffer callback, got %d", calls)
	}
}

func TestClientSniffDisabled(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetSniff(false), SetURL("https://127.0.0.1:9200", "https://127.0.0.1:9201"), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	// The client should not sniff, so it should have two connections.
	if len(client.conns) != 2 {
		t.Fatalf("expected 2 nodes, got: %d (%v)", len(client.conns), client.conns)
	}
	// Make two requests, so that both connections are being used
	for i := 0; i < len(client.conns); i++ {
		client.Flush().Do(context.TODO())
	}
	// The first connection (127.0.0.1:9200) should now be okay.
	if i, found := findConn("https://127.0.0.1:9200", client.conns...); !found {
		t.Fatalf("expected connection to %q to be found", "https://127.0.0.1:9200")
	} else {
		if conn := client.conns[i]; conn.IsDead() {
			t.Fatal("expected connection to be alive, but it is dead")
		}
	}
	// The second connection (127.0.0.1:9201) should now be marked as dead.
	if i, found := findConn("https://127.0.0.1:9201", client.conns...); !found {
		t.Fatalf("expected connection to %q to be found", "https://127.0.0.1:9201")
	} else {
		if conn := client.conns[i]; !conn.IsDead() {
			t.Fatal("expected connection to be dead, but it is alive")
		}
	}
}

func TestClientWillMarkConnectionsAsAliveWhenAllAreDead(t *testing.T) {
	client, err := NewClient(SetURL("http://127.0.0.1:9201"),
		SetSniff(false), SetHealthcheck(false), SetMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	// We should have a connection.
	if len(client.conns) != 1 {
		t.Fatalf("expected 1 node, got: %d (%v)", len(client.conns), client.conns)
	}

	// Make a request, so that the connections is marked as dead.
	client.Flush().Do(context.TODO())

	// The connection should now be marked as dead.
	if i, found := findConn("http://127.0.0.1:9201", client.conns...); !found {
		t.Fatalf("expected connection to %q to be found", "http://127.0.0.1:9201")
	} else {
		if conn := client.conns[i]; !conn.IsDead() {
			t.Fatalf("expected connection to be dead, got: %v", conn)
		}
	}

	// Now send another request and the connection should be marked as alive again.
	client.Flush().Do(context.TODO())

	if i, found := findConn("http://127.0.0.1:9201", client.conns...); !found {
		t.Fatalf("expected connection to %q to be found", "http://127.0.0.1:9201")
	} else {
		if conn := client.conns[i]; conn.IsDead() {
			t.Fatalf("expected connection to be alive, got: %v", conn)
		}
	}
}

func TestClientWithRequiredPlugins(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	_, err := NewClient(SetRequiredPlugins("no-such-plugin"), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err == nil {
		t.Fatal("expected error when creating client")
	}
	if got, want := err.Error(), "opensearch: plugin no-such-plugin not found"; got != want {
		t.Fatalf("expected error %q; got: %q", want, got)
	}
}

func TestClientHealthcheckStartupTimeout(t *testing.T) {
	start := time.Now()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	_, err := NewClient(SetURL("https://localhost:9299"), SetHealthcheckTimeoutStartup(5*time.Second), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	duration := time.Since(start)
	if !IsConnErr(err) {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("expected error to contain %q, have %q", "connection refused", err.Error())
	}
	if duration < 5*time.Second {
		t.Fatalf("expected a timeout in more than 5 seconds; got: %v", duration)
	}
}

func TestClientHealthcheckTimeoutLeak(t *testing.T) {
	// This test test checks if healthcheck requests are canceled
	// after timeout.
	// It contains couple of hacks which won't be needed once we
	// stop supporting Go1.7.
	// On Go1.7 it uses server side effects to monitor if connection
	// was closed,
	// and on Go 1.8+ we're additionally honestly monitoring routine
	// leaks via leaktest.
	mux := http.NewServeMux()

	var reqDoneMu sync.Mutex
	var reqDone bool
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			t.Fatalf("Writer is not CloseNotifier, but %v", reflect.TypeOf(w).Name())
		}
		<-cn.CloseNotify()
		reqDoneMu.Lock()
		reqDone = true
		reqDoneMu.Unlock()
	})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Couldn't setup listener: %v", err)
	}
	addr := lis.Addr().String()

	srv := &http.Server{
		Handler: mux,
	}
	go srv.Serve(lis)

	cli := &Client{
		c: &http.Client{},
		conns: []*conn{
			{
				url: "http://" + addr + "/",
			},
		},
		log: logrus.StandardLogger(),
	}

	type closer interface {
		Shutdown(context.Context) error
	}

	// pre-Go1.8 Server can't Shutdown
	cl, isServerCloseable := (interface{}(srv)).(closer)

	// Since Go1.7 can't Shutdown() - there will be leak from server
	// Monitor leaks on Go 1.8+
	if isServerCloseable {
		defer leaktest.CheckTimeout(t, time.Second*10)()
	}

	cli.healthcheck(context.Background(), time.Millisecond*500, true)

	if isServerCloseable {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cl.Shutdown(ctx)
	}

	<-time.After(time.Second)
	reqDoneMu.Lock()
	if !reqDone {
		reqDoneMu.Unlock()
		t.Fatal("Request wasn't canceled or stopped")
	}
	reqDoneMu.Unlock()
}

func TestClientSniffUpdatingNodeURL(t *testing.T) {
	var (
		nodeID  = "3DWDurZJQvWyWIOFnEB7VA"
		nodeURL string
		n       int
	)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/_nodes/http" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		u, err := url.Parse(nodeURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{
			"cluster_name": "opensearchsearch",
			"nodes": {
				%q: {
					"name": "opensearchsearch",
					"http": {
						"publish_address": %q
					}
				}
			}
		}`, nodeID, u.Host)
		fmt.Fprintln(w)
	})
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	nodeURL = ts.URL

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewSimpleClient(SetURL(ts.URL), SetSniff(true), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 0, n; want != have {
		t.Fatalf("expected %d calls to handler; got %d", want, have)
	}
	if want, have := 1, len(client.conns); want != have {
		t.Fatalf("expected %d connections; got %d", want, have)
	}
	if want, have := nodeURL, client.conns[0].URL(); want != have {
		t.Fatalf("expected URL=%q; got %q", want, have)
	}

	err = client.sniff(context.Background(), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := 1, n; want != have {
		t.Fatalf("expected %d calls to handler; got %d", want, have)
	}
	if want, have := 1, len(client.conns); want != have {
		t.Fatalf("expected %d connections; got %d", want, have)
	}
	if want, have := nodeID, client.conns[0].NodeID(); want != have {
		t.Fatalf("expected NodeID=%q; got %q", want, have)
	}

	if want, have := nodeURL, client.conns[0].URL(); want != have {
		t.Fatalf("expected URL=%q; got %q", want, have)
	}
	oldNodeID := client.conns[0].NodeID()
	oldURL := client.conns[0].URL()

	nodeURL = "https://127.0.0.1:9999" // some other nodeURL to report

	err = client.sniff(context.Background(), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := 2, n; want != have {
		t.Fatalf("expected %d calls to handler; got %d", want, have)
	}
	if want, have := 1, len(client.conns); want != have {
		t.Fatalf("expected %d connections; got %d", want, have)
	}
	newNodeID := client.conns[0].NodeID()
	newURL := client.conns[0].URL()

	// NodeID mustn't change
	if newNodeID != oldNodeID {
		t.Fatalf("expected NodeID=%q; got %q", oldNodeID, newNodeID)
	}
	// URL must have change
	if newURL == oldURL {
		t.Fatalf("expected to update URL=%q to %q", oldURL, newURL)
	}
}

// -- NewSimpleClient --

func TestSimpleClientDefaults(t *testing.T) {
	client, err := NewSimpleClient()
	if err != nil {
		t.Fatal(err)
	}
	if client.healthcheckEnabled != false {
		t.Errorf("expected health checks to be disabled, got: %v", client.healthcheckEnabled)
	}
	if client.healthcheckTimeoutStartup != off {
		t.Errorf("expected health checks timeout on startup = %v, got: %v", off, client.healthcheckTimeoutStartup)
	}
	if client.healthcheckTimeout != off {
		t.Errorf("expected health checks timeout = %v, got: %v", off, client.healthcheckTimeout)
	}
	if client.healthcheckInterval != off {
		t.Errorf("expected health checks interval = %v, got: %v", off, client.healthcheckInterval)
	}
	if client.snifferEnabled != false {
		t.Errorf("expected sniffing to be disabled, got: %v", client.snifferEnabled)
	}
	if client.snifferTimeoutStartup != off {
		t.Errorf("expected sniffer timeout on startup = %v, got: %v", off, client.snifferTimeoutStartup)
	}
	if client.snifferTimeout != off {
		t.Errorf("expected sniffer timeout = %v, got: %v", off, client.snifferTimeout)
	}
	if client.snifferInterval != off {
		t.Errorf("expected sniffer interval = %v, got: %v", off, client.snifferInterval)
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

// -- Start and stop --

func TestClientStartAndStop(t *testing.T) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}

	running := client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}

	// Stop
	client.Stop()
	running = client.IsRunning()
	if running {
		t.Fatalf("expected background processes to be stopped; got: %v", running)
	}

	// Stop again => no-op
	client.Stop()
	running = client.IsRunning()
	if running {
		t.Fatalf("expected background processes to be stopped; got: %v", running)
	}

	// Start
	client.Start()
	running = client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}

	// Start again => no-op
	client.Start()
	running = client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}
}

func TestClientStartAndStopWithSnifferAndHealthchecksDisabled(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetSniff(false), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}

	running := client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}

	// Stop
	client.Stop()
	running = client.IsRunning()
	if running {
		t.Fatalf("expected background processes to be stopped; got: %v", running)
	}

	// Stop again => no-op
	client.Stop()
	running = client.IsRunning()
	if running {
		t.Fatalf("expected background processes to be stopped; got: %v", running)
	}

	// Start
	client.Start()
	running = client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}

	// Start again => no-op
	client.Start()
	running = client.IsRunning()
	if !running {
		t.Fatalf("expected background processes to run; got: %v", running)
	}
}

// -- Sniffing --

func TestClientSniffNode(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan []*conn)
	go func() { ch <- client.sniffNode(context.Background(), DefaultURL) }()

	select {
	case nodes := <-ch:
		if len(nodes) != 1 {
			t.Fatalf("expected %d nodes; got: %d", 1, len(nodes))
		}
		pattern := `https:\/\/[\d\.]+:9200`
		matched, err := regexp.MatchString(pattern, nodes[0].URL())
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Fatalf("expected node URL pattern %q; got: %q", pattern, nodes[0].URL())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected no timeout in sniff node")
		break
	}
}

func TestClientSniffOnDefaultURL(t *testing.T) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, _ := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if client == nil {
		t.Fatal("no client returned")
	}

	ch := make(chan error, 1)
	go func() {
		ch <- client.sniff(context.Background(), DefaultSnifferTimeoutStartup)
	}()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("expected sniff to succeed; got: %v", err)
		}
		if len(client.conns) != 1 {
			t.Fatalf("expected %d nodes; got: %d", 1, len(client.conns))
		}
		pattern := `https:\/\/[\d\.]+:9200`
		matched, err := regexp.MatchString(pattern, client.conns[0].URL())
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Fatalf("expected node URL pattern %q; got: %q", pattern, client.conns[0].URL())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected no timeout in sniff")
		break
	}
}

func TestClientSniffTimeoutLeak(t *testing.T) {
	// This test test checks if sniff requests are canceled
	// after timeout.
	// It contains couple of hacks which won't be needed once we
	// stop supporting Go1.7.
	// On Go1.7 it uses server side effects to monitor if connection
	// was closed,
	// and on Go 1.8+ we're additionally honestly monitoring routine
	// leaks via leaktest.
	mux := http.NewServeMux()

	var reqDoneMu sync.Mutex
	var reqDone bool
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			t.Fatalf("Writer is not CloseNotifier, but %v", reflect.TypeOf(w).Name())
		}
		<-cn.CloseNotify()
		reqDoneMu.Lock()
		reqDone = true
		reqDoneMu.Unlock()
	})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Couldn't setup listener: %v", err)
	}
	addr := lis.Addr().String()

	srv := &http.Server{
		Handler: mux,
	}
	go srv.Serve(lis)

	cli := &Client{
		c: &http.Client{},
		conns: []*conn{
			{
				url: "http://" + addr + "/",
			},
		},
		snifferEnabled: true,
	}

	type closer interface {
		Shutdown(context.Context) error
	}

	// pre-Go1.8 Server can't Shutdown
	cl, isServerCloseable := (interface{}(srv)).(closer)

	// Since Go1.7 can't Shutdown() - there will be leak from server
	// Monitor leaks on Go 1.8+
	if isServerCloseable {
		defer leaktest.CheckTimeout(t, time.Second*10)()
	}

	cli.sniff(context.Background(), time.Millisecond*500)

	if isServerCloseable {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cl.Shutdown(ctx)
	}

	<-time.After(time.Second)
	reqDoneMu.Lock()
	if !reqDone {
		reqDoneMu.Unlock()
		t.Fatal("Request wasn't canceled or stopped")
	}
	reqDoneMu.Unlock()
}

func TestClientExtractHostname(t *testing.T) {
	tests := []struct {
		Scheme  string
		Address string
		Output  string
	}{
		{
			Scheme:  "http",
			Address: "127.0.0.1:9200",
			Output:  "http://127.0.0.1:9200",
		},
		{
			Scheme:  "https",
			Address: "127.0.0.1:9200",
			Output:  "https://127.0.0.1:9200",
		},
		{
			Scheme:  "http",
			Address: "127.0.0.1:19200",
			Output:  "http://127.0.0.1:19200",
		},
		{
			Scheme:  "http",
			Address: "myelk.local/10.1.0.24:9200",
			Output:  "http://myelk.local:9200",
		},
	}

	client, err := NewClient(SetSniff(false), SetHealthcheck(false))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		got := client.extractHostname(test.Scheme, test.Address)
		if want := test.Output; want != got {
			t.Errorf("expected %q; got: %q", want, got)
		}
	}
}

// -- Selector --

func TestClientSelectConnHealthy(t *testing.T) {
	client, err := NewClient(
		SetSniff(false),
		SetHealthcheck(false),
		SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"))
	if err != nil {
		t.Fatal(err)
	}

	// Both are healthy, so we should get both URLs in round-robin
	client.conns[0].MarkAsHealthy()
	client.conns[1].MarkAsHealthy()

	// #1: Return 1st
	c, err := client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[0].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[0].URL())
	}
	// #2: Return 2nd
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[1].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[1].URL())
	}
	// #3: Return 1st
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[0].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[0].URL())
	}
}

func TestClientSelectConnHealthyAndDead(t *testing.T) {
	client, err := NewClient(
		SetSniff(false),
		SetHealthcheck(false),
		SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"))
	if err != nil {
		t.Fatal(err)
	}

	// 1st is healthy, second is dead
	client.conns[0].MarkAsHealthy()
	client.conns[1].MarkAsDead()

	// #1: Return 1st
	c, err := client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[0].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[0].URL())
	}
	// #2: Return 1st again
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[0].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[0].URL())
	}
	// #3: Return 1st again and again
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[0].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[0].URL())
	}
}

func TestClientSelectConnDeadAndHealthy(t *testing.T) {
	client, err := NewClient(
		SetSniff(false),
		SetHealthcheck(false),
		SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"))
	if err != nil {
		t.Fatal(err)
	}

	// 1st is dead, 2nd is healthy
	client.conns[0].MarkAsDead()
	client.conns[1].MarkAsHealthy()

	// #1: Return 2nd
	c, err := client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[1].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[1].URL())
	}
	// #2: Return 2nd again
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[1].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[1].URL())
	}
	// #3: Return 2nd again and again
	c, err = client.next()
	if err != nil {
		t.Fatal(err)
	}
	if c.URL() != client.conns[1].URL() {
		t.Fatalf("expected %s; got: %s", c.URL(), client.conns[1].URL())
	}
}

func TestClientSelectConnAllDead(t *testing.T) {
	client, err := NewClient(
		SetSniff(false),
		SetHealthcheck(false),
		SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"))
	if err != nil {
		t.Fatal(err)
	}

	// Both are dead
	client.conns[0].MarkAsDead()
	client.conns[1].MarkAsDead()

	// If all connections are dead, next should make them alive again, but
	// still return an error when it first finds out.
	c, err := client.next()
	if !IsConnErr(err) {
		t.Fatal(err)
	}
	if c != nil {
		t.Fatalf("expected no connection; got: %v", c)
	}
	// Return a connection
	c, err = client.next()
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if c == nil {
		t.Fatalf("expected connection; got: %v", c)
	}
	// Return a connection
	c, err = client.next()
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if c == nil {
		t.Fatalf("expected connection; got: %v", c)
	}
}

// -- OpensearchVersion --

func TestOpensearchVersion(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	version, err := client.OpensearchVersion(DefaultURL)
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
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
	})
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
	if ret.ClusterName == "" {
		t.Errorf("expected cluster name; got: %q", ret.ClusterName)
	}
}

func TestPerformRequestWithStream(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
		Stream: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
	if res.BodyReader == nil {
		t.Fatal("expected res.BodyReader to be != nil")
	}
	if res.Body != nil {
		t.Fatal("expected res.Body to be == nil")
	}

	ret := new(PingResult)
	if err := json.NewDecoder(res.BodyReader).Decode(ret); err != nil {
		t.Fatalf("expected no error on decode; got: %v", err)
	}
	if ret.ClusterName == "" {
		t.Errorf("expected cluster name; got: %q", ret.ClusterName)
	}
}

func TestPerformRequestWithSimpleClient(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewSimpleClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
	})
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
	if ret.ClusterName == "" {
		t.Errorf("expected cluster name; got: %q", ret.ClusterName)
	}
}

func TestPerformRequestWithMaxResponseSize(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method:          "GET",
		Path:            "/",
		MaxResponseSize: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}

	_, err = client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method:          "GET",
		Path:            "/",
		MaxResponseSize: 100,
	})
	if err != ErrResponseSize {
		t.Fatal("expected response size error")
	}
}

func TestPerformRequestOnNoConnectionsWithHealthcheckRevival(t *testing.T) {
	fail := func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("request failed")
	}
	tr := &failingTransport{
		path: "/fail",
		fail: fail,
	}

	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(SetHttpClient(httpClient), SetMaxRetries(0), SetHealthcheck(true), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}

	// Run against a failing endpoint to mark connection as dead
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/fail",
	})
	if err == nil {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatal("expected no response")
	}

	// Forced healthcheck should bring connection back to life and complete request
	res, err = client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
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

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return t.RoundTrip(r)
}

func TestPerformRequestRetryOnHttpError(t *testing.T) {
	var numFailedReqs int
	fail := func(r *http.Request) (*http.Response, error) {
		numFailedReqs += 1
		return nil, errors.New("request failed")
	}

	// Run against a failing endpoint and see if PerformRequest
	// retries correctly.
	tr := &failingTransport{path: "/fail", fail: fail}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(SetHttpClient(httpClient), SetMaxRetries(5), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/fail",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if res != nil {
		t.Fatal("expected no response")
	}
	// Connection should be marked as dead after it failed
	if numFailedReqs != 5 {
		t.Errorf("expected %d failed requests; got: %d", 5, numFailedReqs)
	}
}

func TestPerformRequestNoRetryOnValidButUnsuccessfulHttpStatus(t *testing.T) {
	var numFailedReqs int
	fail := func(r *http.Request) (*http.Response, error) {
		numFailedReqs++
		return &http.Response{Request: r, StatusCode: 500, Body: http.NoBody}, nil
	}

	// Run against a failing endpoint and see if PerformRequest
	// retries correctly.
	tr := &failingTransport{path: "/fail", fail: fail}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(SetHttpClient(httpClient), SetMaxRetries(5), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/fail",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if res == nil {
		t.Fatal("expected response, got nil")
	}
	if want, got := 500, res.StatusCode; want != got {
		t.Fatalf("expected status code = %d, got %d", want, got)
	}
	// Retry should not have triggered additional requests because
	if numFailedReqs != 1 {
		t.Errorf("expected %d failed requests; got: %d", 1, numFailedReqs)
	}
}

func TestPerformRequestOnSpecifiedHttpStatusCodes(t *testing.T) {
	var numFailedReqs int
	fail := func(r *http.Request) (*http.Response, error) {
		numFailedReqs++
		return &http.Response{Request: r, StatusCode: 429, Body: http.NoBody}, nil
	}

	// Run against a failing endpoint and see if PerformRequest
	// retries correctly.
	tr := &failingTransport{path: "/fail", fail: fail}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(
		SetHttpClient(httpClient),
		SetMaxRetries(5),
		SetRetryStatusCodes(429, 504),
		SetHealthcheck(false),
		SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"),
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/fail",
	})
	if err == nil {
		t.Fatalf("expected error: err=%v, resp=%+v", err, res)
	}
	if res == nil {
		t.Fatal("expected response, got nil")
	}
	if want, got := 429, res.StatusCode; want != got {
		t.Fatalf("expected status code = %d, got %d", want, got)
	}
	// Retry should not have triggered additional requests because
	if numFailedReqs != 5 {
		t.Errorf("expected %d failed requests; got: %d", 1, numFailedReqs)
	}
}

// failingBody will return an error when json.Marshal is called on it.
type failingBody struct{}

// MarshalJSON implements the json.Marshaler interface and always returns an error.
func (fb failingBody) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failing to marshal")
}

func TestPerformRequestWithSetBodyError(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
		Body:   failingBody{},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if res != nil {
		t.Fatal("expected no response")
	}
}

// sleepingTransport will sleep before doing a request.
type sleepingTransport struct {
	timeout time.Duration
}

// RoundTrip implements a "sleepy" transport.
func (tr *sleepingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	time.Sleep(tr.timeout)
	return http.DefaultTransport.RoundTrip(r)
}

func TestPerformRequestWithCancel(t *testing.T) {
	tr := &sleepingTransport{timeout: 3 * time.Second}
	httpClient := &http.Client{Transport: tr}

	client, err := NewSimpleClient(SetHttpClient(httpClient), SetMaxRetries(0), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}

	type result struct {
		res *Response
		err error
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resc := make(chan result, 1)
	go func() {
		res, err := client.PerformRequest(ctx, PerformRequestOptions{
			Method: "GET",
			Path:   "/",
		})
		resc <- result{res: res, err: err}
	}()
	select {
	case <-time.After(1 * time.Second):
		cancel()
	case res := <-resc:
		t.Fatalf("expected response before cancel, got %v", res)
	case <-ctx.Done():
		t.Fatalf("expected no early termination, got ctx.Done(): %v", ctx.Err())
	}
	err = ctx.Err()
	if err != context.Canceled {
		t.Fatalf("expected error context.Canceled, got: %v", err)
	}
}

func TestPerformRequestWithTimeout(t *testing.T) {
	tr := &sleepingTransport{timeout: 3 * time.Second}
	httpClient := &http.Client{Transport: tr}

	client, err := NewSimpleClient(SetHttpClient(httpClient), SetMaxRetries(0), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}

	type result struct {
		res *Response
		err error
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resc := make(chan result, 1)
	go func() {
		res, err := client.PerformRequest(ctx, PerformRequestOptions{
			Method: "GET",
			Path:   "/",
		})
		resc <- result{res: res, err: err}
	}()
	select {
	case res := <-resc:
		t.Fatalf("expected timeout before response, got %v", res)
	case <-ctx.Done():
		err := ctx.Err()
		if err != context.DeadlineExceeded {
			t.Fatalf("expected error context.DeadlineExceeded, got: %v", err)
		}
	}
}

func TestPerformRequestWithCustomHTTPHeadersOnRequest(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/_tasks",
		Params: url.Values{
			"pretty": []string{"true"},
		},
		Headers: http.Header{
			"X-Opaque-Id": []string{"123456"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
	if want, have := "123456", res.Header.Get("X-Opaque-Id"); want != have {
		t.Fatalf("want response header X-Opaque-Id=%q, have %q", want, have)
	}
}

func TestPerformRequestWithCustomHTTPHeadersOnClient(t *testing.T) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client, err := NewClient(SetHeaders(http.Header{
		"Custom-Id": []string{"olivere"},
	}), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"), SetTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/_tasks",
		Params: url.Values{
			"pretty": []string{"true"},
		},
		Headers: http.Header{
			"X-Opaque-Id": []string{"123456"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
	// Request-level headers have preference
	if want, have := "123456", res.Header.Get("X-Opaque-Id"); want != have {
		t.Fatalf("want X-Opaque-Id=%q, have %q", want, have)
	}
}

func TestPerformRequestSetsDefaultUserAgent(t *testing.T) {
	var req *http.Request
	h := func(r *http.Request) (*http.Response, error) {
		req = new(http.Request)
		*req = *r
		return &http.Response{Request: r, StatusCode: http.StatusOK, Body: http.NoBody}, nil
	}
	tr := &failingTransport{path: "/", fail: h}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(SetHttpClient(httpClient), SetSniff(false), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
		Params: url.Values{
			"pretty": []string{"true"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
	// Have a default for User-Agent
	if want, have := "opensearch/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")", req.Header.Get("User-Agent"); want != have {
		t.Fatalf("want User-Agent=%q, have %q", want, have)
	}
}

func TestPerformRequestWithCustomHTTPHeadersPriority(t *testing.T) {
	var req *http.Request
	h := func(r *http.Request) (*http.Response, error) {
		req = new(http.Request)
		*req = *r
		return &http.Response{Request: r, StatusCode: http.StatusOK, Body: http.NoBody}, nil
	}
	tr := &failingTransport{path: "/", fail: h}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(SetHttpClient(httpClient), SetHeaders(http.Header{
		"Custom-Id":   []string{"olivere"},
		"User-Agent":  []string{"My user agent"},
		"X-Opaque-Id": []string{"sandra"}, // <- will be overridden by request-level header
	}), SetSniff(false), SetHealthcheck(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
		Params: url.Values{
			"pretty": []string{"true"},
		},
		Headers: http.Header{
			"X-Opaque-Id": []string{"123456"}, // <- request-level has preference
			"X-Somewhat":  []string{"somewhat"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
	if req == nil {
		t.Fatal("expected to record HTTP request")
	}
	if want, have := "123456", req.Header.Get("X-Opaque-Id"); want != have {
		t.Fatalf("want X-Opaque-Id=%q, have %q", want, have)
	}
	if want, have := "olivere", req.Header.Get("Custom-Id"); want != have {
		t.Fatalf("want Custom-Id=%q, have %q", want, have)
	}
	if want, have := "My user agent", req.Header.Get("User-Agent"); want != have {
		t.Fatalf("want User-Agent=%q, have %q", want, have)
	}
	if want, have := "somewhat", req.Header.Get("X-Somewhat"); want != have {
		t.Fatalf("want X-Somewhat=%q, have %q", want, have)
	}
}

// -- Compression --

// Notice that the trace log does always print "Accept-Encoding: gzip"
// regardless of whether compression is enabled or not. This is because
// of the underlying "httputil.DumpRequestOut".
//
// Use a real HTTP proxy/recorder to convince yourself that
// "Accept-Encoding: gzip" is NOT sent when DisableCompression
// is set to true.
//
// See also:
// https://groups.google.com/forum/#!topic/golang-nuts/ms8QNCzew8Q

func TestPerformRequestWithCompressionEnabled(t *testing.T) {
	testPerformRequestWithCompression(t, &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
}

func TestPerformRequestWithCompressionDisabled(t *testing.T) {
	testPerformRequestWithCompression(t, &http.Client{
		Transport: &http.Transport{
			DisableCompression: false,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
}

func testPerformRequestWithCompression(t *testing.T, hc *http.Client) {
	client, err := NewClient(SetHttpClient(hc), SetSniff(false), SetBasicAuth("admin", "vLPeJYa8.3RqtZCcAK6jNz"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.PerformRequest(context.TODO(), PerformRequestOptions{
		Method: "GET",
		Path:   "/",
	})
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
	if ret.ClusterName == "" {
		t.Errorf("expected cluster name; got: %q", ret.ClusterName)
	}
}
