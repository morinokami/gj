package gj

import (
	"fmt"
	"reflect"
	"testing"
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
		l := newLexer(tt.input)
		p := newParser(l)
		json := p.parse()
		checkParserErrors(t, p)

		boolean, ok := json.Value.(*boolean)
		if !ok {
			t.Fatalf("json.Value not boolean. got=%T.", json.Value)
		}
		if boolean.Value != tt.expected {
			t.Errorf("boolean.Value not %t. got=%t.", tt.expected, boolean.Value)
		}
	}
}

func TestNullExpression(t *testing.T) {
	input := "null"

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	null, ok := json.Value.(*null)
	if !ok {
		t.Fatalf("json.Value not null. got=%T.", json.Value)
	}
	if null.Value != nil {
		t.Errorf("null.Value not %v. got=%v.", nil, null.Value)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "42"

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	testNumber(t, json.Value, int64(42))
}

func TestFloatExpression(t *testing.T) {
	input := "2.7182"

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
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
		l := newLexer(tt.input)
		p := newParser(l)
		json := p.parse()
		checkParserErrors(t, p)

		exp, ok := json.Value.(*prefixExpression)
		if !ok {
			t.Fatalf("json.Value not prefixExpression. got=%T", json.Value)
		}
		if exp.Operator != tt.expectedOperator {
			t.Fatalf("exp.Operator not '-'. got=%s", exp.Operator)
		}
		testNumber(t, exp.Right, tt.expectedValue)
	}
}

func TestStringExpression(t *testing.T) {
	input := `"Hello world!"`

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	str, ok := json.Value.(*stringLiteral)
	if !ok {
		t.Fatalf("json.Value not stringLiteral. got=%T.", json.Value)
	}
	if str.Value != "Hello world!" {
		t.Errorf("str.Value not %q. got=%q.", "Hello world!", str.Value)
	}
}

func TestEmptyObjectExpression(t *testing.T) {
	input := "{}"

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	object, ok := json.Value.(*object)
	if !ok {
		t.Fatalf("json.Value not object. got=%T.", json.Value)
	}
	if len(object.Pairs) != 0 {
		t.Errorf("object.Pairs has wrong length. got=%d.", len(object.Pairs))
	}
}

func TestObjectExpression(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	object, ok := json.Value.(*object)
	if !ok {
		t.Fatalf("json.Value not object. got=%T.", json.Value)
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

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	array, ok := json.Value.(*array)
	if !ok {
		t.Fatalf("json.Value not array. got=%T.", json.Value)
	}
	if len(array.Values) != 0 {
		t.Errorf("array.Values has wrong length. got=%d.", len(array.Values))
	}
}

func TestArrayExpression(t *testing.T) {
	input := "[1, 2, 3]"

	l := newLexer(input)
	p := newParser(l)
	json := p.parse()
	checkParserErrors(t, p)

	array, ok := json.Value.(*array)
	if !ok {
		t.Fatalf("json.Value not array. got=%T.", json.Value)
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
		l := newLexer(tt.input)
		p := newParser(l)
		p.parse()
		errors := p.getErrors()

		if len(errors) == 0 {
			t.Fatalf("expected to have errors, got none.")
		}
		if !reflect.DeepEqual(errors, tt.expected) {
			t.Errorf("%s - got=%q, want=%q.", tt.desc, errors, tt.expected)
		}
	}
}

func testNumber(t *testing.T, exp expression, value interface{}) {
	switch number := value.(type) {
	case int64:
		testInteger(t, exp, number)
	case string: // for float
		testFloat(t, exp, number)
	default:
		t.Errorf("Wrong value type - %T.", value)
	}
}

func testInteger(t *testing.T, exp expression, value int64) {
	i, ok := exp.(*integer)
	if !ok {
		t.Fatalf("exp not integer. got=%T.", exp)
	}
	if i.Value != value {
		t.Errorf("i.Value not %d. got=%d.", value, i.Value)
	}
	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral() not %d. got=%s.", value, i.TokenLiteral())
	}
}

func testFloat(t *testing.T, exp expression, value string) {
	f, ok := exp.(*float)
	if !ok {
		t.Fatalf("exp not float. got=%T.", exp)
	}
	if f.Value != value {
		t.Errorf("f.Value not %s. got=%s.", value, f.Value)
	}
	if f.TokenLiteral() != value {
		t.Errorf("f.TokenLiteral() not %s. got=%s.", value, f.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *parser) {
	errors := p.getErrors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
