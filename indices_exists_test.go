// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIndicesExistsWithoutIndex(t *testing.T) {
	client := setupTestClient(t)

	// No index name -> fail with error
	res, err := NewIndicesExistsService(client).Do(context.TODO())
	if err == nil {
		t.Fatalf("expected IndicesExists to fail without index name")
	}
	if res != false {
		t.Fatalf("expected result to be false; got: %v", res)
	}
}
func TestIndicesExistsService_buildURL(t *testing.T) {
	tests := []struct {
		index             []string
		pretty            bool
		local             bool
		ignoreUnavailable bool
		allowNoIndices    bool
		expandWildcards   string
		includeTypeName   bool
		expectedParams    url.Values
	}{
		{
			pretty:            true,
			local:             false,
			ignoreUnavailable: true,
			allowNoIndices:    false,
			includeTypeName:   true,
			expectedParams: url.Values{
				"pretty":             []string{"true"},
				"local":              []string{"false"},
				"ignore_unavailable": []string{"true"},
				"allow_no_indices":   []string{"false"},
				"include_type_name":  []string{"true"},
			},
		},
	}

	for _, tt := range tests {
		_, params, err := NewIndicesExistsService(nil).
			Pretty(tt.pretty).
			Local(tt.local).
			IgnoreUnavailable(tt.ignoreUnavailable).
			AllowNoIndices(tt.allowNoIndices).
			IncludeTypeName(tt.includeTypeName).
			buildURL()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if want, have := tt.expectedParams, params; !cmp.Equal(want, have) {
			t.Errorf("expected params=%#v; got: %#v\ndiff: %s", want, have, cmp.Diff(want, have))
		}
	}
}
