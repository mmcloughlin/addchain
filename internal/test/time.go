package test

import (
	"flag"
	"testing"
	"time"
)

// Custom test command line flags.
var (
	long   = flag.Bool("long", false, "enable long running tests")
	stress = flag.Bool("stress", false, "enable stress tests (implies -long)")
)

// timeallowed returns how long a single test is allowed to take.
func timeallowed() time.Duration {
	switch {
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

// Long reports whether long tests are enabled.
func Long() bool {
	return *long || *stress
}

// Stress reports whether stress tests are enabled.
func Stress() bool {
	return *stress
}

// RequireLong marks this test as a long test. Test will be skipped if long
// tests are not enabled.
func RequireLong(t *testing.T) {
	if !Long() {
		t.Skipf("long test: use -long or -stress to enable")
	}
}

// RequireStress marks this test as a stress test. Test will be skipped if stress
// tests are not enabled.
func RequireStress(t *testing.T) {
	if !Stress() {
		t.Skipf("stress test: use -stress to enable")
	}
}

type Timer struct {
	start    time.Time
	duration time.Duration
}

func NewTimer() *Timer {
	return &Timer{
		start:    time.Now(),
		duration: timeallowed(),
	}
}

func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

func (t *Timer) Done() bool {
	return t.Elapsed() >= t.duration
}

func (t *Timer) Check(test *testing.T) {
	if t.Done() {
		test.Skip("time limit reached")
	}
}

// Repeat the given trial function. The duration is controlled by custom
// command-line flags. The trial function returns whether it wants to continue
// testing.
//
//	-short		run for less time than usual
//	-long		allow more time
//	-stress		run for an extremely long time
func Repeat(t *testing.T, trial func(t *testing.T) bool) {
	timer := NewTimer()
	n := 1
	for !timer.Done() && trial(t) {
		n++
	}
	t.Logf("%d trials in %s", n, timer.Elapsed())
}

// Trials returns a function that repeats f.
func Trials(f func(t *testing.T) bool) func(t *testing.T) {
	return func(t *testing.T) {
		Repeat(t, f)
	}
}
