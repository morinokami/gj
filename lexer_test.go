package gj

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `{
    "name": "John",
    "profile": {
        "age": 100,
        "dog": true,
        "cat": false,
        "items": [
            null,
            -123,
            "This is an \"escape\" of a double-quote",
            3.1415
        ]
    }
}`

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{tokLBrace, "{"},
		{tokString, "name"},
		{tokColon, ":"},
		{tokString, "John"},
		{tokComma, ","},
		{tokString, "profile"},
		{tokColon, ":"},
		{tokLBrace, "{"},
		{tokString, "age"},
		{tokColon, ":"},
		{tokInt, "100"},
		{tokComma, ","},
		{tokString, "dog"},
		{tokColon, ":"},
		{tokTrue, "true"},
		{tokComma, ","},
		{tokString, "cat"},
		{tokColon, ":"},
		{tokFalse, "false"},
		{tokComma, ","},
		{tokString, "items"},
		{tokColon, ":"},
		{tokLBracket, "["},
		{tokNull, "null"},
		{tokComma, ","},
		{tokMinus, "-"},
		{tokInt, "123"},
		{tokComma, ","},
		{tokString, "This is an \"escape\" of a double-quote"},
		{tokComma, ","},
		{tokFloat, "3.1415"},
		{tokRBracket, "]"},
		{tokRBrace, "}"},
		{tokRBrace, "}"},
		{tokEOF, ""},
	}

	l := newLexer(input)

	for i, tt := range tests {
		tok := l.nextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type worng. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got%q.", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
