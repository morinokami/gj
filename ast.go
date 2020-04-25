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

type jsonValue struct {
	Value expression
}

func (json *jsonValue) TokenLiteral() string {
	return json.Value.TokenLiteral()
}

func (json *jsonValue) String() string {
	return json.Value.String()
}

type boolean struct {
	Token token
	Value bool
}

func (b *boolean) TokenLiteral() string { return b.Token.Literal }
func (b *boolean) String() string       { return b.Token.Literal }

type null struct {
	Token token
	Value interface{}
}

func (n *null) TokenLiteral() string { return n.Token.Literal }
func (n *null) String() string       { return n.Token.Literal }

type integer struct {
	Token token
	Value int64
}

func (i *integer) TokenLiteral() string { return i.Token.Literal }
func (i *integer) String() string       { return i.Token.Literal }

type float struct {
	Token token
	Value string
}

func (f *float) TokenLiteral() string { return f.Token.Literal }
func (f *float) String() string       { return f.Token.Literal }

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

type stringLiteral struct {
	Token token
	Value string
}

func (s *stringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *stringLiteral) String() string       { return fmt.Sprintf(`"%s"`, s.Token.Literal) }

type object struct {
	Token token
	Pairs map[stringLiteral]expression
}

func (o *object) TokenLiteral() string { return o.Token.Literal }
func (o *object) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range o.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type array struct {
	Token  token
	Values []expression
}

func (a *array) TokenLiteral() string { return a.Token.Literal }
func (a *array) String() string {
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
