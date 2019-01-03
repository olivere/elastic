// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDocvalueField(t *testing.T) {
	tests := []struct {
		Doc  DocvalueField
		Want interface{}
	}{
		{
			Doc:  DocvalueField{},
			Want: "",
		},
		{
			Doc:  DocvalueField{Field: "name"},
			Want: "name",
		},
		{
			Doc:  DocvalueField{Field: "name", Format: "epoch_millis"},
			Want: map[string]interface{}{"field": "name", "format": "epoch_millis"},
		},
	}
	for _, tt := range tests {
		have, err := tt.Doc.Source()
		if err != nil {
			t.Fatalf("Source(%#v): err=%v", tt.Doc, err)
		}
		if want := tt.Want; !cmp.Equal(want, have) {
			t.Fatalf("Source(%#v): want %v, have %v", tt.Doc, want, have)
		}
	}
}

func TestDocvalueFields(t *testing.T) {
	doc := DocvalueFields{
		DocvalueField{Field: "retweets"},
		DocvalueField{Field: "name", Format: "epoch_millis"},
	}
	have, err := doc.Source()
	if err != nil {
		t.Fatalf("Source(%#v): err=%v", doc, err)
	}
	want := []interface{}{
		"retweets",
		map[string]interface{}{"field": "name", "format": "epoch_millis"},
	}
	if !cmp.Equal(want, have) {
		t.Fatalf("Source(%#v): want %v, have %v", doc, want, have)
	}
}
