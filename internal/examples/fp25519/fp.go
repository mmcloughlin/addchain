package fp25519

import "math/big"

// p is the field prime modulus.
var p, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)

// Elt is an element of the field modulo 2^255-19.
type Elt struct{ n big.Int }

// Modp reduces z modulo p, ensuring it's in the range [0, p).
func Modp(z *Elt) {
	z.n.Mod(&z.n, p)
}

// Mul computes z = x*y (mod p).
func Mul(z, x, y *Elt) {
	z.n.Mul(&x.n, &y.n)
	Modp(z)
}

// Sqr computes z = x^2 (mod p).
func Sqr(z, x *Elt) { Mul(z, x, x) }
