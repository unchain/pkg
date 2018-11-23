package xslice

func ContainsString(slice []string, item string) bool {
	for i := range slice {
		if slice[i] == item {
			return true
		}
	}

	return false
}

func ContainsInt(slice []int, item int) bool {
	for i := range slice {
		if slice[i] == item {
			return true
		}
	}

	return false
}
