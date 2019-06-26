package addchain

import "testing"

func TestProgramEvaluateDoublings(t *testing.T) {
	// Build a chain of doublings.
	p := Program{}
	n := 17
	for i := 0; i < n; i++ {
		p.Double(i)
	}

	// Evaluate.
	c := p.Evaluate()
	for i, got := range c {
		if !got.IsUint64() {
			t.Fatal("expected to be representable as uint64")
		}
		expect := uint64(1) << uint(i)
		if got.Uint64() != expect {
			t.Fail()
		}
	}
}
