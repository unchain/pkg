package xjson_test

import (
	"testing"

	"github.com/unchainio/pkg/xjson"
	"v/github.com/docker/docker@v0.0.0-20170524085120-eef6495eddab/pkg/testutil/assert"
)

type T struct {
	A string
	B string
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
	}

	for _, tc := range cases {
		v := new(T)
		xjson.Unmarshal(tc.Data, v, tc.Whitelist...)
		assert.DeepEqual(t, v, tc.Expected)
	}
}
