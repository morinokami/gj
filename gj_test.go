package gj

import "testing"

func TestParseString(t *testing.T) {

	t.Run("ParseString errors", func(t *testing.T) {
		tests := []struct {
			input         string
			expectedError string
		}{
			{"$", "no parse function for ILLEGAL found."},
			{"{", "expected next token to be STRING, got EOF instead."},
			{"{!}", "expected next token to be STRING, got ILLEGAL instead."},
		}

		for i, tt := range tests {
			_, err := ParseString(tt.input)

			if err == nil {
				t.Fatalf("[test %d] error expected, got none.", i)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("[test %d] unexpected error - got=%q, want=%q.", i, err.Error(), tt.expectedError)
			}
		}
	})

}

func TestGet(t *testing.T) {

	checkResult := func(t *testing.T, val, got, want interface{}, ok bool) {
		if !ok {
			t.Fatalf("cannot convert to %T: %T (%+v)", want, val, val)
		}
		if got != want {
			t.Errorf("Get result wrong. got=%q, want=%q", got, want)
		}
	}

	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			input    string
			path     string
			expected interface{}
		}{
			{"true", "", true},
			{"false", "", false},
			{"null", "", nil},
			{"1", "", int64(1)},
			{"3.14", "", 3.14},
			{`"foo"`, "", "foo"},
			{"-1", "", int64(-1)},
			{"-3.14", "", -3.14},
			{`{"foo": 2}`, "foo", int64(2)},
			{`{"foo": {"bar": 3}}`, "foo.bar", int64(3)},
			{"[1, 2, 3]", "[1]", int64(2)},
			{`{"foo": [1, 2, 3]}`, "foo.[0]", int64(1)},
			{`{"foo": [1, 2, 3]}`, "foo.[1]", int64(2)},
			{`{"foo": [1, 2, 3]}`, "foo.[2]", int64(3)},
			{`{"foo": [1, 2, {"bar": 3, "baz": 4}]}`, "foo.[2].baz", int64(4)},
		}

		for _, tt := range tests {
			json, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("unexpected error - %q", err.Error())
			}

			val, err := json.Get(tt.path)
			if err != nil {
				t.Fatalf("unexpected error - %q", err.Error())
			}

			switch want := tt.expected.(type) {
			case bool:
				got, ok := val.(bool)
				checkResult(t, val, got, want, ok)
			case nil:
				if val != want {
					t.Errorf("Get result wrong. got=%q, want=%q", val, want)
				}
			case int64:
				got, ok := val.(int64)
				checkResult(t, val, got, want, ok)
			case float64:
				got, ok := val.(float64)
				checkResult(t, val, got, want, ok)
			}
		}
	})

	t.Run("Get errors 1", func(t *testing.T) {
		input := "true"

		tests := []struct {
			path          string
			expectedError string
		}{
			{"meh", `key error - "meh"`},
			{"[0]", `index error - cannot use "[]"`},
		}

		for i, tt := range tests {
			json, _ := ParseString(input)
			_, err := json.Get(tt.path)

			if err == nil {
				t.Fatalf("[test %d] error expected, got none.", i)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("[test %d] unexpected error - got=%q, want=%q.", i, err.Error(), tt.expectedError)
			}
		}
	})

	t.Run("Get errors 2", func(t *testing.T) {
		input := `{
	"foo": 1,
	"bar": [
		true,
		false,
		null
	]
}`
		tests := []struct {
			path          string
			expectedError string
		}{
			{"meh", `key error - "meh"`},
			{"[0]", `index error - cannot use "[]"`},
			{"bar.[", `index error - "["`},
			{"bar.[]", `index error - "[]"`},
			{"bar.[-1]", "index error - index out of bounds"},
			{"bar.[3]", "index error - index out of bounds"},
			{"bar.[wow]", `index error - "[wow]"`},
		}

		for i, tt := range tests {
			json, _ := ParseString(input)
			_, err := json.Get(tt.path)

			if err == nil {
				t.Fatalf("[test %d] error expected, got none.", i)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("[test %d] unexpected error - got=%q, want=%q.", i, err.Error(), tt.expectedError)
			}
		}
	})

}
