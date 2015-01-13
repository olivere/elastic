package elastic

import (
	"encoding/json"
	"testing"
)

func TestQueryStringQuery(t *testing.T) {
	q := NewQueryStringQuery(`this AND that OR thus`)
	q = q.DefaultField("content")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"query_string":{"default_field":"content","query":"this AND that OR thus"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestQueryStringQueryEscaping(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    ``,
			Expected: ``,
		},
		{
			Input:    `no escaping`,
			Expected: `no escaping`,
		},
		{
			Input:    `+ - && || ! ( ) { } [ ] ^ " ~ * ? : \ /`,
			Expected: `\+ \- \&\& \|\| \! \( \) \{ \} \[ \] \^ \" \~ \* \? \: \\ \/`,
		},
		{
			Input:    `What~3 say~3 about~3 me…~3 I~3 don't~3 now:)~3`,
			Expected: `What\~3 say\~3 about\~3 me…\~3 I\~3 don't\~3 now\:\)\~3`,
		},
	}

	for _, test := range tests {
		got := QueryStringEscape(test.Input)
		if got != test.Expected {
			t.Errorf("expected %q on input %q; got: %q", test.Expected, test.Input, got)
		}
	}
}
