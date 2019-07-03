package addchain

import (
	"math/big"
	"testing"
	"time"

	"github.com/mmcloughlin/addchain/internal/bigint"

	"github.com/mmcloughlin/addchain/internal/test"

	"github.com/mmcloughlin/addchain/prime"
)

// References:
//
//	[curvechains]  Brian Smith. The Most Efficient Known Addition Chains for Field Element and
//	               Scalar Inversion for the Most Popular and Most Unpopular Elliptic Curves. 2017.
//	               https://briansmith.org/ecc-inversion-addition-chains-01 (accessed June 30, 2019)

// TestEfficientInversionChains compares our methods against the efficient
// chains listed for inversion in [curvechains].
func TestEfficientInversionChains(t *testing.T) {
	cases := []struct {
		Name        string
		N           *big.Int
		Delta       int64
		HandCrafted int
	}{
		{
			Name:        "curve25519_field",
			N:           prime.P25519.Int(),
			Delta:       2,
			HandCrafted: 265,
		},
		{
			Name:        "p256_field",
			N:           prime.NISTP256.Int(),
			Delta:       3,
			HandCrafted: 266,
		},
		{
			Name:        "p384_field",
			N:           prime.NISTP384.Int(),
			Delta:       3,
			HandCrafted: 396,
		},
		{
			Name:        "secp256k1_field",
			N:           prime.Secp256k1.Int(),
			Delta:       3,
			HandCrafted: 269,
		},
		{
			Name:        "secp256k1_scalar",
			N:           bigint.MustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036413f"),
			HandCrafted: 290,
		},
		{
			Name:        "p256_scalar",
			N:           bigint.MustHex("ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc63254f"),
			HandCrafted: 292,
		},
		{
			Name:        "p384_scalar",
			N:           bigint.MustHex("ffffffffffffffffffffffffffffffffffffffffffffffffc7634d81f4372ddf581a0db248b0a77aecec196accc52971"),
			HandCrafted: 433,
		},
		{
			Name:        "curve25519_scalar",
			N:           bigint.MustHex("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb"),
			HandCrafted: 284,
		},
	}
	as := Ensemble()
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			test.RequireDuration(t, 30*time.Second)

			n := new(big.Int).Sub(c.N, big.NewInt(c.Delta))
			t.Logf("n-%d=%x", c.Delta, n)

			rs := Parallel(n, as)
			best := 0
			for i, r := range rs {
				if r.Err != nil {
					t.Fatal(r.Err)
					continue
				}
				if len(r.Program) <= len(rs[best].Program) {
					best = i
				}
			}
			b := rs[best]
			t.Logf("algorithm: %s", b.Algorithm)
			t.Logf("total: %d", len(b.Program))
			t.Logf(" hand: %d", c.HandCrafted)
			t.Logf("delta: %+d", len(b.Program)-c.HandCrafted)
		})
	}
}
