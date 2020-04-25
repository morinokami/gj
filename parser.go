package gj

import (
	"fmt"
	"strconv"
)

type parseFn func() Expression

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	parseFns map[TokenType]parseFn
}

// New initializes a Parser object and returns it.
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.parseFns = make(map[TokenType]parseFn)
	p.registerParseFn(TRUE, p.parseBoolean)
	p.registerParseFn(FALSE, p.parseBoolean)
	p.registerParseFn(NULL, p.parseNull)
	p.registerParseFn(INT, p.parseInteger)
	p.registerParseFn(FLOAT, p.parseFloat)
	p.registerParseFn(MINUS, p.parsePrefixExpression)
	p.registerParseFn(STRING, p.parseString)
	p.registerParseFn(LBRACE, p.parseObject)
	p.registerParseFn(LBRACKET, p.parseArray)

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
func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs returns true if the type of peekToken is t, false otherwise.
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// If the type of peekToken is t, expectPeek returns true and advance the tokens.
// Otherwise it returns false and append an error message to errors.
func (p *Parser) expectPeek(t TokenType) bool {
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

func (p *Parser) noParseFnError(t TokenType) {
	msg := fmt.Sprintf("no parse function for %s found.", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Parse parses the input string and returns the result as an ast.JSON.
func (p *Parser) Parse() *JSON {
	json := &JSON{}

	json.Value = p.parseExpression()

	return json
}

func (p *Parser) parseExpression() Expression {
	prefix := p.parseFns[p.curToken.Type]
	if prefix == nil {
		p.noParseFnError(p.curToken.Type)
		return nil
	}

	exp := prefix()

	return exp
}

func (p *Parser) parseBoolean() Expression {
	return &Boolean{Token: p.curToken, Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parseNull() Expression {
	return &Null{Token: p.curToken, Value: nil}
}

func (p *Parser) parseInteger() Expression {
	i := &Integer{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer.", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	i.Value = value

	return i
}

func (p *Parser) parseFloat() Expression {
	return &Float{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	exp := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression()

	return exp
}

func (p *Parser) parseString() Expression {
	return &String{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseObject() Expression {
	object := &Object{Token: p.curToken}
	object.Pairs = make(map[String]Expression)

	for !p.peekTokenIs(RBRACE) {
		if !p.expectPeek(STRING) {
			return nil
		}

		key := p.parseString().(*String)

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression()

		object.Pairs[*key] = value

		if !p.peekTokenIs(RBRACE) && !p.expectPeek(COMMA) {
			return nil
		}
	}

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return object
}

func (p *Parser) parseArray() Expression {
	array := &Array{Token: p.curToken}
	array.Values = []Expression{}

	if p.peekTokenIs(RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()
	array.Values = append(array.Values, p.parseExpression())

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		array.Values = append(array.Values, p.parseExpression())
	}

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return array
}

// registerParseFn registers functions to parse each token.
func (p *Parser) registerParseFn(tokenType TokenType, fn parseFn) {
	p.parseFns[tokenType] = fn
}
