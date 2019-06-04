## go-bls12-381

Package go-bls12-381 implements the BLS12 pairing-friendly elliptic curve construction for u = -0xd201000000010000.

Test vectors taken from [Relic](https://github.com/relic-toolkit/relic).

### Benchmarks

(2,3 GHz Intel Core i7)

branch `master`:

```
BenchmarkG1-8        	    3000	    574851 ns/op
BenchmarkG2-8        	    1000	   2080695 ns/op
BenchmarkPairing-8   	     300	   5481567 ns/op
```

branch `lattices`:

```
BenchmarkG1-8
BenchmarkG2-8
BenchmarkPairing-8
```
