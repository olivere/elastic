package elastic

import (
	"net/http"
	"testing"
)

func TestConnection(t *testing.T) {
	conn := NewConnection(http.DefaultClient, "http://localhost:9200")
	if conn == nil {
		t.Fatalf("expected connection, got: %v", conn)
	}
	if conn.url != "http://localhost:9200" {
		t.Errorf("expected url of %s, got: %s", "http://localhost:9200", conn.url)
	}
	broken := conn.IsBroken()
	if broken {
		t.Error("expected connection to not be broken")
	}
}

func TestConnectionBroken(t *testing.T) {
	conn := NewConnection(http.DefaultClient, "http://localhost:19200")
	if conn == nil {
		t.Fatalf("expected connection, got: %v", conn)
	}
	if conn.url != "http://localhost:19200" {
		t.Errorf("expected url of %s, got: %s", "http://localhost:19200", conn.url)
	}
	broken := conn.IsBroken()
	if !broken {
		t.Error("expected broken connection")
	}
}
