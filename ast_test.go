package gj

import (
	"testing"
)

func TestString(t *testing.T) {
	json := &jsonValue{
		Value: &object{
			Token: token{
				Type:    LBRACE,
				Literal: "{",
			},
			Pairs: map[stringLiteral]expression{
				stringLiteral{
					Token: token{Type: STRING, Literal: "foo"},
					Value: "foo",
				}: &integer{
					Token: token{Type: INT, Literal: "123"},
					Value: 123,
				},
			},
		},
	}

	if json.String() != `{"foo": 123}` {
		t.Errorf("json.String() wrong. got=%q.", json.String())
	}
}
