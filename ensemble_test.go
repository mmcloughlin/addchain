package addchain

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/test"
	"github.com/mmcloughlin/addchain/prime"
)

// References:
//
//	[curvechains]  Brian Smith. The Most Efficient Known Addition Chains for Field Element and
//	               Scalar Inversion for the Most Popular and Most Unpopular Elliptic Curves. 2017.
//	               https://briansmith.org/ecc-inversion-addition-chains-01 (accessed June 30, 2019)

// TestEnsembleResultsInversionChains tests our methods against inversion
// exponents for various interesting or popular fields. This is used to set a
// baseline to prevent regressions, as well as to compare against the best
// hand-crafted chains [curvechains].
func TestEnsembleResultsInversionChains(t *testing.T) {
	cases := []struct {
		Name  string
		N     *big.Int
		Delta int64

		// Baseline program length obtained by previous versions, to prevent
		// regressions and highlight improvements.
		Baseline int

		// BestPublished is the length of the most efficient chain seen published
		// somewhere. These are currently all from [curvechains], it's possible there
		// are better results elsewhere.
		BestPublished int
	}{
		{
			Name:          "curve25519_field",
			N:             prime.P25519.Int(),
			Delta:         2,
			Baseline:      266,
			BestPublished: 265,
		},
		{
			Name:          "p256_field",
			N:             prime.NISTP256.Int(),
			Delta:         3,
			Baseline:      266,
			BestPublished: 266,
		},
		{
			Name:          "p384_field",
			N:             prime.NISTP384.Int(),
			Delta:         3,
			Baseline:      397,
			BestPublished: 396,
		},
		{
			Name:          "secp256k1_field",
			N:             prime.Secp256k1.Int(),
			Delta:         3,
			Baseline:      269,
			BestPublished: 269,
		},
		{
			Name:          "secp256k1_scalar",
			N:             bigint.MustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036413f"),
			Baseline:      293,
			BestPublished: 290,
		},
		{
			Name:          "p256_scalar",
			N:             bigint.MustHex("ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc63254f"),
			Baseline:      294,
			BestPublished: 292,
		},
		{
			Name:          "p384_scalar",
			N:             bigint.MustHex("ffffffffffffffffffffffffffffffffffffffffffffffffc7634d81f4372ddf581a0db248b0a77aecec196accc52971"),
			Baseline:      434,
			BestPublished: 433,
		},
		{
			Name:          "curve25519_scalar",
			N:             bigint.MustHex("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb"),
			Baseline:      283,
			BestPublished: 284,
		},
		{
			Name:     "p2213_field",
			N:        prime.P2213.Int(),
			Delta:    2,
			Baseline: 231,
		},
		{
			Name:     "p222117_field",
			N:        prime.P222117.Int(),
			Delta:    2,
			Baseline: 233,
		},
		{
			Name:     "p2519_field",
			N:        prime.P2519.Int(),
			Delta:    2,
			Baseline: 263,
		},
		{
			Name:     "p382105_field",
			N:        prime.P382105.Int(),
			Delta:    2,
			Baseline: 395,
		},
		{
			Name:     "p383187_field",
			N:        prime.P383187.Int(),
			Delta:    2,
			Baseline: 396,
		},
		{
			Name:     "p41417_field",
			N:        prime.P41417.Int(),
			Delta:    2,
			Baseline: 426,
		},
		{
			Name:     "p511187_field",
			N:        prime.P511187.Int(),
			Delta:    2,
			Baseline: 525,
		},
		{
			Name:     "p192_field",
			N:        prime.NISTP192.Int(),
			Delta:    2,
			Baseline: 203,
		},
		{
			Name:     "p224_field",
			N:        prime.NISTP224.Int(),
			Delta:    2,
			Baseline: 234,
		},
		{
			Name:     "goldilocks_field",
			N:        prime.Goldilocks.Int(),
			Delta:    2,
			Baseline: 460,
		},
		{
			Name:     "secp192k1_field",
			N:        prime.Secp192k1.Int(),
			Delta:    2,
			Baseline: 205,
		},
		{
			Name:     "secp224k1_field",
			N:        prime.Secp224k1.Int(),
			Delta:    2,
			Baseline: 238,
		},
	}
	as := Ensemble()
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			// Tests with a best known result are prioritized. Only run all tests in
			// stress test mode.
			if c.BestPublished == 0 {
				test.RequireStress(t)
			} else {
				test.RequireLong(t)
			}

			n := new(big.Int).Sub(c.N, big.NewInt(c.Delta))
			t.Logf("n-%d=%x", c.Delta, n)

			// Execute.
			p := NewParallel()
			p.SetLogger(log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds))
			rs := p.Execute(n, as)

			// Process results.
			best := 0
			for i, r := range rs {
				if r.Err != nil {
					t.Logf("error with %s", r.Algorithm)
					t.Error(r.Err)
					continue
				}
				if len(r.Program) <= len(rs[best].Program) {
					best = i
				}
			}

			// Report the best.
			b := rs[best]
			t.Logf(" best: %s", b.Algorithm)
			t.Logf("total: %d", len(b.Program))

			if c.BestPublished > 0 {
				t.Logf("known: %d", c.BestPublished)
				t.Logf("delta: %+d", len(b.Program)-c.BestPublished)
			}

			// Report on differences from baseline in either direction.
			if c.Baseline == 0 {
				t.Errorf("missing baseline: set to %d", len(b.Program))
			} else if len(b.Program) != c.Baseline {
				t.Errorf("change from baseline: %+d", len(b.Program)-c.Baseline)
			}
		})
	}
}
