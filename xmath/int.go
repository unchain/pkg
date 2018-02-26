package xmath

func MinInt(x int, y int) int {
	if x <= y {
		return x
	}

	return y
}

func MaxInt(x int, y int) int {
	if x >= y {
		return x
	}

	return y
}
