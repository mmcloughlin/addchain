// Package ints provides simple integer utilities.
package ints

// MinMax returns the min and max of x and y.
func MinMax(x, y int) (min, max int) {
	if x > y {
		return y, x
	}
	return x, y
}

// Min returns the minimum of x and y.
func Min(x, y int) int {
	min, _ := MinMax(x, y)
	return min
}

// Max returns the maximum of x and y.
func Max(x, y int) int {
	_, max := MinMax(x, y)
	return max
}
