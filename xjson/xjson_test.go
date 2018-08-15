package xjson_test

import (
	"reflect"
	"testing"

	"github.com/unchainio/pkg/xjson"
)

type U struct {
	C string
	D string
}

type T struct {
	A string
	B string
	U U
}

type TestCase struct {
	Data      []byte
	Whitelist []string
	Expected  *T
}

func TestUnmarshal(t *testing.T) {
	cases := []TestCase{
		{
			Data:      []byte("{\"A\": \"a\", \"B\":\"b\"}"),
			Whitelist: []string{"A"},
			Expected:  &T{A: "a"},
		},
		{
			Data:      []byte("{\"A\": \"a\", \"B\":\"b\"}"),
			Whitelist: []string{"B"},
			Expected:  &T{B: "b"},
		},
		{
			Data:      []byte("{\"A\": \"a\", \"B\":\"b\"}"),
			Whitelist: []string{"A", "B"},
			Expected:  &T{A: "a", B: "b"},
		},
		{
			Data:      []byte("{\"A\": \"a\", \"B\":\"b\", \"U\": {\"C\":\"c\", \"D\":\"d\"}}"),
			Whitelist: []string{"A", "B", "U"},
			Expected:  &T{A: "a", B: "b"},
		},
	}

	for _, tc := range cases {
		v := new(T)
		xjson.Unmarshal(tc.Data, v, tc.Whitelist...)
		reflect.DeepEqual(v, tc.Expected)
	}
}
