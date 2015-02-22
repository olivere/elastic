// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"bytes"
	"encoding/json"
	"sync/atomic"
	"testing"
)

type decoder struct {
	dec json.Decoder

	N int64
}

func (d *decoder) Decode(data []byte, v interface{}) error {
	atomic.AddInt64(&d.N, 1)
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}

func TestDecoder(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	dec := &decoder{}
	client.SetDecoder(dec)

	tweet := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}

	// Add a document
	indexResult, err := client.Index().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		BodyJson(&tweet).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if indexResult == nil {
		t.Errorf("expected result to be != nil; got: %v", indexResult)
	}

	if dec.N != 1 {
		t.Errorf("expected %d call of decoder; got: %d", 1, dec.N)
	}
}
