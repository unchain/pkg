package xutil_test

import (
	"github.com/unchainio/pkg/xutil"
	"testing"
)

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xutil.RandomString(10)
	}
}

func BenchmarkRandom2String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xutil.RandomString2(10)
	}
}
