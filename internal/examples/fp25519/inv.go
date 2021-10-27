// Code generated by addchain. DO NOT EDIT.

package fp25519

// Inv computes z = 1/x (mod p) and returns it.
func (z *Elt) Inv(x *Elt) *Elt {
	// Inversion computation is derived from the addition chain:
	//
	//	_10       = 2*1
	//	_11       = 1 + _10
	//	_1100     = _11 << 2
	//	_1111     = _11 + _1100
	//	_11110000 = _1111 << 4
	//	_11111111 = _1111 + _11110000
	//	x10       = _11111111 << 2 + _11
	//	x20       = x10 << 10 + x10
	//	x30       = x20 << 10 + x10
	//	x60       = x30 << 30 + x30
	//	x120      = x60 << 60 + x60
	//	x240      = x120 << 120 + x120
	//	x250      = x240 << 10 + x10
	//	return      (x250 << 2 + 1) << 3 + _11
	//
	// Operations: 254 squares 12 multiplies
	//
	// Generated by github.com/mmcloughlin/addchain v0.3.0.

	// Allocate Temporaries.
	var (
		t0 = new(Elt)
		t1 = new(Elt)
		t2 = new(Elt)
	)

	// Step 1: z = x^0x2
	z.Sqr(x)

	// Step 2: z = x^0x3
	z.Mul(x, z)

	// Step 4: t0 = x^0xc
	t0.Sqr(z)
	for s := 1; s < 2; s++ {
		t0.Sqr(t0)
	}

	// Step 5: t0 = x^0xf
	t0.Mul(z, t0)

	// Step 9: t1 = x^0xf0
	t1.Sqr(t0)
	for s := 1; s < 4; s++ {
		t1.Sqr(t1)
	}

	// Step 10: t0 = x^0xff
	t0.Mul(t0, t1)

	// Step 12: t0 = x^0x3fc
	for s := 0; s < 2; s++ {
		t0.Sqr(t0)
	}

	// Step 13: t0 = x^0x3ff
	t0.Mul(z, t0)

	// Step 23: t1 = x^0xffc00
	t1.Sqr(t0)
	for s := 1; s < 10; s++ {
		t1.Sqr(t1)
	}

	// Step 24: t1 = x^0xfffff
	t1.Mul(t0, t1)

	// Step 34: t1 = x^0x3ffffc00
	for s := 0; s < 10; s++ {
		t1.Sqr(t1)
	}

	// Step 35: t1 = x^0x3fffffff
	t1.Mul(t0, t1)

	// Step 65: t2 = x^0xfffffffc0000000
	t2.Sqr(t1)
	for s := 1; s < 30; s++ {
		t2.Sqr(t2)
	}

	// Step 66: t1 = x^0xfffffffffffffff
	t1.Mul(t1, t2)

	// Step 126: t2 = x^0xfffffffffffffff000000000000000
	t2.Sqr(t1)
	for s := 1; s < 60; s++ {
		t2.Sqr(t2)
	}

	// Step 127: t1 = x^0xffffffffffffffffffffffffffffff
	t1.Mul(t1, t2)

	// Step 247: t2 = x^0xffffffffffffffffffffffffffffff000000000000000000000000000000
	t2.Sqr(t1)
	for s := 1; s < 120; s++ {
		t2.Sqr(t2)
	}

	// Step 248: t1 = x^0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	t1.Mul(t1, t2)

	// Step 258: t1 = x^0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc00
	for s := 0; s < 10; s++ {
		t1.Sqr(t1)
	}

	// Step 259: t0 = x^0x3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	t0.Mul(t0, t1)

	// Step 261: t0 = x^0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc
	for s := 0; s < 2; s++ {
		t0.Sqr(t0)
	}

	// Step 262: t0 = x^0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd
	t0.Mul(x, t0)

	// Step 265: t0 = x^0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe8
	for s := 0; s < 3; s++ {
		t0.Sqr(t0)
	}

	// Step 266: z = x^0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeb
	z.Mul(z, t0)

	return z
}
