package lexer

import (
	"testing"

	"github.com/morinokami/go-json-parser/token"
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
          -123
      ]
   }
}`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.STRING, "name"},
		{token.COLON, ":"},
		{token.STRING, "John"},
		{token.COMMA, ","},
		{token.STRING, "profile"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "age"},
		{token.COLON, ":"},
		{token.INT, "100"},
		{token.COMMA, ","},
		{token.STRING, "dog"},
		{token.COLON, ":"},
		{token.TRUE, "true"},
		{token.COMMA, ","},
		{token.STRING, "cat"},
		{token.COLON, ":"},
		{token.FALSE, "false"},
		{token.COMMA, ","},
		{token.STRING, "items"},
		{token.COLON, ":"},
		{token.LBRACKET, "["},
		{token.NULL, "null"},
		{token.COMMA, ","},
		{token.MINUS, "-"},
		{token.INT, "123"},
		{token.RBRACKET, "]"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type worng. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got%q.", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
