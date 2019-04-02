package bls12

import (
	"fmt"
	"testing"
)

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
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fq2
			fq2Add(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
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
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fq2
			fq2Sub(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFq2String(t *testing.T) {
	testCases := []struct {
		input  fq2
		output string
	}{
		/*
			{fq0, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
			{fq1, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"},
			{fqLastElement, "1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaaa"},
		*/
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("fq2: %v\n", testCase.input), func(t *testing.T) {
			if str := testCase.input.String(); str != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output, str)
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
				c0: Fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
				c1: Fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
			},
			y: fq2{
				c0: Fq{5508758831087832138, 6448303779119275098, 16710190169160573786, 13542242618704742751, 563980702369916322, 37152010398653157},
				c1: Fq{12520284671833321565, 1777275927576994268, 9704602344324656032, 8739618045342622522, 16651875250601773805, 804950956836789234},
			},
			output: fq2{
				c0: Fq{15827099533983733295, 5585758986002856792, 17089326063188852987, 9744188138385726684, 16136916081236286470, 539400444303524798},
				c1: Fq{1625211833195412533, 6260280147765452499, 3544356884998474521, 9883248032979034166, 8703231040866562622, 339556463953896598},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("A: %v, B: %v\n", testCase.x, testCase.y), func(t *testing.T) {
			result := new(fq2)
			fq2Mul(result, &testCase.x, &testCase.y)
			if *result != testCase.output {
				t.Errorf("expected %v, got %v\n", testCase.output, *result)
			}
		})
	}
}
