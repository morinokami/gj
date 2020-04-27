# gj
 
gj is a simple JSON parser written in Go.

## Usage

```go
import "github.com/morinokami/gj"

input := `{
  "firstName": "John",
  "lastName": "Smith",
  "isAlive": true,
  "age": 27,
  "address": {
    "streetAddress": "21 2nd Street",
    "city": "New York",
    "state": "NY",
    "postalCode": "10021-3100"
  },
  "phoneNumbers": [
    {
      "type": "home",
      "number": "212 555-1234"
    },
    {
      "type": "office",
      "number": "646 555-4567"
    }
  ],
  "children": [],
  "spouse": null
}`

json, err := gj.ParseString(input)

city, err := json.Get("address.city")
number, err := json.Get("phoneNumbers.[1].number")
```
