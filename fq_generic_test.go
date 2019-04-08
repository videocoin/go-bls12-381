package bls12

import (
	"fmt"
	"testing"
)

func TestFqAdd(t *testing.T) {
	testCases := []struct {
		a, b, output fq
	}{
		{a: fq0, b: fq0, output: fq0},
		{a: fq0, b: fq1, output: fq1},
		{a: fq0, b: fqLastElement, output: fqLastElement},
		{a: fqLastElement, b: fq1, output: fq0},
		{a: fqLastElement, b: fq100, output: fq99},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fq
			fqAdd(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqNeg(t *testing.T) {
	testCases := []struct {
		input, output fq
	}{
		{input: fqLastElement, output: fq1},
		{input: fq1, output: fqLastElement},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			var result fq
			fqNeg(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqSub(t *testing.T) {
	// TODO complete
	two, _ := fqMontgomeryFromBase10("2")

	negFqMont1 := new(fq)
	fqNeg(negFqMont1, &fqMont1)

	testCases := []struct {
		a, b, output fq
	}{
		{a: two, b: fqMont1, output: fqMont1},
		{a: fqLastElement, b: fqLastElement, output: fq0},
		{a: fq0, b: fqLastElement, output: fq1},
		{a: fq100, b: fq99, output: fq1},
		{a: fqMont1, b: two, output: *negFqMont1},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fq
			fqSub(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqBasicMul(t *testing.T) {
	testCases := []struct {
		a, b   fq
		output fqLarge
	}{
		{
			a:      fq{0, 0, 0, 0, 1},
			b:      fq{0, 0, 0, 0, 1},
			output: fqLarge{0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			a:      fq0,
			b:      fq1,
			output: fqLarge{0},
		},
		{
			a:      fq{0, 0, 0, 0, 1},
			b:      fq1,
			output: fqLarge{0, 0, 0, 0, 1},
		},
		{
			a:      fq1,
			b:      r2,
			output: fqLarge{r2[0], r2[1], r2[2], r2[3], r2[4], r2[5]},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fqLarge
			fqBasicMul(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}
func TestFqREDC(t *testing.T) {
	testCases := []struct {
		input  fqLarge
		output fq
	}{
		{
			input:  fqLarge{r2[0], r2[1], r2[2], r2[3], r2[4], r2[5]},
			output: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s\n", testCase.input.String()), func(t *testing.T) {
			var result fq
			fqREDC(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqMul(t *testing.T) {
	testCases := []struct {
		a, b, output fq
	}{
		{
			a:      fqMont1,
			b:      fqMont1,
			output: fqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result fq
			fqMul(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqExp(t *testing.T) {
	testCases := []struct {
		base, output fq
		exponent     []uint64
	}{
		{
			base:     fqMont1,
			exponent: []uint64{3},
			output:   fqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("base: %s, exponent: %v", testCase.base.String(), testCase.exponent), func(t *testing.T) {
			var result fq
			fqExp(&result, &testCase.base, testCase.exponent)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqCube(t *testing.T) {
	testCases := []struct {
		input, output fq
	}{
		{
			input:  fqMont1,
			output: fqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			var result fq
			fqCube(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqSqrt(t *testing.T) {
	testCases := []struct {
		input, output fq
	}{
		{
			input:  fqMont1,
			output: fqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			result := new(fq)
			fqSqrt(result, &testCase.input)
			if *result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqInv(t *testing.T) {
	negFqMont1 := new(fq)
	fqNeg(negFqMont1, &fqMont1)

	// TODO complete
	testCases := []struct {
		base     fq
		exponent []uint64
		output   fq
	}{
		{
			base:     *negFqMont1,
			exponent: []uint64{123123},
			output:   *negFqMont1,
		},
		{
			base:     fqMont1,
			exponent: []uint64{123123},
			output:   fqMont1,
		},
	}
	for _, testCase := range testCases {
		result := new(fq)
		fqExp(result, &testCase.base, testCase.exponent)
		if *result != testCase.output {
			t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
		}
	}
}
