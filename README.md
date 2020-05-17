<p align="center">
  <img src="logo.svg" width="40%" border="0" alt="addchain" />
</p>

<p align="center">Cryptographic Addition Chain Generation in Go</p>

`addchain` generates short addition chains for exponents of cryptographic
interest with [results](#results) rivaling the best hand-optimized chains.
Intended to form a building block in elliptic curve or other cryptographic
code generators.

* Suite of algorithms from academic research: continued fractions,
  dictionary-based and Bos-Coster heuristics
* Custom run-length techniques exploit structure of cryptographic exponents
  with excellent results on Solinas primes
* Generic optimization methods eliminate redundant operations
* Simple domain-specific language for addition chain computations
* Command-line tool or library: use as a building block in cryptographic code
  generators

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

### Command-line Tool

Install with:

```
go get -u github.com/mmcloughlin/addchain/cmd/addchain
```
