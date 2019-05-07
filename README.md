## go-bls12-381

Package go-bls12-381 implements the BLS12 pairing-friendly elliptic curve construction for u = -0xd201000000010000.

### Benchmarks

(2,3 GHz Intel Core i7)

`cloudflare/bn256` branch `master`:

```
BenchmarkG1-8 5000 232965 ns/op
BenchmarkG2-8 2000 796422 ns/op
BenchmarkGT-8 1000 2052766 ns/op
BenchmarkPairing-8 500 2561803 ns/op
```
