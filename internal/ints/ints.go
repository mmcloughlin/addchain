// Package ints provides simple integer utilities.
package ints

// Min returns the minimum of x and y.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum of x and y.
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
