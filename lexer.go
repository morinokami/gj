package gj

import (
	"strings"
)

const eof = 0

type lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// newLexer returns a lexer object.
func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

// nextToken returns the next token.
func (l *lexer) nextToken() token {
	var tok token

	l.skipWhitespace()

	switch l.ch {
	case ',':
		tok = newToken(COMMA, l.ch)
	case ':':
		tok = newToken(COLON, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case '-':
		tok = newToken(MINUS, l.ch)
	case eof:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readKeyword()
			tok.Type = lookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = FLOAT
			} else {
				tok.Type = INT
			}
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

// skipWhitespace kips whitespace characters.
func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// TODO: Support Unicode
// readChar reads the next character and advances the position in the input string.
func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = eof
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return eof
	} else {
		return l.input[l.readPosition]
	}
}

// readString returns a series of characters surrounded by double quotes.
// It advances the position until it encounters either a closing double quote or the end of the input.
func (l *lexer) readString() string {
	start := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == eof {
			break
		}
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar()
		}
	}

	result := strings.Replace(l.input[start:l.position], `\"`, `"`, -1)

	return result
}

// readNumber returns an integer as a string.
// It advances the position until it encounters a non-digit character.
func (l *lexer) readNumber() string {
	start := l.position
	readDecimalPoint := false

	for isDigit(l.ch) || (l.ch == '.' && !readDecimalPoint) {
		if l.ch == '.' {
			readDecimalPoint = true
		}
		l.readChar()
	}

	return l.input[start:l.position]
}

// readKeyword returns a string of keywords.
// It advances the position until it encounters a non-alphabetic character.
func (l *lexer) readKeyword() string {
	start := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

// isLetter returns true if the character is either an alphabet or an underscore character, false otherwise.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns true if the character is a digit string, false otherwise.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken initializes a token and returns it.
func newToken(tokenType tokenType, ch byte) token {
	return token{Type: tokenType, Literal: string(ch)}
}
