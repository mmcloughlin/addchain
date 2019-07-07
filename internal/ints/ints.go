// Package ints provides simple integer utilities.
package ints

// NextMultiple returns the next multiple of n greater than or equal to a.
func NextMultiple(a, n int) int {
	a += n - 1
	a -= a % n
	return a - (a % n)
}

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
