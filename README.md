## go-bls12-381

Package go-bls12-381 implements the [BLS12-381](https://electriccoin.co/blog/new-snark-curve/) pairing-friendly elliptic curve construction that targets the 128-bit security level. This project used the [constant-time hashing to the BLS12-381 elliptic curve proposed by Wahby, Boneh](https://eprint.iacr.org/2019/403.pdf). For G1/G2 signatures, use sig1/sig2 respectively.

There is a `lattices` branch that implements the 2-GLV method on G1 and 4-GLS method on G2 - both methods use an efficent endomorphism and scalar decomposition to speed up elliptic curve scalar multiplication.

Test vectors taken from [Relic](https://github.com/relic-toolkit/relic).
Inspiration taken from Cloudflare's [bn256](https://github.com/cloudflare/bn256) implementation.

## Benchmarks

branch `master`:

```
BenchmarkG1-8        	    3000	    595602 ns/op
BenchmarkG2-8        	    1000	   2117188 ns/op
BenchmarkPairing-8   	     300	   5985261 ns/op
```

branch `lattices`:

```
BenchmarkG1-8        	    5000	    359453 ns/op
BenchmarkG2-8        	    1000	   1433320 ns/op
BenchmarkPairing-8   	     300	   5702674 ns/op
```
