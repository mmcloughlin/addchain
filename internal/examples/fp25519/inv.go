package fp25519

func (z *Elt) Inv(x *Elt) *Elt {
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
