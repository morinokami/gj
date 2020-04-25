package gj

import (
	"testing"
)

func TestString(t *testing.T) {
	json := &JSON{
		Value: &Object{
			Token: Token{
				Type:    LBRACE,
				Literal: "{",
			},
			Pairs: map[String]Expression{
				String{
					Token: Token{Type: STRING, Literal: "foo"},
					Value: "foo",
				}: &Integer{
					Token: Token{Type: INT, Literal: "123"},
					Value: 123,
				},
			},
		},
	}

	if json.String() != `{"foo": 123}` {
		t.Errorf("json.String() wrong. got=%q.", json.String())
	}
}
