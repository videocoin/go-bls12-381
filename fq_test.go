package bls12

import (
	"fmt"
	"math/big"
	"testing"
)

var (
	fqLastElement = fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
	fq100         = fq{100}
	fq99          = fq{99}
)

func TestFqFromBig(t *testing.T) {
	testCases := []struct {
		input       *big.Int
		output      fq
		expectedErr error
	}{
		{input: bigQ, output: fq{}, expectedErr: errOutOfBounds},
		{input: big0, output: fq0, expectedErr: nil},
		{input: big1, output: fq1, expectedErr: nil},
		{input: new(big.Int).Sub(bigQ, big1), output: fqLastElement, expectedErr: nil},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("big integer: %s\n", testCase.input), func(t *testing.T) {
			result, err := fqFromBig(testCase.input)
			if err != nil {
				if err != testCase.expectedErr {
					t.Errorf("expected %s, got %s\n", testCase.expectedErr, err)
				}
				return
			}
			if result != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
			}
		})
	}
}

func TestFqMontgomeryFromBig(t *testing.T) {

}

func TestFqHex(t *testing.T) {
	testCases := []struct {
		input  fq
		output string
	}{
		{fq0, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{fq1, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"},
		{fqLastElement, "1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaaa"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("fq: %v\n", testCase.input), func(t *testing.T) {
			if hex := testCase.input.Hex(); hex != testCase.output {
				t.Errorf("expected %s, got %s\n", testCase.output, hex)
			}
		})
	}
}
func TestFqEqual(t *testing.T) {
	testCases := []struct {
		a, b   fq
		output bool
	}{
		{fq0, fq0, true},
		{fq0, fq{0}, true},
		{fq1, fq{1}, true},
		{fq0, fq1, false},
		{fq{0x0, 0x1, 0x2, 0x4, 0x5, 0x5}, fq{0x0, 0x1, 0x2, 0x4, 0x5, 0x5}, true},
		{fq{0x5, 0x4, 0x3, 0x2, 0x1, 0x0}, fq{0x5, 0x4, 0x3, 0x2, 0x1, 0x0}, true},
		{fq{0x5, 0x4, 0x3, 0x2, 0x1, 0x0}, fq{0x5, 0x4, 0x3, 0x2, 0x1}, true},
		{fq{0x5, 0x4, 0x3, 0x2, 0x1, 0x3}, fq{0x5, 0x4, 0x3, 0x2, 0x1}, false},
		{fq{0x5, 0x4, 0x3, 0x2, 0x1, 0x0}, fq{0x6, 0x4, 0x3, 0x2, 0x1, 0x0}, false},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {
			if result := testCase.a.equal(testCase.b); result != testCase.output {
				t.Errorf("expected %t, got %t\n", testCase.output, result)
			}
		})
	}

}
