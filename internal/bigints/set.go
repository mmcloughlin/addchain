package bigints

import (
	"math/big"
	"strconv"
	"sync"
)

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
	s.set[setkey(x)] = true
}

func (s *Set) Contains(x *big.Int) bool {
	b := bufpool.Get().(*buffer)
	defer bufpool.Put(b)
	b.buf = appendkey(b.buf[:0], x)
	return s.set[string(b.buf)]
}

func setkey(x *big.Int) string {
	var buf []byte
	buf = appendkey(buf, x)
	return string(buf)
}

func appendkey(buf []byte, x *big.Int) []byte {
	w := x.Bits()

	// Skip leading zeros.
	i := len(w) - 1
	for ; i >= 0 && w[i] == 0; i-- {
	}

	// Format remaining.
	for ; i >= 0; i-- {
		buf = strconv.AppendUint(buf, uint64(w[i]), 16)
	}

	return buf
}

type buffer struct {
	buf []byte
}

var bufpool = sync.Pool{
	New: func() interface{} {
		return &buffer{
			buf: make([]byte, 0, 32),
		}
	},
}
