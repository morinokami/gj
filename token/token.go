package token

type Type string

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
	STRING = "STRING"
)

type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

// LookupKeyword returns the token if kw is in keywords, ILLEGAL otherwise.
func LookupKeyword(kw string) Type {
	if tok, ok := keywords[kw]; ok {
		return tok
	}
	return ILLEGAL
}
