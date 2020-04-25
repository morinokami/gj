package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/morinokami/gj/token"
)

type Expression interface {
	TokenLiteral() string
	String() string
}

type JSON struct {
	Value Expression
}

func (json *JSON) TokenLiteral() string {
	return json.Value.TokenLiteral()
}

func (json *JSON) String() string {
	return json.Value.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type Null struct {
	Token token.Token
	Value interface{}
}

func (n *Null) TokenLiteral() string { return n.Token.Literal }
func (n *Null) String() string       { return n.Token.Literal }

type Integer struct {
	Token token.Token
	Value int64
}

func (i *Integer) TokenLiteral() string { return i.Token.Literal }
func (i *Integer) String() string       { return i.Token.Literal }

type Float struct {
	Token token.Token
	Value string
}

func (f *Float) TokenLiteral() string { return f.Token.Literal }
func (f *Float) String() string       { return f.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())

	return out.String()
}

type String struct {
	Token token.Token
	Value string
}

func (s *String) TokenLiteral() string { return s.Token.Literal }
func (s *String) String() string       { return fmt.Sprintf(`"%s"`, s.Token.Literal) }

type Object struct {
	Token token.Token
	Pairs map[String]Expression
}

func (o *Object) TokenLiteral() string { return o.Token.Literal }
func (o *Object) String() string {
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

type Array struct {
	Token  token.Token
	Values []Expression
}

func (a *Array) TokenLiteral() string { return a.Token.Literal }
func (a *Array) String() string {
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
