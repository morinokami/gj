package gj

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Delimiters
	COMMA = "COMMA"
	COLON = "COLON"

	// Brackets
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Values
	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"

	// Key & Literals
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	MINUS = "-"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

// LookupKeyword returns the token if kw is in keywords, ILLEGAL otherwise.
func LookupKeyword(kw string) TokenType {
	if tok, ok := keywords[kw]; ok {
		return tok
	}
	return ILLEGAL
}
