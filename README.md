## go-bls12-381

Package go-bls12-381 implements the BLS12 pairing-friendly elliptic curve construction for u = -0xd201000000010000.

### Benchmarks

(2,3 GHz Intel Core i7)

branch `master`:

```
BenchmarkG1-8
BenchmarkG2-8
BenchmarkPairing-8
```

branch `master` -tags generic:

```
BenchmarkG1-8        	    1000	   2049394 ns/op
BenchmarkG2-8        	     200	   6370759 ns/op
BenchmarkPairing-8   	     100	  15894423 ns/op
```

branch `lattices`:

```
BenchmarkG1-8
BenchmarkG2-8
BenchmarkPairing-8
```
