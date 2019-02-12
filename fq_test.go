package bls12

import (
	"fmt"
	"math/big"
	"testing"
)

func TestNewFq(t *testing.T) {
	testCases := []struct {
		input  *big.Int
		output string
	}{
		{
			input:  bigFromBase10("1"),
			output: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			input:  bigFromBase10("0"),
			output: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			input:  bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"),
			output: "1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab",
		},
		{
			input:  bigFromBase10("40312312312302409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"),
			output: "c365b13802f37e3028339cc48d630919b974b3740485f0a23abf8f8762dc16666199fffeb153ffffb9feffffffffaaab",
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("input: %s", testCase.input), func(t *testing.T) {
			if fq := newFq(testCase.input); fq.String() != testCase.output {
				t.Errorf("expected %s, got %s", testCase.output, fq)
			}
		})
	}
}
