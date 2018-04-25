package xrand_test

import (
	"testing"
)

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xrand.String(10)
	}
}

func BenchmarkRandom2String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xrand.String2(10)
	}
}
