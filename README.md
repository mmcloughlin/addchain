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

Results for common cryptographic exponents and delta compared to [best known
hand-optimized addition
chains](https://briansmith.org/ecc-inversion-addition-chains-01).

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


See [full results listing](doc/results.md) for more detail and additional
exponents.

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

Algorithms all conform to the [`alg.ChainAlgorithm`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg#ChainAlgorithm) or
[`alg.SequenceAlgorithm`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg#SequenceAlgorithm) interfaces and can be used directly. However the
most user-friendly method uses the [`alg/ensemble`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/ensemble) package to
instantiate a sensible default set of algorithms and the [`alg/exec`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/exec)
helper to execute them in parallel. The following code uses this method to
find an addition chain for curve25519 field inversion:

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

## Algorithms

### Binary

The [`alg/binary`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/binary) package implements the addition chain equivalent
of the basic [square-and-multiply exponentiation
method](https://en.wikipedia.org/wiki/Exponentiation_by_squaring). It is
included for completeness, but is almost always outperformed by more advanced
algorithms below.

### Continued Fractions

The [`alg/contfrac`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/contfrac) package implements the continued fractions
methods for addition sequence search introduced by
Bergeron-Berstel-Brlek-Duboc in 1989 and later extended. This approach
utilizes a decomposition of an addition chain akin to continued fractions,
namely

```
(1,..., k,..., n) = (1,...,n mod k,..., k) ⊗ (1,..., n/k) ⊕ (n mod k).
```

for certain special operators ⊗ and ⊕. This
decomposition lends itself to a recursive algorithm for efficient addition
sequence search, with results dependent on the _strategy_ for choosing the
auxillary integer _k_. The [`alg/contfrac`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/contfrac) package provides a
laundry list of strategies from the literature: binary, co-binary,
dichotomic, dyadic, fermat, square-root and total.

#### References

* F Bergeron, J Berstel, S Brlek and C Duboc. Addition chains using continued fractions. Journal of Algorithms. 1989. http://www-igm.univ-mlv.fr/~berstel/Articles/1989AdditionChainDuboc.pdf
* Bergeron, F., Berstel, J. and Brlek, S. Efficient computation of addition chains. Journal de theorie des nombres de Bordeaux. 1994. http://www.numdam.org/item/JTNB_1994__6_1_21_0
* Amadou Tall and Ali Yassin Sanghare. Efficient computation of addition-subtraction chains using generalized continued Fractions. Cryptology ePrint Archive, Report 2013/466. 2013. https://eprint.iacr.org/2013/466
* Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve Cryptography, chapter 9. 2006. http://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

### Bos-Coster Heuristics

Bos and Coster described an iterative algorithm for efficient addition
sequence generation in which at each step a heuristic proposes new numbers
for the sequence in such a way that the _maximum_ number always decreases.
The [original Bos-Coster paper](https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf) defined four
heuristics: Approximation, Divison, Halving and Lucas. Package
[`alg/heuristic`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/heuristic) implements a variation on these heuristics:

* **Approximation**: looks for two elements a, b in the current sequence with sum close to the largest element.
* **Halving**: applies when the target is at least twice as big as the next largest, and if so it will propose adding a sequence of doublings.
* **Delta Largest**: proposes adding the delta between the largest two entries in the current sequence.

Divison and Lucas are not implemented due to disparities in the literature
about their precise definition and poor results from early experiments.
Furthermore, this library does not apply weights to the heuristics as
suggested in the paper, rather it simply uses the first that applies. However
both of these remain [possible avenues for
improvement](https://github.com/mmcloughlin/addchain/issues/26).

#### References

* Bos, Jurjen and Coster, Matthijs. Addition Chain Heuristics. In Advances in Cryptology --- CRYPTO' 89 Proceedings, pages 400--407. 1990. https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf
* Riad S. Wahby. kwantam/addchain. Github Repository. Apache License, Version 2.0. 2018. https://github.com/kwantam/addchain
* Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve Cryptography, chapter 9. 2006. http://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf
* Ayan Nandy. Modifications of Bos and Coster’s Heuristics in search of a shorter addition chain for faster exponentiation. Masters thesis, Indian Statistical Institute Kolkata. 2011. http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf
* F. L. Ţiplea, S. Iftene, C. Hriţcu, I. Goriac, R. Gordân and E. Erbiceanu. MpNT: A Multi-Precision Number Theory Package, Number Theoretical Algorithms (I). Technical Report TR03-02, Faculty of Computer Science, "Alexandru Ioan Cuza" University, Iasi. 2003. https://profs.info.uaic.ro/~tr/tr03-02.pdf
* Stam, Martijn. Speeding up subgroup cryptosystems. PhD thesis, Technische Universiteit Eindhoven. 2003. https://cr.yp.to/bib/2003/stam-thesis.pdf

### Dictionary

Dictionary methods decompose the binary representation of a target integer _n_ into a set of dictionary _terms_, such that _n_
may be written as a sum

```
n = ∑ 2^{e_i} d_i
```

for exponents _e_ and elements _d_ from a dictionary _D_. Given such a decomposition we can construct an addition chain for _n_ by

1. Find a short addition _sequence_ containing every element of the dictionary _D_. Continued fractions and Bos-Coster heuristics can be used here.
2. Build _n_ from the dictionary terms using the sum decomposition.

The efficiency of this approach depends on the decomposition method. The [`alg/dict`](https://pkg.go.dev/github.com/mmcloughlin/addchain/alg/dict) package provides:

* **Fixed Window**:
* **Sliding Window**:
* **Run Length**:
* **Hybrid**:

#### References

* Martin Otto. Brauer addition-subtraction chains. PhD thesis, Universitat Paderborn. 2001. http://www.martin-otto.de/publications/docs/2001_MartinOtto_Diplom_BrauerAddition-SubtractionChains.pdf
* Kunihiro, Noboru and Yamamoto, Hirosuke. New Methods for Generating Short Addition Chains. IEICE Transactions on Fundamentals of Electronics Communications and Computer Sciences. 2000. https://pdfs.semanticscholar.org/b398/d10faca35af9ce5a6026458b251fd0a5640c.pdf
* Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve Cryptography, chapter 9. 2006. http://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

### Runs

### Optimization

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
