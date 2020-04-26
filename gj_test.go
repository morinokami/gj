package gj

import "testing"

func TestParseString(t *testing.T) {

	checkResult := func(t *testing.T, val, got, want interface{}, ok bool) {
		if !ok {
			t.Fatalf("Cannot convert to %T: %T (%+v)", want, val, val)
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
			json, _ := ParseString(tt.input)
			val, _ := json.Get(tt.path)

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

	t.Run("Get errors", func(t *testing.T) {})

}
