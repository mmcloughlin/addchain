package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigints"
	"github.com/mmcloughlin/addchain/internal/container/queue"
)

// SequenceAlgorithm is a method of generating an addition sequence for a set of
// target values.
type SequenceAlgorithm interface {
	fmt.Stringer
	FindSequence(targets []*big.Int) ([]*big.Int, error)
}

// SequenceState represents a current state in a search for a addition sequence.
type SequenceState struct {
	Proto []*big.Int // remaining elements to produce
	Chain []*big.Int // chain generated so far
}

// NewInitialSequenceState builds an initial sequence state for the given
// targets. The protosequence is the list of targets together with {1,2}, sorted
// and uniqued. The initial chain is empty.
func NewInitialSequenceState(targets []*big.Int) *SequenceState {
	proto := append([]*big.Int{big.NewInt(1), big.NewInt(2)}, targets...)
	bigints.Sort(proto)
	proto = bigints.Unique(proto)

	return &SequenceState{
		Proto: proto,
		Chain: []*big.Int{},
	}
}

// Complete reports whether the state is complete. That is, if the protosequence only contains {1,2}.
func (s *SequenceState) Complete() bool {
	return len(s.Proto) <= 2
}

// Target returns the next target: the top element of the protosequence.
func (s *SequenceState) Target() *big.Int {
	top := len(s.Proto) - 1
	return s.Proto[top]
}

// MoveTargetToChain moves the target element (top element in protosequence) to the chain.
func (s *SequenceState) MoveTargetToChain() {
	top := len(s.Proto) - 1
	target := s.Proto[top]
	s.Proto = s.Proto[:top]
	s.Chain = bigints.InsertSortedUnique(s.Chain, target)
}

// Score is an estimate for how long the final chain will be.
func (s *SequenceState) Score() float64 {
	log := bigints.Max(s.Proto).BitLen()
	remaining := 1.5*float64(log) + float64(len(s.Proto))
	return remaining + float64(len(s.Chain))
}

// Proposal is a suggested "move" to apply to a sequence state.
//
// TODO(mbm): are all proposals insertions?
type Proposal struct {
	Insert []*big.Int
}

func (p Proposal) Apply(s *SequenceState) *SequenceState {
	return &SequenceState{
		Proto: bigints.MergeUnique(p.Insert, s.Proto),
		Chain: bigints.Clone(s.Chain),
	}
}

// Heuristic suggests moves from a given sequence state.
type Heuristic interface {
	fmt.Stringer
	Suggest(*SequenceState) []*Proposal
}

// HeuristicSequenceAlgorithm searches for an addition sequence with a
// collection of heuristics.
type HeuristicSequenceAlgorithm struct {
	heuristics []Heuristic
}

func NewHeuristicSequenceAlgorithm(heuristics ...Heuristic) *HeuristicSequenceAlgorithm {
	h := &HeuristicSequenceAlgorithm{}
	for _, heuristic := range heuristics {
		h.AddHeuristic(heuristic)
	}
	return h
}

func (h *HeuristicSequenceAlgorithm) AddHeuristic(heuristic Heuristic) {
	h.heuristics = append(h.heuristics, heuristic)
}

func (h HeuristicSequenceAlgorithm) String() string {
	return fmt.Sprintf("heuristic(%v)", h.heuristics)
}

// FindSequence searches for an addition sequence for the given targets.
func (h HeuristicSequenceAlgorithm) FindSequence(targets []*big.Int) ([]*big.Int, error) {
	// Initialize priority queue.
	initial := NewInitialSequenceState(targets)
	q := queue.NewPriority()
	q.Insert(initial, initial.Score())

	for !q.Empty() {
		// Pop off the current best state.
		s := q.Pop().(*SequenceState)
		if s.Complete() {
			return s.Chain, nil
		}

		// Apply heuristics.
		for _, heuristic := range h.heuristics {
			proposals := heuristic.Suggest(s)
			for _, proposal := range proposals {
				t := proposal.Apply(s)
				t.MoveTargetToChain()
				if t.Complete() {
					return t.Chain, nil
				}
				q.Insert(t, t.Score())
			}
		}
	}

	return nil, errors.New("failed to find sequence")
}
