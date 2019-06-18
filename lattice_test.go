package bls12

/*
func TestGLV2(t *testing.T) {
	// See Guide to Pairing-Based Cryptography - 6.3.2
	lambda := bigFromBase10("228988810152649578064853576960394133503")

	// n ≡ n0 + n1λ mod r
	n, _ := rand.Int(rand.Reader, r)
	subScalars := glvLattice.Decompose(n)

	sum, term := new(big.Int), new(big.Int)
	sum.Add(sum, subScalars[0])
	sum.Add(sum, term.Mul(subScalars[1], lambda))

	mod := new(big.Int)
	new(big.Int).DivMod(sum, r, mod)
	if mod.Cmp(n) != 0 {
		t.Fatalf("expected: %v, got: %v", n, mod)
	}
}

func TestGLVEndomorphism(t *testing.T) {
	p0 := g1Gen.p
	lambda := bigFromBase10("228988810152649578064853576960394133503")
	p11 := new(curvePoint).ScalarMult(&p0, lambda)
	p12 := p0
	fqMul(&p12.x, &p12.x, &frobFq6C2[2].c0)
	if !p11.ToAffine().Equal(p12.ToAffine()) {
		t.Fatalf("expected: %v, got: %v", p11, p12)
	}
}

func TestGLS4(t *testing.T) {
	x := bigFromBase10("15132376222941642752")
	lambda1 := x
	lambda2 := new(big.Int).Mul(x, x)
	lambda3 := new(big.Int).Mul(lambda1, lambda2)

	//n, _ := randFieldElement(rand.Reader)
	n := bigFromBase10("10123123123123")
	subScalars := glsLattice.Decompose(n)

	for _, si := range subScalars {
		if si.Sign() == -1 {
			si.Neg(si)
		}
	}

	sum, term := new(big.Int), new(big.Int)
	sum.Add(sum, subScalars[0])
	sum.Add(sum, term.Mul(subScalars[1], lambda1))
	sum.Add(sum, term.Mul(subScalars[2], lambda2))
	sum.Add(sum, term.Mul(subScalars[3], lambda3))

	mod := new(big.Int)
	new(big.Int).DivMod(sum, r, mod)
	if mod.Cmp(n) != 0 {
		t.Fatalf("expected: %v, got: %v", n, mod)
	}
}

func TestGLSEndomorphism(t *testing.T) {
	/*
		//54.... {[8505329371266088957 17002214543764226050 6865905132761471162 8632934651105793861 6631298214892334189 1582556514881692819] [0 0 0 0 0 0]} {[8505329371266088957 17002214543764226050 6865905132761471162 8632934651105793861 6631298214892334189 1582556514881692819] [0 0 0 0 0 0]}},
		qx := fq2{
			c0: fq{7945946479894423835, 15934652502800266572, 2488440080670104002, 9386780446441512881, 11030534594722426706, 816011729840636128},
			c1: fq{13033736016918358455, 10575140707633074893, 712033873462394481, 15426744332523829345, 15149289313990657609, 527804130256077952},
		}
		qy := fq2{
			c0: fq{3979628422115981119, 8799464239825509261, 8769850388554038001, 14088617588560639466, 16028322177399626617, 1349288934766117873},
			c1: fq{4943015918193553641, 13632206254062479526, 686991123441432951, 12362413863511458327, 14482031900939104829, 591053520888768029},
		}
		p := g2Gen.p

		fmt.Println("Endomorphisms")
		fmt.Println()
		//fmt.Println(new(fq2).Mul(&qx, new(fq2).Inv(new(fq2).Conjugate(&p.x))))
		//fmt.Println(new(fq2).Mul(&qy, new(fq2).Inv(new(fq2).Conjugate(&p.y))))
		fmt.Println(new(fq2).Mul(&qx, new(fq2).Inv(&p.x)))
		fmt.Println(new(fq2).Mul(&qy, new(fq2).Inv(&p.y)))


			//fmt.Println(new(fq2).Mul(&qx, new(fq2).Inv(new(fq2).Conjugate(&p.x))))
			//fmt.Println(new(fq2).Mul(&qy, new(fq2).Inv(new(fq2).Conjugate(&p.y))))

				fmt.Println("got")
				fmt.Println()
				fmt.Println(new(fq2).Mul(new(fq2).Conjugate(&p.x), &fq2{c1: frobFq6C2[1].c0}))
				fmt.Println(new(fq2).Mul(new(fq2).Conjugate(&p.y), frobFq12C1[3]))
				t.Fatal()
*/

/*
	x := bigFromBase10("15132376222941642752")
	xFinal := new(big.Int).Sub(r, x)
	xFinal.Add(xFinal, r)
	mod := new(big.Int).Set(r)
	new(big.Int).DivMod(xFinal, r, mod)
	fmt.Println(mod)
	// [λ^i]P0 = Pi

*/

//p0 := g2Gen.p
//r := new(big.Int).SetUint64(15132376222941642751)
//fmt.Println(glsLattice.Decompose(r))
//fmt.Println(new(twistPoint).ScalarMult(&p0, r))
//lambda := bigFromBase10("52435875175126190479447740508185965837690552500527637822588526323715639541761")
//p11 := new(twistPoint).ScalarMult(&p0, lambda)
//p12 := glsEnd(&p0)
/*
	if !p11.ToAffine().Equal(p12.ToAffine()) {
		t.Fatalf("expected: %v, got: %v", p11, p12)
	}
*/

//lambda = bigFromBase10("228988810152649578064853576960394133504")
//p11 = new(twistPoint).ScalarMult(&p0, lambda)
//p12 = glsEnd(p12)
//fmt.Println(p12)
/*
	if !p11.ToAffine().Equal(p12.ToAffine()) {
		t.Fatalf("expected: %v, got: %v", p11, p12)
	}
*/

//lambda = bigFromBase10("52435875175126190475982595682112313518914282969839895044333406231173219221505")
//p11 = new(twistPoint).ScalarMult(&p0, lambda)
//p12 = glsEnd(p12)

/*
	if !p11.ToAffine().Equal(p12.ToAffine()) {
		t.Fatalf("expected: %v, got: %v", p11, p12)
	}
*/
//fmt.Println(p12)
//n, _ := randFieldElement(rand.Reader)
//p := new(twistPoint).ScalarMult(&g2Gen.p, n)
//fmt.Println(n)
//fmt.Println(p)

//t.Fatal()
//}
