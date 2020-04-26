package gj

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type JSON struct {
	json *jsonExpression
}

func ParseString(input string) (*JSON, error) {
	l := newLexer(input)
	p := newParser(l)
	json := &JSON{json: p.parse()}

	err := p.getErrors()
	if len(err) > 0 {
		return nil, errors.New(strings.Join(err, "\n"))
	}

	return json, nil
}

func (j *JSON) String() string {
	return j.json.String()
}

func (j *JSON) Get(path string) (interface{}, error) {
	json := evalExpression(j.json.Value)
	if len(path) == 0 {
		return json, nil
	}

	keys := strings.Split(path, ".")
	for _, key := range keys {
		if strings.HasPrefix(key, "[") {
			if key == "[" || key == "[]" {
				return nil, fmt.Errorf(`index error - "%s"`, key)
			}

			index, err := strconv.ParseInt(key[1:len(key)-1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf(`index error - "%s"`, key)
			}
			if index < 0 {
				return nil, errors.New("index error - index out of bounds")
			}

			obj, ok := json.([]interface{})
			if !ok {
				return nil, errors.New(`index error - cannot use "[]"`)
			}

			if index >= int64(len(obj)) {
				return nil, errors.New("index error - index out of bounds")
			}

			json = obj[index]
		} else {
			obj, ok := json.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf(`key error - "%s"`, key)
			}

			json, ok = obj[key]
			if !ok {
				return nil, fmt.Errorf(`key error - "%s"`, key)
			}
		}
	}

	return json, nil
}

func evalExpression(exp expression) interface{} {
	switch value := exp.(type) {
	case *booleanExpression:
		return value.Value
	case *nullExpression:
		return value.Value
	case *integerExpression:
		return value.Value
	case *floatExpression:
		f, _ := strconv.ParseFloat(value.Value, 64)
		return f
	case *prefixExpression:
		switch right := value.Right.(type) {
		case *integerExpression:
			return -right.Value
		case *floatExpression:
			f, _ := strconv.ParseFloat(right.Value, 64)
			return -f
		default:
			// TODO: not implemented
			return nil
		}
	case *stringExpression:
		return value.Value
	case *objectExpression:
		o := make(map[string]interface{})
		for key := range value.Pairs {
			o[key] = evalExpression(value.Pairs[key])
		}
		return o
	case *arrayExpression:
		var a []interface{}
		for _, elem := range value.Values {
			a = append(a, evalExpression(elem))
		}
		return a
	default:
		// TODO: not implemented
		return nil
	}
}
