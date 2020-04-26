package gj

import (
	"fmt"
	"strconv"
)

type parseFn func() expression

type parser struct {
	l      *lexer
	errors []string

	curToken  token
	peekToken token

	parseFns map[tokenType]parseFn
}

// newParser initializes a parser object and returns it.
func newParser(l *lexer) *parser {
	p := &parser{
		l:      l,
		errors: []string{},
	}

	p.parseFns = make(map[tokenType]parseFn)
	p.registerParseFn(tokTrue, p.parseBoolean)
	p.registerParseFn(tokFalse, p.parseBoolean)
	p.registerParseFn(tokNull, p.parseNull)
	p.registerParseFn(tokInt, p.parseInteger)
	p.registerParseFn(tokFloat, p.parseFloat)
	p.registerParseFn(tokMinus, p.parsePrefixExpression)
	p.registerParseFn(tokString, p.parseString)
	p.registerParseFn(tokLBrace, p.parseObject)
	p.registerParseFn(tokLBracket, p.parseArray)

	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the tokens.
func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.nextToken()
}

// curTokenIs returns true if the type of curToken is t, false otherwise.
func (p *parser) curTokenIs(t tokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs returns true if the type of peekToken is t, false otherwise.
func (p *parser) peekTokenIs(t tokenType) bool {
	return p.peekToken.Type == t
}

// If the type of peekToken is t, expectPeek returns true and advance the tokens.
// Otherwise it returns false and append an error message to errors.
func (p *parser) expectPeek(t tokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// getErrors returns the slice of error messages.
func (p *parser) getErrors() []string {
	return p.errors
}

func (p *parser) noParseFnError(t tokenType) {
	msg := fmt.Sprintf("no parse function for %s found.", t)
	p.errors = append(p.errors, msg)
}

func (p *parser) peekError(t tokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parse parses the input string and returns the result as an ast.jsonExpression.
func (p *parser) parse() *jsonExpression {
	json := &jsonExpression{}

	json.Value = p.parseExpression()

	return json
}

func (p *parser) parseExpression() expression {
	prefix := p.parseFns[p.curToken.Type]
	if prefix == nil {
		p.noParseFnError(p.curToken.Type)
		return nil
	}

	exp := prefix()

	return exp
}

func (p *parser) parseBoolean() expression {
	return &booleanExpression{Token: p.curToken, Value: p.curTokenIs(tokTrue)}
}

func (p *parser) parseNull() expression {
	return &nullExpression{Token: p.curToken, Value: nil}
}

func (p *parser) parseInteger() expression {
	i := &integerExpression{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer.", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	i.Value = value

	return i
}

func (p *parser) parseFloat() expression {
	return &floatExpression{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parsePrefixExpression() expression {
	exp := &prefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression()

	return exp
}

func (p *parser) parseString() expression {
	return &stringExpression{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parseObject() expression {
	object := &objectExpression{Token: p.curToken}
	object.Pairs = make(map[string]expression)

	for !p.peekTokenIs(tokRBrace) {
		if !p.expectPeek(tokString) {
			return nil
		}

		key := p.parseString().(*stringExpression)

		if !p.expectPeek(tokColon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression()

		object.Pairs[key.Value] = value

		if !p.peekTokenIs(tokRBrace) && !p.expectPeek(tokComma) {
			return nil
		}
	}

	if !p.expectPeek(tokRBrace) {
		return nil
	}

	return object
}

func (p *parser) parseArray() expression {
	array := &arrayExpression{Token: p.curToken}
	array.Values = []expression{}

	if p.peekTokenIs(tokRBracket) {
		p.nextToken()
		return array
	}

	p.nextToken()
	array.Values = append(array.Values, p.parseExpression())

	for p.peekTokenIs(tokComma) {
		p.nextToken()
		p.nextToken()
		array.Values = append(array.Values, p.parseExpression())
	}

	if !p.expectPeek(tokRBracket) {
		return nil
	}

	return array
}

// registerParseFn registers functions to parse each token.
func (p *parser) registerParseFn(tokenType tokenType, fn parseFn) {
	p.parseFns[tokenType] = fn
}
