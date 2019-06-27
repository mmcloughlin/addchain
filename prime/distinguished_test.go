package prime

import "testing"

func TestDistinguished(t *testing.T) {
	for _, p := range Distinguished {
		if !p.Int().ProbablyPrime(20) {
			t.Errorf("%s is not prime", p)
		}
	}
}
