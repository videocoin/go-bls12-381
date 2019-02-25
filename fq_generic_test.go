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
	testCases := []struct {
		a, b, output fq
	}{
		{a: fqLastElement, b: fqLastElement, output: fq0},
		{a: fq0, b: fqLastElement, output: fq1},
		{a: fq100, b: fq99, output: fq1},
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
			a:      fq{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			b:      fq{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			output: fqLarge{0x0000000000000001, 0x0000000000000001, 0x0000000000000001, 0x0000000000000001, 0x0000000000000001, 0000000000000001, 0000000000000001, 0000000000000001, 0000000000000001},
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

/*
func TestFqMul(t *testing.T) {
	testCases := []struct {
		a, b, output fq
	}{
		{
			a:      fq{0xcc6200000020aa8a, 0x422800801dd8001a, 0x7f4f5e619041c62c, 0x8a55171ac70ed2ba, 0x3f69cc3a3d07d58b, 0xb972455fd09b8ef},
			b:      fq{0x329300000030ffcf, 0x633c00c02cc40028, 0xbef70d925862a942, 0x4f7fa2a82a963c17, 0xdf1eb2575b8bc051, 0x1162b680fb8e9566},
			output: fq{0x9dc4000001ebfe14, 0x2850078997b00193, 0xa8197f1abb4d7bf, 0xc0309573f4bfe871, 0xf48d0923ffaf7620, 0x11d4b58c7a926e66},
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
*/
