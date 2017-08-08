// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"testing"
)

const (
	testKey   = "key"
	testValue = "value"
)

func TestAddHeadersToNil(t *testing.T) {
	newHeaders := addHeader(nil, testKey, testValue)

	l := len(newHeaders)
	if l != 1 {
		t.Errorf("expected one item in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 1 || actualVal[0] != testValue {
		t.Errorf("expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}

func TestAddHeadersToExistingMap(t *testing.T) {
	newHeaders := addHeader(headers{"existing": {"other"}}, testKey, testValue)

	l := len(newHeaders)
	if l != 2 {
		t.Errorf("expected two items in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 1 || actualVal[0] != testValue {
		t.Errorf("expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}

func TestAddHeadersDuplicatedKey(t *testing.T) {
	newHeaders := addHeader(headers{testKey: {"other"}}, testKey, testValue)

	l := len(newHeaders)
	if l != 1 {
		t.Errorf("expected one item in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 2 || actualVal[1] != testValue {
		t.Errorf("expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}
