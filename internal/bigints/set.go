package bigints

import "math/big"

type Set struct {
	set map[string]bool
}

func NewSetEmpty(xs ...*big.Int) *Set {
	return &Set{set: map[string]bool{}}
}

func NewSet(xs ...*big.Int) *Set {
	s := NewSetEmpty()
	for _, x := range xs {
		s.Add(x)
	}
	return s
}

func (s *Set) Add(x *big.Int) {
	s.set[s.key(x)] = true
}

func (s *Set) Contains(x *big.Int) bool {
	return s.set[s.key(x)]
}

func (s *Set) key(x *big.Int) string {
	return x.Text(32)
}
