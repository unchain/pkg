package xmath

import "github.com/ethereum/go-ethereum/common/math"

func MinInt(nums ...int) int {
	min := math.MaxInt32

	for _, num := range nums {
		if min > num {
			min = num
		}
	}

	return min
}

func MaxInt(nums ...int) int {
	max := math.MinInt32

	for _, num := range nums {
		if max < num {
			max = num
		}
	}

	return max
}
