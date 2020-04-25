package gj

import (
	"testing"
)

func TestString(t *testing.T) {
	json := &jsonValue{
		Value: &object{
			Token: token{
				Type:    tokLBrace,
				Literal: "{",
			},
			Pairs: map[stringLiteral]expression{
				stringLiteral{
					Token: token{Type: tokString, Literal: "foo"},
					Value: "foo",
				}: &integer{
					Token: token{Type: tokInt, Literal: "123"},
					Value: 123,
				},
			},
		},
	}

	if json.String() != `{"foo": 123}` {
		t.Errorf("json.String() wrong. got=%q.", json.String())
	}
}
