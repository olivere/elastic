package elastic

import (
	"net/http"
	"testing"
)

func TestClientSingleConnection(t *testing.T) {
	client, err := NewClient(http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.pool.conns) != 1 {
		t.Fatalf("expected a pool of 1 connection, got: %v", len(client.pool.conns))
	}
	if client.pool.conns[0].url != defaultUrl {
		t.Errorf("expected default client connection url of %s, got: %s", defaultUrl, client.pool.conns[0].url)
	}
}

func TestClientMultipleConnections(t *testing.T) {
	client, err := NewClient(http.DefaultClient, "http://localhost:9200", "http://localhost:9201")
	if err != nil {
		t.Fatal(err)
	}
	if len(client.pool.conns) != 2 {
		t.Fatalf("expected a pool of 2 connections, got: %v", len(client.pool.conns))
	}
	if client.pool.conns[0].url != "http://localhost:9200" {
		t.Errorf("expected 1st connection url of %s, got: %s", "http://localhost:9200", client.pool.conns[0].url)
	}
	broken := client.pool.conns[0].IsBroken()
	if broken {
		t.Errorf("expected 1st connection url of %s to not be broken, got: %s", client.pool.conns[0].url, broken)
	}
	if client.pool.conns[1].url != "http://localhost:9201" {
		t.Errorf("expected 2nd connection url of %s, got: %s", "http://localhost:9201", client.pool.conns[1].url)
	}
	broken = client.pool.conns[1].IsBroken()
	if !broken {
		t.Errorf("expected 2nd connection url of %s to be broken, got: %s", client.pool.conns[1].url, broken)
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
