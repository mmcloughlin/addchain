package fp25519

func Inv(z, x *Elt) {
	var (
		t0 = new(Elt)
		t1 = new(Elt)
		t2 = new(Elt)
	)

	Sqr(z, x)

	Mul(z, x, z)

	Sqr(t0, z)
	for s := 1; s < 2; s++ {
		Sqr(t0, t0)
	}

	Mul(t0, z, t0)

	Sqr(t1, t0)
	for s := 1; s < 4; s++ {
		Sqr(t1, t1)
	}

	Mul(t0, t0, t1)

	for s := 0; s < 2; s++ {
		Sqr(t0, t0)
	}

	Mul(t0, z, t0)

	Sqr(t1, t0)
	for s := 1; s < 10; s++ {
		Sqr(t1, t1)
	}

	Mul(t1, t0, t1)

	for s := 0; s < 10; s++ {
		Sqr(t1, t1)
	}

	Mul(t1, t0, t1)

	Sqr(t2, t1)
	for s := 1; s < 30; s++ {
		Sqr(t2, t2)
	}

	Mul(t1, t1, t2)

	Sqr(t2, t1)
	for s := 1; s < 60; s++ {
		Sqr(t2, t2)
	}

	Mul(t1, t1, t2)

	Sqr(t2, t1)
	for s := 1; s < 120; s++ {
		Sqr(t2, t2)
	}

	Mul(t1, t1, t2)

	for s := 0; s < 10; s++ {
		Sqr(t1, t1)
	}

	Mul(t0, t0, t1)

	for s := 0; s < 2; s++ {
		Sqr(t0, t0)
	}

	Mul(t0, x, t0)

	for s := 0; s < 3; s++ {
		Sqr(t0, t0)
	}

	Mul(z, z, t0)
}
