package gj

import (
	"testing"
)

func TestString(t *testing.T) {
	json := &jsonExpression{
		Value: &objectExpression{
			Token: token{
				Type:    tokLBrace,
				Literal: "{",
			},
			Pairs: map[string]expression{
				"foo": &integerExpression{
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
