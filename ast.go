package gj

import (
	"bytes"
	"fmt"
	"strings"
)

type expression interface {
	TokenLiteral() string
	String() string
}

type jsonExpression struct {
	Value expression
}

func (json *jsonExpression) TokenLiteral() string {
	return json.Value.TokenLiteral()
}

func (json *jsonExpression) String() string {
	return json.Value.String()
}

type booleanExpression struct {
	Token token
	Value bool
}

func (b *booleanExpression) TokenLiteral() string { return b.Token.Literal }
func (b *booleanExpression) String() string       { return b.Token.Literal }

type nullExpression struct {
	Token token
	Value interface{}
}

func (n *nullExpression) TokenLiteral() string { return n.Token.Literal }
func (n *nullExpression) String() string       { return n.Token.Literal }

type integerExpression struct {
	Token token
	Value int64
}

func (i *integerExpression) TokenLiteral() string { return i.Token.Literal }
func (i *integerExpression) String() string       { return i.Token.Literal }

type floatExpression struct {
	Token token
	Value string
}

func (f *floatExpression) TokenLiteral() string { return f.Token.Literal }
func (f *floatExpression) String() string       { return f.Token.Literal }

type prefixExpression struct {
	Token    token
	Operator string
	Right    expression
}

func (pe *prefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *prefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())

	return out.String()
}

type stringExpression struct {
	Token token
	Value string
}

func (s *stringExpression) TokenLiteral() string { return s.Token.Literal }
func (s *stringExpression) String() string       { return fmt.Sprintf(`"%s"`, s.Token.Literal) }

type objectExpression struct {
	Token token
	Pairs map[string]expression
}

func (o *objectExpression) TokenLiteral() string { return o.Token.Literal }
func (o *objectExpression) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range o.Pairs {
		pairs = append(pairs, `"`+key+`": `+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type arrayExpression struct {
	Token  token
	Values []expression
}

func (a *arrayExpression) TokenLiteral() string { return a.Token.Literal }
func (a *arrayExpression) String() string {
	var out bytes.Buffer

	var values []string
	for _, v := range a.Values {
		values = append(values, v.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")

	return out.String()
}
