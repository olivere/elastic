package elastic

import (
	"net/http"
	"testing"
)

func TestConnectionPoolRoundRobin(t *testing.T) {
	client, err := NewClient(http.DefaultClient,
		"http://localhost:19200", // broken
		"http://localhost:9200",  // ok
		"http://localhost:19201", // broken
		"http://127.0.0.1:9200",  // ok
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.pool.conns) != 4 {
		t.Fatalf("expected 4 connections, got: %v", len(client.pool.conns))
	}

	// 1st request must return http://localhost:9200
	url, err := client.pool.GetNextRequestURL()
	if err != nil {
		t.Fatal(err)
	}
	if url != "http://localhost:9200" {
		t.Fatalf("expected 1st request to return %s, got: %s", "http://localhost:9200", url)
	}

	// 2nd request must return http://127.0.0.1:9200
	url, err = client.pool.GetNextRequestURL()
	if err != nil {
		t.Fatal(err)
	}
	if url != "http://127.0.0.1:9200" {
		t.Fatalf("expected 2nd request to return %s, got: %s", "http://127.0.0.1:9200", url)
	}

	// 3rd request must return http://localhost:9200 again
	url, err = client.pool.GetNextRequestURL()
	if err != nil {
		t.Fatal(err)
	}
	if url != "http://localhost:9200" {
		t.Fatalf("expected 3rd request to return %s, got: %s", "http://localhost:9200", url)
	}
}

func TestConnectionPoolWithAllBrokenConnectionReturnsErrNoClient(t *testing.T) {
	client, err := NewClient(http.DefaultClient,
		"http://localhost:19200", // broken
		"http://localhost:19201", // broken
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.pool.conns) != 2 {
		t.Fatalf("expected 2 connections, got: %v", len(client.pool.conns))
	}

	// 1st request must return ErrNoClient
	url, err := client.pool.GetNextRequestURL()
	if err != ErrNoClient {
		t.Errorf("expected ErrNoClient, got: %v", err)
	}
	if url != "" {
		t.Errorf("expected to return blank URL, got: %s", url)
	}

	// 2nd request must also return ErrNoClient
	url, err = client.pool.GetNextRequestURL()
	if err != ErrNoClient {
		t.Errorf("expected ErrNoClient, got: %v", err)
	}
	if url != "" {
		t.Errorf("expected to return blank URL, got: %s", url)
	}

	// 3rd request must, again, return ErrNoClient
	url, err = client.pool.GetNextRequestURL()
	if err != ErrNoClient {
		t.Errorf("expected to return ErrNoClient, got: %v", err)
	}
	if url != "" {
		t.Errorf("expected to return blank URL, got: %s", url)
	}
}
