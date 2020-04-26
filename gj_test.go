package gj

import "testing"

func TestParser(t *testing.T) {
	input := `{
  "data": {
    "id": 888,
    "foo": "bar"
  }
}`

	t.Run("Get", func(t *testing.T) {
		json, _ := ParseString(input)

		got := json.Get("data.id")
		want := 888

		if got != want {
			t.Errorf("Get result wrong. got=%q, want=%d", got, want)
		}
	})

}
