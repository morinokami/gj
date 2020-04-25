package parser

import (
	"fmt"
	"strconv"

	"github.com/morinokami/gj/ast"
	"github.com/morinokami/gj/lexer"
	"github.com/morinokami/gj/token"
)

type parseFn func() ast.Expression

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	parseFns map[token.Type]parseFn
}

// New initializes a Parser object and returns it.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.parseFns = make(map[token.Type]parseFn)
	p.registerParseFn(token.TRUE, p.parseBoolean)
	p.registerParseFn(token.FALSE, p.parseBoolean)
	p.registerParseFn(token.NULL, p.parseNull)
	p.registerParseFn(token.INT, p.parseInteger)
	p.registerParseFn(token.FLOAT, p.parseFloat)
	p.registerParseFn(token.MINUS, p.parsePrefixExpression)
	p.registerParseFn(token.STRING, p.parseString)
	p.registerParseFn(token.LBRACE, p.parseObject)
	p.registerParseFn(token.LBRACKET, p.parseArray)

	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the tokens.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs returns true if the type of curToken is t, false otherwise.
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// peekTokenIs returns true if the type of peekToken is t, false otherwise.
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// If the type of peekToken is t, expectPeek returns true and advance the tokens.
// Otherwise it returns false and append an error message to errors.
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Errors returns the slice of error messages.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noParseFnError(t token.Type) {
	msg := fmt.Sprintf("no parse function for %s found.", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Parse parses the input string and returns the result as an ast.JSON.
func (p *Parser) Parse() *ast.JSON {
	json := &ast.JSON{}

	json.Value = p.parseExpression()

	return json
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.parseFns[p.curToken.Type]
	if prefix == nil {
		p.noParseFnError(p.curToken.Type)
		return nil
	}

	exp := prefix()

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{Token: p.curToken, Value: nil}
}

func (p *Parser) parseInteger() ast.Expression {
	i := &ast.Integer{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer.", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	i.Value = value

	return i
}

func (p *Parser) parseFloat() ast.Expression {
	return &ast.Float{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression()

	return exp
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseObject() ast.Expression {
	object := &ast.Object{Token: p.curToken}
	object.Pairs = make(map[ast.String]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		if !p.expectPeek(token.STRING) {
			return nil
		}

		key := p.parseString().(*ast.String)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression()

		object.Pairs[*key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return object
}

func (p *Parser) parseArray() ast.Expression {
	array := &ast.Array{Token: p.curToken}
	array.Values = []ast.Expression{}

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()
	array.Values = append(array.Values, p.parseExpression())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		array.Values = append(array.Values, p.parseExpression())
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return array
}

// registerParseFn registers functions to parse each token.
func (p *Parser) registerParseFn(tokenType token.Type, fn parseFn) {
	p.parseFns[tokenType] = fn
}
