package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigints"
)

type LastTwoDelta struct{}

func (d LastTwoDelta) String() string { return "last_two_delta" }

func (d LastTwoDelta) Suggest(s *SequenceState) []*Proposal {
	f := s.Proto
	n := len(f)
	delta := new(big.Int).Sub(f[n-1], f[n-2])
	propose := &Proposal{
		Insert: []*big.Int{delta},
	}
	return []*Proposal{propose}
}

func TestHeuristicSequenceAlgorithm(t *testing.T) {
	a := NewHeuristicSequenceAlgorithm(LastTwoDelta{})
	targets := bigints.Int64s([]int64{1, 2, 47, 117, 343, 499, 933, 5689})
	c, err := a.FindSequence(targets)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("length=%d sequence=%v", len(c), c)
}
