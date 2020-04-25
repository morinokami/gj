package gj

import "testing"

func TestParser(t *testing.T) {
	input := `{
  "data": {
    "id": 888,
    "foo": "bar"
  }
}`

	json, err := ParseString(input)
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if json.String() != `{"data": {"id": 888, "foo": "bar"}}` {
		t.Errorf("Parsed result wrong. got=%s, want=%s", json.String(), `{"data": {"id": 888, "foo": "bar"}}`)
	}
}
