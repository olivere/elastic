package elastic

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestResponseError(t *testing.T) {
	raw := "HTTP/1.1 404 Not Found\r\n" +
		"\r\n" +
		`{"error":{"root_cause":[{"type":"index_missing_exception","reason":"no such index","index":"elastic-test"}],"type":"index_missing_exception","reason":"no such index","index":"elastic-test"},"status":404}` + "\r\n"
	r := bufio.NewReader(strings.NewReader(raw))

	resp, err := http.ReadResponse(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = checkResponse(resp)
	if err == nil {
		t.Fatalf("expected error; got: %v", err)
	}

	// Check for correct error message
	expected := fmt.Sprintf("elastic: Error %d (%s): no such index [type=index_missing_exception]", resp.StatusCode, http.StatusText(resp.StatusCode))
	got := err.Error()
	if got != expected {
		t.Fatalf("expected %q; got: %q", expected, got)
	}

	// Check that error is of type *elastic.Error, which contains additional information
	e, ok := err.(*Error)
	if !ok {
		t.Fatal("expected error to be of type *elastic.Error")
	}
	if e.Status != resp.StatusCode {
		t.Fatalf("expected status code %d; got: %d", resp.StatusCode, e.Status)
	}
	if e.Details == nil {
		t.Fatalf("expected error details; got: %v", e.Details)
	}
	if got, want := e.Details.Index, "elastic-test"; got != want {
		t.Fatalf("expected error details index %q; got: %q", want, got)
	}
	if got, want := e.Details.Type, "index_missing_exception"; got != want {
		t.Fatalf("expected error details type %q; got: %q", want, got)
	}
	if got, want := e.Details.Reason, "no such index"; got != want {
		t.Fatalf("expected error details reason %q; got: %q", want, got)
	}
	if got, want := len(e.Details.RootCause), 1; got != want {
		t.Fatalf("expected %d error details root causes; got: %d", want, got)
	}

	if got, want := e.Details.RootCause[0].Index, "elastic-test"; got != want {
		t.Fatalf("expected root cause index %q; got: %q", want, got)
	}
	if got, want := e.Details.RootCause[0].Type, "index_missing_exception"; got != want {
		t.Fatalf("expected root cause type %q; got: %q", want, got)
	}
	if got, want := e.Details.RootCause[0].Reason, "no such index"; got != want {
		t.Fatalf("expected root cause reason %q; got: %q", want, got)
	}
}
