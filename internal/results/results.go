package results

import (
	"math/big"

	"github.com/mmcloughlin/addchain/internal/bigint"
	"github.com/mmcloughlin/addchain/internal/prime"
)

// References:
//
//	[curvechains]  Brian Smith. The Most Efficient Known Addition Chains for Field Element and
//	               Scalar Inversion for the Most Popular and Most Unpopular Elliptic Curves. 2017.
//	               https://briansmith.org/ecc-inversion-addition-chains-01 (accessed June 30, 2019)

// Result represents the best result of this library on a particular input.
type Result struct {
	Name  string
	N     *big.Int
	Delta int64

	// Best is the length of the most efficient chain produced by this library.
	Best int

	// BestPublished is the length of the most efficient chain seen published
	// somewhere. These are currently all from [curvechains], it's possible there
	// are better results elsewhere.
	BestPublished int
}

// Results on inversion exponents for various interesting or popular fields.
// These results set a baseline to prevent regressions, as well as to compare
// against the best hand-crafted chains [curvechains].
var Results = []Result{
	{
		Name:          "curve25519_field",
		N:             prime.P25519.Int(),
		Delta:         2,
		Best:          266,
		BestPublished: 265,
	},
	{
		Name:          "p256_field",
		N:             prime.NISTP256.Int(),
		Delta:         3,
		Best:          266,
		BestPublished: 266,
	},
	{
		Name:          "p384_field",
		N:             prime.NISTP384.Int(),
		Delta:         3,
		Best:          397,
		BestPublished: 396,
	},
	{
		Name:          "secp256k1_field",
		N:             prime.Secp256k1.Int(),
		Delta:         3,
		Best:          269,
		BestPublished: 269,
	},
	{
		Name:          "secp256k1_scalar",
		N:             bigint.MustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036413f"),
		Best:          293,
		BestPublished: 290,
	},
	{
		Name:          "p256_scalar",
		N:             bigint.MustHex("ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc63254f"),
		Best:          294,
		BestPublished: 292,
	},
	{
		Name:          "p384_scalar",
		N:             bigint.MustHex("ffffffffffffffffffffffffffffffffffffffffffffffffc7634d81f4372ddf581a0db248b0a77aecec196accc52971"),
		Best:          434,
		BestPublished: 433,
	},
	{
		Name:          "curve25519_scalar",
		N:             bigint.MustHex("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb"),
		Best:          283,
		BestPublished: 284,
	},
	{
		Name:  "p2213_field",
		N:     prime.P2213.Int(),
		Delta: 2,
		Best:  231,
	},
	{
		Name:  "p222117_field",
		N:     prime.P222117.Int(),
		Delta: 2,
		Best:  233,
	},
	{
		Name:  "p2519_field",
		N:     prime.P2519.Int(),
		Delta: 2,
		Best:  263,
	},
	{
		Name:  "p382105_field",
		N:     prime.P382105.Int(),
		Delta: 2,
		Best:  395,
	},
	{
		Name:  "p383187_field",
		N:     prime.P383187.Int(),
		Delta: 2,
		Best:  396,
	},
	{
		Name:  "p41417_field",
		N:     prime.P41417.Int(),
		Delta: 2,
		Best:  426,
	},
	{
		Name:  "p511187_field",
		N:     prime.P511187.Int(),
		Delta: 2,
		Best:  525,
	},
	{
		Name:  "p192_field",
		N:     prime.NISTP192.Int(),
		Delta: 2,
		Best:  203,
	},
	{
		Name:  "p224_field",
		N:     prime.NISTP224.Int(),
		Delta: 2,
		Best:  234,
	},
	{
		Name:  "goldilocks_field",
		N:     prime.Goldilocks.Int(),
		Delta: 2,
		Best:  460,
	},
	{
		Name:  "secp192k1_field",
		N:     prime.Secp192k1.Int(),
		Delta: 2,
		Best:  205,
	},
	{
		Name:  "secp224k1_field",
		N:     prime.Secp224k1.Int(),
		Delta: 2,
		Best:  238,
	},
}
