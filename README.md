<p align="center">
  <img src="logo.svg" width="40%" border="0" alt="addchain" />
</p>

<p align="center">Cryptographic Addition Chain Generation in Go</p>

## Results

| Name | _N_ | _d_ | Length | Best | Delta |
| ---- | --- | --- | -----: | ---: | ----: |
| Curve25519 Field Inversion | `2^255-19` | 2 | 266 | 265 | +1 |
| NIST P-256 Field Inversion | `2^256-2^224+2^192+2^96-1` | 3 | 266 | 266 | +0 |
| NIST P-384 Field Inversion | `2^384-2^128-2^96+2^32-1` | 3 | 397 | 396 | +1 |
| secp256k1 (Bitcoin) Field Inversion | _too long_ | 3 | 269 | 269 | +0 |
| Curve25519 Scalar Inversion | _too long_ | 0 | 293 | 290 | +3 |
| NIST P-256 Scalar Inversion | _too long_ | 0 | 294 | 292 | +2 |
| NIST P-384 Scalar Inversion | _too long_ | 0 | 434 | 433 | +1 |
| secp256k1 (Bitcoin) Scalar Inversion | _too long_ | 0 | 283 | 284 | -1 |

