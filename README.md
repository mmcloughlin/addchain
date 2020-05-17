<p align="center">
  <img src="logo.svg" width="40%" border="0" alt="addchain" />
</p>

<p align="center">Cryptographic Addition Chain Generation in Go</p>

`addchain` generates short addition chains for exponents of cryptographic
interest, with [results](#results) close to or exceeding the best
hand-optimized chains.

* Suite of algorithms from **academic research**: continued fractions methods
  of Bergeron-Berstel-Brlek-Duboc, dictionary-based approaches and Bos-Coster
  heuristics.
* **Novel techniques** exploiting structure of cryptographic exponents: custom
  run-length method with excellent results on Solinas primes

## Results

| Name | _N_ | _d_ | Length | Best | Delta |
| ---- | --- | --- | -----: | ---: | ----: |
| Curve25519 Field Inversion | `2^255-19` | 2 | 266 | 265 | +1 |
| NIST P-256 Field Inversion | `2^256-2^224+2^192+2^96-1` | 3 | 266 | 266 | +0 |
| NIST P-384 Field Inversion | `2^384-2^128-2^96+2^32-1` | 3 | 397 | 396 | +1 |
| secp256k1 (Bitcoin) Field Inversion | _too long_ | 3 | 269 | 269 | +0 |
| Curve25519 Scalar Inversion | _too long_ | 2 | 283 | 284 | -1 |
| NIST P-256 Scalar Inversion | _too long_ | 2 | 294 | 292 | +2 |
| NIST P-384 Scalar Inversion | _too long_ | 2 | 434 | 433 | +1 |
| secp256k1 (Bitcoin) Scalar Inversion | _too long_ | 2 | 293 | 290 | +3 |

