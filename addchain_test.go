package addchain

import "testing"

func TestEvaluateDoubleChain(t *testing.T) {
	// Build a chain of doublings.
	c := Chain{}
	n := 17
	for i := 0; i < n; i++ {
		c.Add(i, i)
	}

	// Evaluate.
	x := c.Evaluate()
	for i, got := range x {
		if !got.IsUint64() {
			t.Fatal("expected to be representable as uint64")
		}
		expect := uint64(1) << uint(i)
		if got.Uint64() != expect {
			t.Fail()
		}
	}
}
