package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/morinokami/go-json-parser/ast"
	"github.com/morinokami/go-json-parser/lexer"
)

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		json := p.Parse()
		checkParserErrors(t, p)

		boolean, ok := json.Value.(*ast.Boolean)
		if !ok {
			t.Fatalf("json.Value not ast.Boolean. got=%T.", json.Value)
		}
		if boolean.Value != tt.expected {
			t.Errorf("boolean.Value not %t. got=%t.", tt.expected, boolean.Value)
		}
	}
}

func TestNullExpression(t *testing.T) {
	input := "null"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	null, ok := json.Value.(*ast.Null)
	if !ok {
		t.Fatalf("json.Value not ast.Null. got=%T.", json.Value)
	}
	if null.Value != nil {
		t.Errorf("null.Value not %v. got=%v.", nil, null.Value)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "42"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	testNumber(t, json.Value, int64(42))
}

func TestFloatExpression(t *testing.T) {
	input := "2.7182"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	testNumber(t, json.Value, "2.7182")
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedOperator string
		expectedValue    interface{}
	}{
		{"-273", "-", int64(273)},
		{"-3.14", "-", "3.14"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		json := p.Parse()
		checkParserErrors(t, p)

		exp, ok := json.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("json.Value not ast.PrefixExpression. got=%T", json.Value)
		}
		if exp.Operator != tt.expectedOperator {
			t.Fatalf("exp.Operator not '-'. got=%s", exp.Operator)
		}
		testNumber(t, exp.Right, tt.expectedValue)
	}
}

func TestStringExpression(t *testing.T) {
	input := `"Hello world!"`

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	str, ok := json.Value.(*ast.String)
	if !ok {
		t.Fatalf("json.Value not ast.String. got=%T.", json.Value)
	}
	if str.Value != "Hello world!" {
		t.Errorf("str.Value not %q. got=%q.", "Hello world!", str.Value)
	}
}

func TestEmptyObjectExpression(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	object, ok := json.Value.(*ast.Object)
	if !ok {
		t.Fatalf("json.Value not ast.Object. got=%T.", json.Value)
	}
	if len(object.Pairs) != 0 {
		t.Errorf("object.Pairs has wrong length. got=%d.", len(object.Pairs))
	}
}

func TestObjectExpression(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	object, ok := json.Value.(*ast.Object)
	if !ok {
		t.Fatalf("json.Value not ast.Object. got=%T.", json.Value)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	if len(object.Pairs) != 3 {
		t.Errorf("object.Pairs has wrong length. got=%d.", len(object.Pairs))
	}
	for key, value := range object.Pairs {
		expectedValue := expected[key.Value]
		testNumber(t, value, expectedValue)
	}
}

func TestEmptyArrayExpression(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	array, ok := json.Value.(*ast.Array)
	if !ok {
		t.Fatalf("json.Value not ast.Array. got=%T.", json.Value)
	}
	if len(array.Values) != 0 {
		t.Errorf("array.Values has wrong length. got=%d.", len(array.Values))
	}
}

func TestArrayExpression(t *testing.T) {
	input := "[1, 2, 3]"

	l := lexer.New(input)
	p := New(l)
	json := p.Parse()
	checkParserErrors(t, p)

	array, ok := json.Value.(*ast.Array)
	if !ok {
		t.Fatalf("json.Value not ast.Array. got=%T.", json.Value)
	}
	if len(array.Values) != 3 {
		t.Errorf("array.Values has wrong length. got=%d.", len(array.Values))
	}
	for i := 0; i < 3; i++ {
		testNumber(t, array.Values[i], int64(i+1))
	}
}

func TestIllegalExpression(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected []string
	}{
		{
			"Illegal token test",
			"!",
			[]string{"no parse function for ILLEGAL found."}},
		{
			"Unclosed object",
			"{",
			[]string{"expected next token to be STRING, got EOF instead."}},
		{
			"Keys other than string not allowed",
			"{1: 1}",
			[]string{"expected next token to be STRING, got INT instead."},
		},
		{
			"Wrong bracket usage",
			`["foo": "bar"]`,
			[]string{"expected next token to be ], got COLON instead."},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		p.Parse()
		errors := p.Errors()

		if len(errors) == 0 {
			t.Fatalf("expected to have errors, got none.")
		}
		if !reflect.DeepEqual(errors, tt.expected) {
			t.Errorf("%s - got=%q, want=%q.", tt.desc, errors, tt.expected)
		}
	}
}

func testNumber(t *testing.T, exp ast.Expression, value interface{}) {
	switch number := value.(type) {
	case int64:
		testInteger(t, exp, number)
	case string: // for Float
		testFloat(t, exp, number)
	default:
		t.Errorf("Wrong value type - %T.", value)
	}
}

func testInteger(t *testing.T, exp ast.Expression, value int64) {
	i, ok := exp.(*ast.Integer)
	if !ok {
		t.Fatalf("exp not ast.Integer. got=%T.", exp)
	}
	if i.Value != value {
		t.Errorf("i.Value not %d. got=%d.", value, i.Value)
	}
	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral() not %d. got=%s.", value, i.TokenLiteral())
	}
}

func testFloat(t *testing.T, exp ast.Expression, value string) {
	f, ok := exp.(*ast.Float)
	if !ok {
		t.Fatalf("exp not ast.Float. got=%T.", exp)
	}
	if f.Value != value {
		t.Errorf("f.Value not %s. got=%s.", value, f.Value)
	}
	if f.TokenLiteral() != value {
		t.Errorf("f.TokenLiteral() not %s. got=%s.", value, f.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
