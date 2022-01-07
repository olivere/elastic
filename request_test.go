// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"strings"
	"testing"
)

var testReq *Request // used as a temporary variable to avoid compiler optimizations in tests/benchmarks

func TestRequestSetContentType(t *testing.T) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "application/json", req.Header.Get("Content-Type"); want != have {
		t.Fatalf("want Content-Type=%q, have %q", want, have)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	if want, have := "application/x-ndjson", req.Header.Get("Content-Type"); want != have {
		t.Fatalf("want Content-Type=%q, have %q", want, have)
	}
	if want, have := "", req.Header.Get("User-Agent"); want != have {
		t.Fatalf("want User-Agent=%q, have %q", want, have)
	}
}

func TestRequestContentLength(t *testing.T) {
	tests := []struct {
		Body                interface{}
		Compress            bool
		ExpectContentLength bool
		ContentLength       int64
	}{
		{
			Body:                `{}`,
			Compress:            false,
			ExpectContentLength: true,
			ContentLength:       2,
		},
		{
			Body:                strings.NewReader(`{}`),
			Compress:            false,
			ExpectContentLength: true,
			ContentLength:       2,
		},
		{
			Body:                bytes.NewBuffer([]byte(`{}`)),
			Compress:            false,
			ExpectContentLength: true,
			ContentLength:       2,
		},
		{
			Body:                bytes.NewReader([]byte(`{}`)),
			Compress:            false,
			ExpectContentLength: true,
			ContentLength:       2,
		},
	}

	for i, tt := range tests {
		req, err := NewRequest("POST", "/")
		if err != nil {
			t.Fatal(err)
		}
		if err := req.SetBody(tt.Body, tt.Compress); err != nil {
			t.Fatal(err)
		}
		if tt.ExpectContentLength {
			contentLength := req.ContentLength
			if want, have := tt.ContentLength, contentLength; want != have {
				t.Fatalf("%d. want Content-Length=%d, have %d", i, want, have)
			}
		} else {
			if want, have := int64(0), req.ContentLength; want != have {
				t.Fatalf("%d. want Content-Length=%d, have %d", i, want, have)
			}
			hdrContentLength := req.Header.Get("Content-Length")
			if hdrContentLength != "" {
				t.Fatalf("%d. want no Content-Length, have %q", i, hdrContentLength)
			}
		}
	}
}

func BenchmarkRequestSetBodyString(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := `{"query":{"match_all":{}}}`
		err = req.SetBody(body, false)
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
		err = req.SetBody(body, true)
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
		err = req.SetBody(body, false)
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
		err = req.SetBody(body, true)
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
		err = req.SetBody(body, false)
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
		err = req.SetBody(body, true)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyStringReader(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := strings.NewReader(`{"query":{"match_all":{}}}`)
		err = req.SetBody(body, false)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}

func BenchmarkRequestSetBodyStringReaderGzip(b *testing.B) {
	req, err := NewRequest("GET", "/")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		body := strings.NewReader(`{"query":{"match_all":{}}}`)
		err = req.SetBody(body, true)
		if err != nil {
			b.Fatal(err)
		}
	}
	testReq = req
	b.ReportAllocs()
}
