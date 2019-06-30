package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigint"
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
		N         *big.Int
		Delta     int64
		KnownBest int
	}{
		{
			Name:      "curve25519_field",
			N:         prime.P25519.Int(),
			Delta:     2,
			KnownBest: 265,
		},
		{
			Name:      "p256_field",
			N:         prime.NISTP256.Int(),
			Delta:     3,
			KnownBest: 266,
		},
		{
			Name:      "p384_field",
			N:         prime.NISTP384.Int(),
			Delta:     3,
			KnownBest: 396,
		},
		{
			Name:      "secp256k1_scalar",
			N:         bigint.MustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036413f"),
			KnownBest: 290,
		},
		{
			Name:      "p256_scalar",
			N:         bigint.MustHex("ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc63254f"),
			KnownBest: 292,
		},
		{
			Name:      "p384_scalar",
			N:         bigint.MustHex("ffffffffffffffffffffffffffffffffffffffffffffffffc7634d81f4372ddf581a0db248b0a77aecec196accc52971"),
			KnownBest: 433,
		},
		{
			Name:      "curve25519_scalar",
			N:         bigint.MustHex("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb"),
			KnownBest: 284,
		},
	}
	as := Ensemble()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			n := new(big.Int).Sub(c.N, big.NewInt(c.Delta))
			t.Logf("n-%d=%x", c.Delta, n)

			rs := Parallel(n, as)
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
