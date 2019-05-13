package bls12

import (
	"math/big"
	"testing"
)

var (
	fqOneHardcoded          = &fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}
	fqLastHardcoded         = &fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}
	fqLastStandardHardcoded = &fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
	fqTwoStandard           = &fq{2}
	bigOne                  = big.NewInt(1)
	bigZero                 = big.NewInt(0)
	bigLast, _              = bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786")
)

func TestFqIsOne(t *testing.T) {
	tests := map[string]struct {
		input *fq
		want  bool
	}{
		"non-one":   {input: fqZero, want: false},
		"non-one 2": {input: fqLastHardcoded, want: false},
		"one":       {input: fqOneHardcoded, want: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.IsOne()
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSet(t *testing.T) {
	tests := map[string]struct {
		input, want *fq
	}{
		"zero":     {input: fqZero, want: fqZero},
		"non-zero": {input: fqLastHardcoded, want: fqLastHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).Set(tc.input)
			if *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetZero(t *testing.T) {
	tests := map[string]struct {
		input *fq
	}{
		"zero":     {input: fqZero},
		"non-zero": {input: fqLastHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := *tc.input
			if *got.SetZero() != *fqZero {
				t.Fatalf("expected: %v, got: %v", fqZero, got)
			}
		})
	}
}

func TestFqSetOne(t *testing.T) {
	tests := map[string]struct {
		input *fq
	}{
		"zero":     {input: fqZero},
		"one":      {input: fqOneHardcoded},
		"non-zero": {input: fqLastHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := *tc.input
			if *got.SetOne() != *fqOneHardcoded {
				t.Fatalf("expected: %v, got: %v", fqOneHardcoded, got)
			}
		})
	}
}

func TestFqSetString(t *testing.T) {
	tests := map[string]struct {
		input string
		want  *fq
		err   error
	}{
		"zero":               {input: "0", want: fqZero, err: nil},
		"one":                {input: "1", want: fqOneHardcoded, err: nil},
		"last field element": {input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", want: fqLastHardcoded, err: nil},
		"out of bounds":      {input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", err: errOutOfBounds},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := new(fq).SetString(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if err == nil && *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetInt(t *testing.T) {
	tests := map[string]struct {
		input *big.Int
		want  *fq
		err   error
	}{
		"zero":               {input: bigZero, want: fqZero, err: nil},
		"one":                {input: bigOne, want: fqOneHardcoded, err: nil},
		"last field element": {input: new(big.Int).Sub(q, bigOne), want: fqLastHardcoded, err: nil},
		"out of bounds":      {input: q, err: errOutOfBounds},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := new(fq).SetInt(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if err == nil && *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetUint64(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  *fq
	}{
		"zero":     {input: 0, want: fqZero},
		"non-zero": {input: 1, want: fqOne},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).SetUint64(tc.input)
			if *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", fqZero, got)
			}
		})
	}
}

func TestFqMontgomeryEncode(t *testing.T) {
	tests := map[string]struct {
		input, want *fq
	}{
		"zero":     {input: fqZero, want: fqZero},
		"non-zero": {input: fqLastStandardHardcoded, want: fqLastHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).MontgomeryEncode(tc.input)
			if *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqMontgomeryDecode(t *testing.T) {
	tests := map[string]struct {
		input, want *fq
	}{
		"zero":     {input: fqZero, want: fqZero},
		"non-zero": {input: fqLastHardcoded, want: fqLastStandardHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).MontgomeryDecode(tc.input)
			if *got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqEqual(t *testing.T) {
	tests := map[string]struct {
		x, y *fq
		want bool
	}{
		"equal":     {x: fqZero, y: fqZero, want: true},
		"equal 2":   {x: fqZero, y: &fq{}, want: true},
		"equal 3":   {x: fqOneStandard, y: &fq{1}, want: true},
		"not equal": {x: fqOneHardcoded, y: fqLastHardcoded, want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.x.Equal(tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqInt(t *testing.T) {
	tests := map[string]struct {
		input *fq
		want  *big.Int
	}{
		"zero": {input: fqZero, want: bigZero},
		"last element (montgomery form)": {input: fqLastHardcoded, want: bigLast},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Int()
			if got.Cmp(tc.want) != 0 {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
func TestFqAdd(t *testing.T) {
	tests := map[string]struct {
		x, y, want *fq
	}{
		"0 + 0 = 0":           {x: fqZero, y: fqZero, want: fqZero},
		"0 + y = y for y < q": {x: fqZero, y: fqLastHardcoded, want: fqLastHardcoded},
		"last elem + 1 = 0":   {x: fqLastStandardHardcoded, y: fqOneStandard, want: fqZero},
		"last elem + 2 = 1":   {x: fqLastStandardHardcoded, y: fqTwoStandard, want: fqOneStandard},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqAdd(&got, tc.x, tc.y)
			if got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqNeg(t *testing.T) {
	tests := map[string]struct {
		input, want *fq
	}{
		"q - last elem = 1": {input: fqLastStandardHardcoded, want: fqOneStandard},
		"q - 1 = last elem": {input: fqOneStandard, want: fqLastStandardHardcoded},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqNeg(&got, tc.input)
			if got != *tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

/*
func TestFqSub(t *testing.T) {
	tests := []struct {
		x, y, want *fq
	}{
		{a: two, b: fqMontOne, output: fqMontOne},
		{a: fqLastElement, b: fqLastElement, output: fqZero},
		{a: fqZero, b: fqLastElement, output: fqOne},
		{a: fq100, b: fq99, output: fqOne},
		{a: fqMontOne, b: two, output: *negFqMontOne},
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

/*


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
			a:      fqZero,
			b:      fqOne,
			output: fqLarge{0},
		},
		{
			a:      fq{0, 0, 0, 0, 1},
			b:      fqOne,
			output: fqLarge{0, 0, 0, 0, 1},
		},
		{
			a:      fqOne,
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
			a:      fqMontOne,
			b:      fqMontOne,
			output: fqMontOne,
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
			base:     fqMontOne,
			exponent: []uint64{3},
			output:   fqMontOne,
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
	negFqMontOne := new(fq)
	fqNeg(negFqMontOne, &fqMontOne)

	// TODO complete
	testCases := []struct {
		base     fq
		exponent []uint64
		output   fq
	}{
		{
			base:     *negFqMontOne,
			exponent: []uint64{123123},
			output:   *negFqMontOne,
		},
		{
			base:     fqMontOne,
			exponent: []uint64{123123},
			output:   fqMontOne,
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
*/
