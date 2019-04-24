package bls12

import (
	"fmt"
	"testing"
)

func TestFq2SetZero(t *testing.T) {
	// TODO
}

func TestFq2SetOne(t *testing.T) {
	// TODO
}

func TestFq2Set(t *testing.T) {
	// TODO
}

func TestFq2Neg(t *testing.T) {
	// TODO
}

func TestFq2Add(t *testing.T) {
	testCases := []struct {
		a, b, output fq2
	}{
		{
			a: fq2{
				c0: fqLastElement,
				c1: fq100,
			},
			b: fq2{
				c0: fq100,
				c1: fqLastElement,
			},
			output: fq2{
				c0: fq99,
				c1: fq99,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %v, b: %v\n", testCase.a, testCase.b), func(t *testing.T) {
			if result := new(fq2).Add(&testCase.a, &testCase.b); *result != testCase.output {
				t.Errorf("expected %v, got %v\n", testCase.output, result)
			}
		})
	}
}

func TestFq2Sub(t *testing.T) {
	testCases := []struct {
		a, b, output fq2
	}{
		{
			a: fq2{
				c0: fq100,
				c1: fq99,
			},
			b: fq2{
				c0: fq99,
				c1: fq100,
			},
			output: fq2{
				c0: fq1,
				c1: fqLastElement,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %v, b: %v\n", testCase.a, testCase.b), func(t *testing.T) {
			if result := new(fq2).Sub(&testCase.a, &testCase.b); *result != testCase.output {
				t.Errorf("expected %v, got %v\n", testCase.output, result)
			}
		})
	}
}

func TestFq2Mul(t *testing.T) {
	testCases := []struct {
		x, y, output fq2
	}{
		{
			x: fq2{
				c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
				c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
			},
			y: fq2{
				c0: fq{5508758831087832138, 6448303779119275098, 16710190169160573786, 13542242618704742751, 563980702369916322, 37152010398653157},
				c1: fq{12520284671833321565, 1777275927576994268, 9704602344324656032, 8739618045342622522, 16651875250601773805, 804950956836789234},
			},
			output: fq2{
				c0: fq{15827099533983733295, 5585758986002856792, 17089326063188852987, 9744188138385726684, 16136916081236286470, 539400444303524798},
				c1: fq{1625211833195412533, 6260280147765452499, 3544356884998474521, 9883248032979034166, 8703231040866562622, 339556463953896598},
			},
		},
		{
			x: fq2{
				c0: fq{2148348102093263616, 12232197882281708926, 11330363351339265390, 8919790901940522406, 17524282994943615806, 496043075549758450},
				c1: fq{406374380327561821, 16222300308001049590, 2744191801523148582, 9384378502465456106, 6088477103101489105, 1478173219842328727},
			},
			y: fq2{
				c0: fq{3656661975483800321, 2343032803162206357, 3148888360041182374, 13718521627034877454, 7582705438920995243, 1502510331453350514},
				c1: fq{12908108411122378854, 12884292956005326590, 14415309872191119313, 17916475641229970923, 8768506067738497560, 1577813823053251215},
			},
			output: fq2{
				c0: fq{16633317224619977984, 8114711876341744521, 5129895554745463763, 3028104440572865277, 600965037607129992, 1357055630732818745},
				c1: fq{16769752770056333816, 9539351424971908300, 1361078859430182954, 6919544897282332375, 8797177321797491027, 1057830777381385737},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("A: %v, B: %v\n", testCase.x, testCase.y), func(t *testing.T) {
			if result := new(fq2).Mul(&testCase.x, &testCase.y); *result != testCase.output {
				t.Errorf("expected %v, got %v\n", testCase.output, *result)
			}
		})
	}
}

func TestFq2MulXi(t *testing.T) {
	// TODO
}

func TestFq2Sqr(t *testing.T) {
	// TODO
}
