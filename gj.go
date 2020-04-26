package gj

import (
	"fmt"
	"strings"
)

type JSON struct {
	json *jsonValue
}

func ParseString(input string) (*JSON, error) {
	l := newLexer(input)
	p := newParser(l)
	json := &JSON{json: p.parse()}

	errors := p.getErrors()
	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return json, nil
}

func (j *JSON) String() string {
	return j.json.String()
}

// TODO
func (j *JSON) Get(path string) interface{} {
	return nil
}
