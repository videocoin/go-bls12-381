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

func TestFqMul(t *testing.T) {

}
