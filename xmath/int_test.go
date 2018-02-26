package xmath

import "testing"

var numsList = [][]int{
	{-10, 2, 4, 5, 15},
	{9, 8, 7, 6, 5, 4, 3, 2, 1},
	{123, 433, 4545, 512, -99999},
}

var expectedMins = []int{-10, 1, -99999}
var expectedMaxs = []int{15, 9, 4545}

func TestMinInt(t *testing.T) {
	for i := range numsList {
		min := MinInt(numsList[i]...)

		if expectedMins[i] != min {
			t.Fatalf("expected %d, received %d\n", expectedMins[i], min)
		}
	}
}

func TestMaxInt(t *testing.T) {
	for i := range numsList {
		max := MaxInt(numsList[i]...)

		if expectedMaxs[i] != max {
			t.Fatalf("expected %d, received %d\n", expectedMaxs[i], max)
		}
	}
}
