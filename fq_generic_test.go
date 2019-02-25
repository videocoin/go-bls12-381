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
		t.Run(fmt.Sprintf("a: %s, b: %s\n", testCase.a.String(), testCase.b.String()), func(t *testing.T) {})
	}
}

func TestFqMul(t *testing.T) {
	testCases := []struct {
		a, b, output fq
	}{}
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
