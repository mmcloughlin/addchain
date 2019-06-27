package addchain

import "testing"

/*
func TestHeuristicSequenceAlgorithm(t *testing.T) {
	a := NewHeuristicSequenceAlgorithm(LastTwoDelta{})
	targets := bigints.Int64s([]int64{1, 2, 47, 117, 343, 499, 933, 5689})
	c, err := a.FindSequence(targets)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("length=%d sequence=%v", len(c), c)
}
*/

func TestSequenceAlgorithms(t *testing.T) {
	as := []SequenceAlgorithm{
		// NewHeuristicSequenceAlgorithm(LastTwoDelta{}),
	}
	for _, a := range as {
		t.Run(a.String(), SequenceAlgorithmSuite(a))
	}
}

func SequenceAlgorithmSuite(a SequenceAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
	}
}
