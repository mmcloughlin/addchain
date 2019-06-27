package test

import (
	"flag"
	"testing"
	"time"
)

// Custom test command line flags.
var (
	long     = flag.Bool("long", false, "allow more time for trial-based tests")
	stress   = flag.Bool("stress", false, "run stress tests for trial-based tests")
	duration = flag.Duration("duration", 0, "duration for trial-based tests")
)

// trialsduration returns how long trial-based tests can take.
func trialsduration() time.Duration {
	switch {
	case *duration > 0:
		return *duration
	case testing.Short():
		return time.Second / 10
	case *long:
		return 30 * time.Second
	case *stress:
		return 2 * time.Minute
	default:
		return time.Second
	}
}

// Repeat the given trial function. The duration is controlled by custom
// command-line flags. The trial function returns whether it wants to continue
// testing.
//
//	-short		run for less time than usual
//	-long		allow more time
//	-stress		run for an extremely long time
//	-duration	set a custom duration
func Repeat(t *testing.T, trial func(t *testing.T) bool) {
	start := time.Now()
	d := trialsduration()
	n := 1
	for time.Since(start) < d && trial(t) {
		n++
	}
	t.Logf("%d trials in %s", n, time.Since(start))
}

// Trials returns a function that repeats f.
func Trials(f func(t *testing.T) bool) func(t *testing.T) {
	return func(t *testing.T) {
		Repeat(t, f)
	}
}
