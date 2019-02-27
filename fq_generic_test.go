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
			output: fqLarge{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
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
func TestFqREDC(t *testing.T) {
	testCases := []struct {
		input  fqLarge
		output fq
	}{
		{
			input:  fqLarge{0x00f630eca1bd6630, 0xfd5552c7c9513c59, 0x4bd35517a56fa1a5, 0xda0eea7dbbcbb3f2, 0x6ab765fc88b83da5, 0xd5bb7a80c6cbb4a7, 0xb8ace446959099f0, 0x873087f81d8beac0, 0xcd85184272c3ee73, 0x27d2462594a0a1e5, 0x2ad107a9544cc968, 0x075cfe63e862ac9b},
			output: fq{0xf5f70713b717914c, 0x355ea5ac64cbbab1, 0xce60dd43417ec960, 0xf16b9d77b0ad7d10, 0xa44c204c1de7cdb7, 0x1684487772bc9a5a},
		},
	}
	// bls12: 144fb1019f8e91b25f7b8ffde88a30adf33ca1b49e03a60beb66a8c54e38c359bdb913d49fb352404e450d3690ca3e8a

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

/*
OUR:
11159732349363015024
13402431016077863595
69296637154780720
821004262248500634

11159732349363015024
2210141511517208575
18254587682645687385
11159732349363015024
2210141511517208575
9718722949989020304
9067282531769192061
*/

/*
BLS12NATIVE:
11159732349363015024
13402431016077863595
69296637154780720
8108072751082139500

11159732349363015024
2210141511517208575
18254587682645687385
11159732349363015024
2210141511517208575
17634639310007295573
1337069979623170057
*/
