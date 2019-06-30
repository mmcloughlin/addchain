package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/prime"
)

// References:
//
//	[curvechains]  Brian Smith. The Most Efficient Known Addition Chains for Field Element & Scalar
//	               Inversion for the Most Popular & Most Unpopular Elliptic Curves. 2017.
//	               https://briansmith.org/ecc-inversion-addition-chains-01

// TestEfficientInversionChains compares our methods against the efficient
// chains listed for inversion in [curvechains].
func TestEfficientInversionChains(t *testing.T) {
	cases := []struct {
		Name      string
		P         prime.Prime
		Delta     int64
		KnownBest int
	}{
		{
			Name:      "curve25519_field",
			P:         prime.P25519,
			Delta:     2,
			KnownBest: 265,
		},
		{
			Name:      "p256_field",
			P:         prime.NISTP256,
			Delta:     3,
			KnownBest: 266,
		},
		{
			Name:      "p384_field",
			P:         prime.NISTP384,
			Delta:     3,
			KnownBest: 396,
		},
	}
	as := Ensemble()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			q := new(big.Int).Sub(c.P.Int(), big.NewInt(c.Delta))
			t.Logf("p=%x", c.P.Int())
			t.Logf("q=%x", q)

			rs := Parallel(q, as)
			var best Program
			for _, r := range rs {
				t.Logf("algorithm: %s", r.Algorithm)
				if r.Err != nil {
					t.Fatal(r.Err)
					continue
				}
				doubles, adds := r.Program.Count()
				total := doubles + adds
				t.Logf("total: %d\tdoubles: \t%d adds: %d", total, doubles, adds)
				if best == nil || len(r.Program) < len(best) {
					best = r.Program
				}
			}
			t.Logf("min: %d bestknown: %d", len(best), c.KnownBest)
		})
	}
}
