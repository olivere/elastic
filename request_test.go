// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"compress/gzip"
	"testing"
)

var testReq *Request // used as a temporary variable to avoid compiler optimizations in tests/benchmarks

func TestRequestSetContentType(t *testing.T) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "application/json", req.Header.Get("Content-Type"); want != have {
		t.Fatalf("want %q, have %q", want, have)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	if want, have := "application/x-ndjson", req.Header.Get("Content-Type"); want != have {
		t.Fatalf("want %q, have %q", want, have)
	}
}

func BenchmarkRequestSetBodyString(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := `{"query":{"match_all":{}}}`
		err = req.SetBody(body, false, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyStringGzip(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := `{"query":{"match_all":{}}}`
		err = req.SetBody(body, true, gzip.NewWriter(new(bytes.Buffer)))
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyBytes(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := []byte(`{"query":{"match_all":{}}}`)
		err = req.SetBody(body, false, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyBytesGzip(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := []byte(`{"query":{"match_all":{}}}`)
		err = req.SetBody(body, true, gzip.NewWriter(new(bytes.Buffer)))
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyMap(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
		err = req.SetBody(body, false, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyMapGzip(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
		err = req.SetBody(body, true, gzip.NewWriter(new(bytes.Buffer)))
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}
