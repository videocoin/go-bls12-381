package bls12

import (
	"fmt"
	"testing"
)

func TestFqMod(t *testing.T) {
	qPlus1 := [6]uint64{0xB9FEFFFFFFFFAAAC, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}

	// note: carry is always 0 for fq elements
	testCases := []struct {
		value, output fq
		carry         uint64
	}{
		{value: fqLastElement, carry: 0, output: fqLastElement},
		{value: fq(qU64), carry: 0, output: fq0},
		{value: fq(qPlus1), carry: 0, output: fq1},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("value: %s, carry: %d\n", testCase.value.String(), testCase.carry), func(t *testing.T) {
			if fqMod(&testCase.value, testCase.carry); testCase.value != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), testCase.value.String())
			}
		})
	}
}

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
