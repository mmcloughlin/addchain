package ensemble

import (
	"testing"

	"github.com/mmcloughlin/addchain/internal/results"
)

func BenchmarkResults(b *testing.B) {
	as := Ensemble()
	for _, c := range results.Results {
		c := c // scopelint
		n := c.Target()
		b.Run(c.Slug, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, a := range as {
					a.FindChain(n)
				}
			}
		})
	}
}
