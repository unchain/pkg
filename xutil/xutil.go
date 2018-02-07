package xutil

import (
	"math/rand"
	"time"
)

// MultiAppendString appends multiple string slices to a string slice
func MultiAppendString(slice []string, elems ...[]string) []string {
	for _, elem := range elems {
		slice = append(slice, elem...)
	}
	
	return slice
}

// Flatten transforms a two dimensional slice to a single dimensional one
func Flatten(strss [][]string) []string {
	slice := make([]string, 10)
	
	for _, strs := range strss {
		slice = append(slice, strs...)
		//
		//for _, str := range strs {
		//	slice = append(slice, str)
		//}
	}
	
	return slice
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString generates a random string with the specified length
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	
	return string(b)
}
