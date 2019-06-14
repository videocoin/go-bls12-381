package bls12

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestDecomposeGLV(t *testing.T) {
	// See Guide to Pairing-Based Cryptography - 6.3.2
	lambda := bigFromBase10("228988810152649578064853576960394133503")

	// n ≡ n0 + n1λ mod r
	n, _ := rand.Int(rand.Reader, r)
	subScalars := glvLattice.Decompose(n)
	fmt.Println(subScalars)

	sum, term := new(big.Int), new(big.Int)
	sum.Add(sum, subScalars[0])
	sum.Add(sum, term.Mul(subScalars[1], lambda))

	mod := new(big.Int)
	new(big.Int).DivMod(sum, r, mod)
	if mod.Cmp(n) != 0 {
		t.Fatalf("expected: %v, got: %v", n, sum)
	}
}

func TestDecomposeGLS2(t *testing.T) {
	n := bigFromBase10("16352079474354410340595364877818663142712328869634074743858828567442853094216")
	subScalars := glsLattice.Decompose(n)
	fmt.Println(subScalars)

	/*
		1- [7551454145663700495 343877308400216775 24371303287397002030 -4719017615457912682]

		2-	[7551454145663700494 142819899919596172490354505392508578502 -199108166460466086 4719017615457912682]

		7551454145663700494
			32704158948708820681111985699948264916707782120674977914261840309922543752506
			3012979683939863664635314476423708672
			-16352079474354410340516620822129601773927056280765044947830193197059723886592
			--- FAIL: TestDecomposeGLS (0.00s)
				/Users/ricardogeraldes/code/src/github.com/VideoCoin/go-bls12-381/lattice_test.go:56:
				expected: 16352079474354410340595364877818663142712328869634074743858828567442853094216,
				got: 16352079474354410340595364877818663142783738819593872830103833881484907275080
				16352079474354410340595364877818663142640918919674276657613823253400798913352
	*/
}

func TestDecomposeGLS(t *testing.T) {
	// See Guide to Pairing-Based Cryptography - 6.3.2
	x := bigFromBase10("-15132376222941642752")
	lambda1 := new(big.Int).Mul(x, x)
	lambda1.Sub(lambda1, bigOne)
	lambda2 := new(big.Int).Set(x)
	lambda3 := new(big.Int).Mul(lambda1, lambda2)

	// n ≡ n0 + n1λ1 + n2λ2 + n3λ1λ2 mod r
	//n, _ := rand.Int(rand.Reader, r)
	n := bigFromBase10("16352079474354410340595364877818663142712328869634074743858828567442853094216")
	subScalars := glsLattice.Decompose(n)

	sum := new(big.Int)
	sum.Add(sum, subScalars[0])
	sum.Add(sum, new(big.Int).Mul(subScalars[1], lambda1))
	sum.Add(sum, new(big.Int).Mul(subScalars[2], lambda2))
	sum.Add(sum, new(big.Int).Mul(subScalars[3], lambda3))

	mod := new(big.Int)
	new(big.Int).DivMod(sum, r, mod)
	if mod.Cmp(n) != 0 {
		t.Fatalf("expected: %v, got: %v", n, mod)
	}
}

func TestOther(t *testing.T) {
	/*
		subScalars := []*big.Int{
			bigFromBase10("-2796572967770866043"),
			bigFromBase10("2036585984363498129"),
			bigFromBase10("27986574728709285892"),
			bigFromBase10("4787294474459546991"),
		}*/

	sum := new(big.Int).Add(bigFromBase10("-2796572967770866043"), bigFromBase10("466355401332960035202818963869017731292095396582916715887"))
	sum.Add(sum, bigFromBase10("423503377986299853829296325399297654784"))
	sum.Add(sum, bigFromBase10("-16588668679064483740876286834069811886906678135563818078618369863374901149696"))

	mod := new(big.Int)
	new(big.Int).DivMod(sum, r, mod)
	fmt.Println(mod)

	/*
				35847206496061706738105098272783193915157552023108502158861100541613694798103
				35847206496061706739037809075449113986410196706819137329103883985578123539445

		 35847206496061706739037809075449113985490746809769006065114226805236249670645,
		 got: 35847206496061706739037809075449113985563189950846537621445291334779528229877
	*/

}

func TestOther2(t *testing.T) {
	adj := []*big.Int{
		bigFromBase10("3465144826073652318776269530687742778270252468765361963008"),
		bigFromBase10("-228988810152649578064853576960394133503"),
		bigFromBase10("15132376222941642752"),
		bigFromBase10("1"),
	}

	for _, ai := range adj {
		value, mod := new(big.Int), new(big.Int)
		value.DivMod(new(big.Int).Mul(ai, v1), r, mod)
		round(value, mod)
		fmt.Println(value)
	}
}
