package bls12

import (
	"fmt"
	"testing"
)

func TestFqAdd(t *testing.T) {
	testCases := []struct {
		a, b, output Fq
	}{
		{a: Fq0, b: Fq0, output: Fq0},
		{a: Fq0, b: fq1, output: fq1},
		{a: Fq0, b: fqLastElement, output: fqLastElement},
		{a: fqLastElement, b: fq1, output: Fq0},
		{a: fqLastElement, b: fq100, output: fq99},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result Fq
			FqAdd(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqNeg(t *testing.T) {
	testCases := []struct {
		input, output Fq
	}{
		{input: fqLastElement, output: fq1},
		{input: fq1, output: fqLastElement},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			var result Fq
			FqNeg(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqSub(t *testing.T) {
	testCases := []struct {
		a, b, output Fq
	}{
		{a: fqLastElement, b: fqLastElement, output: Fq0},
		{a: Fq0, b: fqLastElement, output: fq1},
		{a: fq100, b: fq99, output: fq1},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result Fq
			FqSub(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqBasicMul(t *testing.T) {
	testCases := []struct {
		a, b   Fq
		output FqLarge
	}{
		{
			a:      Fq{0, 0, 0, 0, 1},
			b:      Fq{0, 0, 0, 0, 1},
			output: FqLarge{0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			a:      Fq0,
			b:      fq1,
			output: FqLarge{0},
		},
		{
			a:      Fq{0, 0, 0, 0, 1},
			b:      fq1,
			output: FqLarge{0, 0, 0, 0, 1},
		},
		{
			a:      fq1,
			b:      r2,
			output: FqLarge{r2[0], r2[1], r2[2], r2[3], r2[4], r2[5]},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result FqLarge
			FqBasicMul(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}
func TestFqREDC(t *testing.T) {
	testCases := []struct {
		input  FqLarge
		output Fq
	}{
		{
			input:  FqLarge{r2[0], r2[1], r2[2], r2[3], r2[4], r2[5]},
			output: Fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s\n", testCase.input.String()), func(t *testing.T) {
			var result Fq
			FqREDC(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqMul(t *testing.T) {
	testCases := []struct {
		a, b, output Fq
	}{
		{
			a:      FqMont1,
			b:      FqMont1,
			output: FqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			var result Fq
			FqMul(&result, &testCase.a, &testCase.b)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqExp(t *testing.T) {
	testCases := []struct {
		base, output Fq
		exponent     []uint64
	}{
		{
			base:     FqMont1,
			exponent: []uint64{3},
			output:   FqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("base: %s, exponent: %v", testCase.base.String(), testCase.exponent), func(t *testing.T) {
			var result Fq
			FqExp(&result, &testCase.base, testCase.exponent)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqCube(t *testing.T) {
	testCases := []struct {
		input, output Fq
	}{
		{
			input:  FqMont1,
			output: FqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			var result Fq
			FqCube(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqSqrt(t *testing.T) {
	testCases := []struct {
		input, output Fq
	}{
		{
			input:  FqMont1,
			output: FqMont1,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input.String()), func(t *testing.T) {
			var result Fq
			FqSqrt(&result, &testCase.input)
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}
