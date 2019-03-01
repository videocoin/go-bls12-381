package bls12

import (
	"fmt"
	"math/big"
	"testing"
)

var (
	big0 = new(big.Int).SetUint64(0)
	big1 = new(big.Int).SetUint64(1)

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
		{input: q, output: fq{}, expectedErr: errOutOfBounds},
		{input: big0, output: fq0, expectedErr: nil},
		{input: big1, output: fq1, expectedErr: nil},
		{input: new(big.Int).Sub(q, big1), output: fqLastElement, expectedErr: nil},
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

func TestFqFromBase10(t *testing.T) {
	testCases := []struct {
		input       string
		output      fq
		expectedErr error
	}{
		{input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", output: fq{}, expectedErr: errOutOfBounds},
		{input: "0", output: fq0, expectedErr: nil},
		{input: "1", output: fq1, expectedErr: nil},
		{input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", output: fqLastElement, expectedErr: nil},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("big integer: %s\n", testCase.input), func(t *testing.T) {
			result, err := fqFromBase10(testCase.input)
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

func TestFqMontgomeryFromBig(t *testing.T)    {}
func TestFqMontgomeryFromBase10(t *testing.T) {}

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

func TestFqLargeHex(t *testing.T) {}

func TestMontEncode(t *testing.T) {
	testCases := []struct {
		input, output fq
	}{
		{
			input:  fq1,
			output: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
	}
	for _, testCase := range testCases {
		var result fq
		montgomeryEncode(&result, &testCase.input)
		if result != testCase.output {
			t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
		}
	}
}

func TestMontDecode(t *testing.T) {
	testCases := []struct {
		input, output fq
	}{
		{
			input:  fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			output: fq1,
		},
	}
	for _, testCase := range testCases {
		var result fq
		montgomeryDecode(&result, &testCase.input)
		if result != testCase.output {
			t.Errorf("expected %s, got %s\n", testCase.output.String(), result.String())
		}
	}
}
