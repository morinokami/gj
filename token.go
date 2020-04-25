package gj

type tokenType string

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

type token struct {
	Type    tokenType
	Literal string
}

var keywords = map[string]tokenType{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

// lookupKeyword returns the token if kw is in keywords, ILLEGAL otherwise.
func lookupKeyword(kw string) tokenType {
	if tok, ok := keywords[kw]; ok {
		return tok
	}
	return ILLEGAL
}
