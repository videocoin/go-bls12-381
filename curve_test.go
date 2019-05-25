package bls12

import (
	"math/big"
	"testing"
)

func TestCurvePointSet(t *testing.T) {
	tests := map[string]struct {
		input, want curvePoint
	}{
		"zero":     {input: curvePoint{}, want: curvePoint{}},
		"non-zero": {input: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}}, want: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(curvePoint).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointEqual(t *testing.T) {
	tests := map[string]struct {
		a, b curvePoint
		want bool
	}{
		"zero": {
			a:    curvePoint{},
			b:    curvePoint{},
			want: true,
		},
		"tc 2": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			want: true,
		},
		"tc 3": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{2}, y: fq{2}, z: fq{3}},
			want: false,
		},
		"tc 4": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{3}, z: fq{3}},
			want: false,
		},
		"tc 5": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{2}, z: fq{1}},
			want: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(&tc.b)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointIsInfinity(t *testing.T) {
	tests := map[string]struct {
		input curvePoint
		want  bool
	}{
		"zero": {
			input: curvePoint{},
			want:  true,
		},
		"z = 0": {
			input: curvePoint{x: fq{1}, y: fq{2}, z: fq{}},
			want:  true,
		},

		"z != 0": {
			input: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			want:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.IsInfinity()
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointAdd(t *testing.T) {
	// TODO
}

func TestCurvePointDouble(t *testing.T) {
	// TODO
}

func TestCurvePointScalarMult(t *testing.T) {
	tests := map[string]struct {
		point  curvePoint
		scalar *big.Int
		want   curvePoint
	}{
		"testcase 1": {
			point:  g1Gen.p,
			scalar: new(big.Int).SetUint64(57),
			want: curvePoint{
				x: fq{10167668087098498840, 12007495882678423899, 3444586411392841405, 17496946105194000634, 15545891604672389846, 1004146385472758665},
				y: fq{17399866279200227458, 13488772900345834999, 9150629517599614720, 12422932654397576046, 10596503356839999349, 1360682455110082503},
				z: fq{16288557835584950008, 11072915878001992196, 1009666109118314550, 14317758622910970758, 9177236567440559526, 1669672164031368809},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(curvePoint).ScalarMult(&tc.point, tc.scalar)
			if !tc.want.Equal(got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

/*
func TestCurvePointToAffine(t *testing.T) {
	tests := map[string]struct {
		input, want curvePoint
	}{
		"gen + gen": {

			input: curvePoint{
				x: fq{1170480656948039723, 15860523744001223122, 3723658602231154259, 16074848579805679309, 16321276149294171578, 928717658087505831},
				y: fq{1230614714935444159, 12163446147838635782, 1397133617996746430, 4234108216169586935, 7101055810179665440, 1568086071381718864},
				z: fq{8455833386895688930, 1748740486030555933, 13453024110247299997, 11770351495059383081, 2033683641984398208, 1691240166868468948},
			},
			want: curvePoint{
				x: fq{6046496802367715900, 4512703842675942905, 5557647857818872160, 11911007586355426777, 2789226406901363231, 2402832991291269},
				y: fq{8075247918781118784, 15723127573743364860, 13289805640942397317, 12593984073093990549, 2724610382811436832, 447576566110657301},
				z: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.ToAffine()
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
*/

/*
TODO REVIEW RAND
func TestCurvePointMarshalUnmarshal(t *testing.T) {
	p1 := &curvePoint{
		x: *new(fq).Rand(),
		y: *new(fq).Rand(),
		z: *new(fq).SetUint64(1),
	}
	p2 := new(curvePoint)
	p2.Unmarshal(p1.Marshal())
	if !p2.Equal(p1) {
		t.Errorf("Marshaling/unmarshaling failed: expected %v, got %v", p1, p2)
	}
}
*/

func TestCurvePointSetBytes(t *testing.T) {
	// TODO
}

func TestCurvePointSWEncode(t *testing.T) {
	// TODO
}
