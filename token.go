package gj

type tokenType string

const (
	tokIllegal = "ILLEGAL"
	tokEOF     = "EOF"

	// Delimiters
	tokComma = "COMMA"
	tokColon = "COLON"

	// Brackets
	tokLBrace   = "{"
	tokRBrace   = "}"
	tokLBracket = "["
	tokRBracket = "]"

	// Values
	tokTrue  = "TRUE"
	tokFalse = "FALSE"
	tokNull  = "NULL"

	// Key & Literals
	tokInt    = "INT"
	tokFloat  = "FLOAT"
	tokString = "STRING"

	// Operators
	tokMinus = "-"
)

// TODO: Add position
type token struct {
	Type    tokenType
	Literal string
}

var keywords = map[string]tokenType{
	"true":  tokTrue,
	"false": tokFalse,
	"null":  tokNull,
}

// lookupKeyword returns the token if kw is in keywords, tokIllegal otherwise.
func lookupKeyword(kw string) tokenType {
	if tok, ok := keywords[kw]; ok {
		return tok
	}
	return tokIllegal
}
