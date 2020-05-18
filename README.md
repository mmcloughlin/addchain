<p align="center">
  <img src="logo.svg" width="40%" border="0" alt="addchain" />
</p>

<p align="center">Cryptographic Addition Chain Generation in Go</p>

`addchain` generates short addition chains for exponents of cryptographic
interest with [results](#results) rivaling the best hand-optimized chains.
Intended as a building block in elliptic curve or other cryptographic code
generators.

* Suite of algorithms from academic research: continued fractions,
  dictionary-based and Bos-Coster heuristics
* Custom run-length techniques exploit structure of cryptographic exponents
  with excellent results on Solinas primes
* Generic optimization methods eliminate redundant operations
* Simple domain-specific language for addition chain computations
* Command-line interface or library

## Results

| Name | Length | Best | Delta |
| ---- | -----: | ---: | ----: |
| [Curve25519 Field Inversion](doc/results.md#curve25519-field-inversion) | 266 | 265 | +1 |
| [NIST P-256 Field Inversion](doc/results.md#nist-p-256-field-inversion) | 266 | 266 | +0 |
| [NIST P-384 Field Inversion](doc/results.md#nist-p-384-field-inversion) | 397 | 396 | +1 |
| [secp256k1 (Bitcoin) Field Inversion](doc/results.md#secp256k1-bitcoin-field-inversion) | 269 | 269 | +0 |
| [Curve25519 Scalar Inversion](doc/results.md#curve25519-scalar-inversion) | 283 | 284 | -1 |
| [NIST P-256 Scalar Inversion](doc/results.md#nist-p-256-scalar-inversion) | 294 | 292 | +2 |
| [NIST P-384 Scalar Inversion](doc/results.md#nist-p-384-scalar-inversion) | 434 | 433 | +1 |
| [secp256k1 (Bitcoin) Scalar Inversion](doc/results.md#secp256k1-bitcoin-scalar-inversion) | 293 | 290 | +3 |


## Usage

### Command-line Interface

Install:

```
go get -u github.com/mmcloughlin/addchain/cmd/addchain
```

Search for a curve25519 field inversion addition chain with:

```sh
addchain search '2^255 - 19 - 2'
```

Output:

```
addchain: expr: "2^255 - 19 - 2"
addchain: hex: 7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeb
addchain: dec: 57896044618658097711785492504343953926634992332820282019728792003956564819947
addchain: best: opt(runs(continued_fractions(dichotomic)))
_10       = 2*1
_11       = 1 + _10
_1100     = _11 << 2
_1111     = _11 + _1100
_11110000 = _1111 << 4
_11111111 = _1111 + _11110000
x10       = _11111111 << 2 + _11
x20       = x10 << 10 + x10
x30       = x20 << 10 + x10
x60       = x30 << 30 + x30
x120      = x60 << 60 + x60
x240      = x120 << 120 + x120
x250      = x240 << 10 + x10
return      (x250 << 2 + 1) << 3 + _11
```

### Library

Install:

```
go get -u github.com/mmcloughlin/addchain
```

Algorithms all conform to the `alg.ChainAlgorithm` or `alg.SequenceAlgorithm`
interfaces and can be used directly. However the most user-friendly method
uses the `alg/ensemble` package to instantiate a sensible default set of
algorithms and the `alg/exec` helper to execute them in parallel. The
following code uses this method to find an addition chain for curve25519
field inversion:

```go
func Example() {
	// Target number: 2²⁵⁵ - 21.
	n := new(big.Int)
	n.SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeb", 16)

	// Default ensemble of algorithms.
	algorithms := ensemble.Ensemble()

	// Use parallel executor.
	ex := exec.NewParallel()
	results := ex.Execute(n, algorithms)

	// Output best result.
	best := 0
	for i, r := range results {
		if r.Err != nil {
			log.Fatal(r.Err)
		}
		if len(results[i].Program) < len(results[best].Program) {
			best = i
		}
	}
	r := results[best]
	fmt.Printf("best: %d\n", len(r.Program))
	fmt.Printf("algorithm: %s\n", r.Algorithm)

	// Output:
	// best: 266
	// algorithm: opt(runs(continued_fractions(dichotomic)))
}
```

## Thanks

Thank you to [Tom Dean](https://web.stanford.edu/~trdean/), [Riad
Wahby](https://wahby.org/) and [Brian Smith](https://briansmith.org/) for
advice and encouragement.

## Contributing

Contributions to `addchain` are welcome:

* [Submit bug reports](https://github.com/mmcloughlin/addchain/issues/new) to
  the issues page.
* Suggest [test cases](https://github.com/mmcloughlin/addchain/blob/e6c070065205efcaa02627ab1b23e8ce6aeea1db/internal/results/results.go#L62)
  or update best-known hand-optimized results.
* Pull requests accepted. Please discuss in the [issues section](https://github.com/mmcloughlin/addchain/issues)
  before starting significant work.

## License

`addchain` is available under the [BSD 3-Clause License](LICENSE).
