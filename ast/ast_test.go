package ast

import (
	"testing"

	"github.com/morinokami/gj/token"
)

func TestString(t *testing.T) {
	json := &JSON{
		Value: &Object{
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
			},
			Pairs: map[String]Expression{
				String{
					Token: token.Token{Type: token.STRING, Literal: "foo"},
					Value: "foo",
				}: &Integer{
					Token: token.Token{Type: token.INT, Literal: "123"},
					Value: 123,
				},
			},
		},
	}

	if json.String() != `{"foo": 123}` {
		t.Errorf("json.String() wrong. got=%q.", json.String())
	}
}
