package lexer

import "github.com/morinokami/go-json-parser/token"

const eof = 0

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// New returns a Lexer object.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns the next token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case eof:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readKeyword()
			tok.Type = token.LookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

// skipWhitespace kips whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readChar reads the next character and advances the position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = eof
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// readString returns a series of characters surrounded by double quotes.
// It advances the position until it encounters either a closing double quote or the end of the input.
func (l *Lexer) readString() string {
	start := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == eof {
			break
		}
	}

	return l.input[start:l.position]
}

// TODO: support floats
// readNumber returns an integer as a string.
// It advances the position until it encounters a non-digit character.
func (l *Lexer) readNumber() string {
	start := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

// TODO: support escaping double quotes
// readKeyword returns a string of keywords.
// It advances the position until it encounters a non-alphabetic character.
func (l *Lexer) readKeyword() string {
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
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
