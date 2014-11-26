package elastic

import (
	"net/http"
	"testing"
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
