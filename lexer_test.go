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
		{LBRACE, "{"},
		{STRING, "name"},
		{COLON, ":"},
		{STRING, "John"},
		{COMMA, ","},
		{STRING, "profile"},
		{COLON, ":"},
		{LBRACE, "{"},
		{STRING, "age"},
		{COLON, ":"},
		{INT, "100"},
		{COMMA, ","},
		{STRING, "dog"},
		{COLON, ":"},
		{TRUE, "true"},
		{COMMA, ","},
		{STRING, "cat"},
		{COLON, ":"},
		{FALSE, "false"},
		{COMMA, ","},
		{STRING, "items"},
		{COLON, ":"},
		{LBRACKET, "["},
		{NULL, "null"},
		{COMMA, ","},
		{MINUS, "-"},
		{INT, "123"},
		{COMMA, ","},
		{STRING, "This is an \"escape\" of a double-quote"},
		{COMMA, ","},
		{FLOAT, "3.1415"},
		{RBRACKET, "]"},
		{RBRACE, "}"},
		{RBRACE, "}"},
		{EOF, ""},
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
