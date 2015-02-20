package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSingleUrl(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 1 {
		t.Fatalf("expected 1 default client url, got: %v", client.urls)
	}
	if client.urls[0] != defaultUrl {
		t.Errorf("expected default client url of %s, got: %s", defaultUrl, client.urls[0])
	}
}

func TestMultipleUrls(t *testing.T) {
	client, err := NewClient(http.DefaultClient, "http://localhost:9200", "http://localhost:9201")
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 2 {
		t.Fatalf("expected 2 default client urls, got: %v", client.urls)
	}
	if client.urls[0] != "http://localhost:9200" {
		t.Errorf("expected 1st client url of %s, got: %s", "http://localhost:9200", client.urls[0])
	}
	if client.urls[1] != "http://localhost:9201" {
		t.Errorf("expected 2nd client url of %s, got: %s", "http://localhost:9201", client.urls[0])
	}
}

func TestFindingActiveClient(t *testing.T) {
	client, err := NewClient(http.DefaultClient, "http://localhost:19200", "http://localhost:9200")
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 2 {
		t.Fatalf("expected 2 default client urls, got: %v", client.urls)
	}
	if !client.hasActive {
		t.Errorf("expected to have active connection, got: %v", client.hasActive)
	}
	expected := "http://localhost:9200"
	if client.activeUrl != expected {
		t.Errorf("expected active url to be %s, got: %v", expected, client.activeUrl)
	}
}

func TestFindingNoActiveClient(t *testing.T) {
	client, err := NewClient(http.DefaultClient, "http://localhost:19200", "http://localhost:19201")
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 2 {
		t.Fatalf("expected 2 default client urls, got: %v", client.urls)
	}
	if client.hasActive {
		t.Errorf("expected to not have an active connection, got: %v", client.hasActive)
	}
	if client.activeUrl != "" {
		t.Errorf("expected no active url, got: %v", client.activeUrl)
	}
	req, err := client.NewRequest("HEAD", "/")
	if err != ErrNoClient {
		t.Errorf("expected ErrNoClient, got: %v", err)
	}
	if req != nil {
		t.Errorf("expected no request, got: %v", req)
	}
}

func TestElasticsearchVersion(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	version, err := client.ElasticsearchVersion(defaultUrl)
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Errorf("expected a version number, got: %q", version)
	}
}

func TestPerformRequest(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
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
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	var w bytes.Buffer
	out := log.New(&w, "LOGGER ", log.LstdFlags)
	client.SetLogger(out)

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
	pattern := `^LOGGER \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} GET ` + defaultUrl + `/ \[status:200, request:\d+\.\d{3}s\]\n`
	matched, err := regexp.MatchString(pattern, got)
	if err != nil {
		t.Fatalf("expected log line to match %q; got: %v", pattern, err)
	}
	if !matched {
		t.Errorf("expected log line to match %q", pattern)
	}
}

func TestPerformRequestWithLoggerAndTracer(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	var lw bytes.Buffer
	lout := log.New(&lw, "LOGGER ", log.LstdFlags)
	client.SetLogger(lout)

	var tw bytes.Buffer
	tout := log.New(&tw, "TRACER ", log.LstdFlags)
	client.SetTracer(tout)

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
		t.Error("expected logger output; got: %q", lgot)
	}

	tgot := tw.String()
	if tgot == "" {
		t.Error("expected tracer output; got: %q", tgot)
	}
}

func TestSniffNode(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan sniffResult, 1)
	go func() {
		ch <- client.sniffNode(defaultUrl)
	}()

	select {
	case res := <-ch:
		if len(res.URLs) != 1 {
			t.Fatalf("expected %d node URL; got: %d", 1, len(res.URLs))
		}
		if res.URLs[0] != "http://127.0.0.1:9200" {
			t.Fatalf("expected node URL %q; got: %q", "http://127.0.0.1:9200", res.URLs[0])
		}
		break
	case <-time.After(2 * time.Second):
		t.Fatal("expected no timeout in sniff node")
		break
	}
}

func TestSniff(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan []string, 1)
	go func() {
		ch <- client.sniff()
	}()

	select {
	case urls := <-ch:
		if len(urls) != 1 {
			t.Fatalf("expected %d URL; got: %d", 1, len(urls))
		}
		if urls[0] != "http://127.0.0.1:9200" {
			t.Fatalf("expected node URL %q; got: %q", "http://127.0.0.1:9200", urls[0])
		}
		break
	case <-time.After(2 * time.Second):
		t.Fatal("expected no timeout in sniff")
		break
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

func TestPerformRequestWithMaxRetries(t *testing.T) {
	var numFailedReqs int
	fail := func(r *http.Request) (*http.Response, error) {
		numFailedReqs += 1
		return &http.Response{Request: r, StatusCode: 400}, nil
	}

	// Run against a failing endpoint and see if PerformRequest
	// retries correctly.
	tr := &failingTransport{path: "/fail", fail: fail}
	httpClient := &http.Client{Transport: tr}

	client, err := NewClient(httpClient)
	if err != nil {
		t.Fatal(err)
	}

	// Retry 5 times
	client.SetMaxRetries(5)

	res, err := client.PerformRequest("GET", "/fail", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if res != nil {
		t.Fatal("expected no response")
	}
	// Check if really tried 5 times
	if numFailedReqs != 5 {
		t.Errorf("expected %d failed requests; got: %d", 5, numFailedReqs)
	}
}
