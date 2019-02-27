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

func TestFq2Mul(t *testing.T) {}
