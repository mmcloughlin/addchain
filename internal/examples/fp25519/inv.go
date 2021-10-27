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

	// Allocate Temporaries.
	var (
		t0 = new(Elt)
		t1 = new(Elt)
		t2 = new(Elt)
	)

	z.Sqr(x)

	z.Mul(x, z)

	t0.Sqr(z)
	for s := 1; s < 2; s++ {
		t0.Sqr(t0)
	}

	t0.Mul(z, t0)

	t1.Sqr(t0)
	for s := 1; s < 4; s++ {
		t1.Sqr(t1)
	}

	t0.Mul(t0, t1)

	for s := 0; s < 2; s++ {
		t0.Sqr(t0)
	}

	t0.Mul(z, t0)

	t1.Sqr(t0)
	for s := 1; s < 10; s++ {
		t1.Sqr(t1)
	}

	t1.Mul(t0, t1)

	for s := 0; s < 10; s++ {
		t1.Sqr(t1)
	}

	t1.Mul(t0, t1)

	t2.Sqr(t1)
	for s := 1; s < 30; s++ {
		t2.Sqr(t2)
	}

	t1.Mul(t1, t2)

	t2.Sqr(t1)
	for s := 1; s < 60; s++ {
		t2.Sqr(t2)
	}

	t1.Mul(t1, t2)

	t2.Sqr(t1)
	for s := 1; s < 120; s++ {
		t2.Sqr(t2)
	}

	t1.Mul(t1, t2)

	for s := 0; s < 10; s++ {
		t1.Sqr(t1)
	}

	t0.Mul(t0, t1)

	for s := 0; s < 2; s++ {
		t0.Sqr(t0)
	}

	t0.Mul(x, t0)

	for s := 0; s < 3; s++ {
		t0.Sqr(t0)
	}

	z.Mul(z, t0)

	return z
}
