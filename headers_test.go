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
		t.Errorf("Expected one item in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("Expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 1 || actualVal[0] != testValue {
		t.Errorf("Expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}

func TestAddHeadersToExistingMap(t *testing.T) {
	newHeaders := addHeader(map[string][]string{"existing": {"other"}}, testKey, testValue)

	l := len(newHeaders)
	if l != 2 {
		t.Errorf("Expected two items in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("Expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 1 || actualVal[0] != testValue {
		t.Errorf("Expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}

func TestAddHeadersDuplicatedKey(t *testing.T) {
	newHeaders := addHeader(map[string][]string{testKey: {"other"}}, testKey, testValue)

	l := len(newHeaders)
	if l != 1 {
		t.Errorf("Expected one item in the headers map, found %d", l)
	}

	actualVal, ok := newHeaders["key"]
	if !ok {
		t.Errorf("Expected key: %s in the map, but not found", testKey)
	}
	if len(actualVal) != 2 || actualVal[1] != testValue {
		t.Errorf("Expected value: %s associated with the key: %s, but found: %s", testValue, testKey, actualVal)
	}
}
