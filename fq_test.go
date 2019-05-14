package bls12

import (
	"math/big"
	"testing"
)

var (
	fqZero         = fq{}
	fqOne          = fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}
	fqLast         = fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}
	fqLastStandard = fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}
	fqTwoStandard  = fq{2}
	bigOne         = big.NewInt(1)
	bigZero        = big.NewInt(0)
	bigLast, _     = bigFromBase10("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786")
)

// note: fq operates, internally, on the montgomery form.
func TestFqIsOne(t *testing.T) {
	tests := map[string]struct {
		input fq
		want  bool
	}{
		"mont(0)":    {input: fq{}, want: false},
		"mont(1)":    {input: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, want: true},
		"mont(last)": {input: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, want: false},
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
		input, want fq
	}{
		"zero": {input: fq{}, want: fq{}},
		"non-zero": {
			input: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
			want:  fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetString(t *testing.T) {
	tests := map[string]struct {
		input string
		want  fq
		err   error
	}{
		"0 > mont(0)":       {input: "0", want: fq{}, err: nil},
		"1 > mont(1)":       {input: "1", want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, err: nil},
		"last > mont(last)": {input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", want: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, err: nil},
		"out of bounds":     {input: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", err: errOutOfBounds},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := new(fq).SetString(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if err == nil && *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetInt(t *testing.T) {
	tests := map[string]struct {
		input *big.Int
		want  fq
		err   error
	}{
		"0 > mont(0)":       {input: new(big.Int).SetUint64(0), want: fq{}, err: nil},
		"1 > mont(1)":       {input: new(big.Int).SetUint64(1), want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, err: nil},
		"last > mont(last)": {input: new(big.Int).Sub(q, bigOne), want: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, err: nil},
		"out of bounds":     {input: q, err: errOutOfBounds},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := new(fq).SetInt(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if err == nil && *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSetUint64(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  fq
	}{
		"0 > mont(0)": {input: 0, want: fq{}},
		"1 > mont(1)": {input: 1, want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).SetUint64(tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", fqZero, got)
			}
		})
	}
}

func TestFqMontgomeryEncode(t *testing.T) {
	tests := map[string]struct {
		input, want fq
	}{
		"0 > mont(0)": {input: fq{}, want: fq{}},
		"1 > mont(1)": {input: fq{1}, want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).MontgomeryEncode(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqMontgomeryDecode(t *testing.T) {
	tests := map[string]struct {
		input, want fq
	}{
		"mont(0) > 0": {input: fq{}, want: fq{}},
		"mont(1) > 1": {input: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, want: fq{1}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq).MontgomeryDecode(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqInt(t *testing.T) {
	last, valid := new(big.Int).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", 10)
	if !valid {
		t.Fatal(valid)
	}
	tests := map[string]struct {
		input fq
		want  *big.Int
	}{
		"mont(0) > big(0)":       {input: fq{}, want: new(big.Int).SetUint64(0)},
		"mont(1) > big(1)":       {input: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, want: new(big.Int).SetUint64(1)},
		"mont(last) > big(last)": {input: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, want: last},
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

// note: addition in the montgomery form is the same as ordinary modular addition.
func TestFqAdd(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq
	}{
		"0 + 0 = 0":    {x: fq{}, y: fq{}, want: fq{}},
		"0 + y = y":    {x: fq{}, y: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, want: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
		"1 + 2 = 3":    {x: fq{1}, y: fq{2}, want: fq{3}},
		"last + 1 = 0": {x: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, y: fq{1}, want: fq{}},
		"last + 2 = 1": {x: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, y: fq{2}, want: fq{1}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqAdd(&got, &tc.x, &tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// note: neg input = field order - input
func TestFqNeg(t *testing.T) {
	tests := map[string]struct {
		input, want fq
	}{
		"-last = 1": {input: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, want: fq{1}},
		"-1 = last": {input: fq{1}, want: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqNeg(&got, &tc.input)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// note: subtraction in the montgomery form is the same as ordinary modular subtraction.
func TestFqSub(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq
	}{
		"2 - 1 = 1":    {x: fq{2}, y: fq{1}, want: fq{1}},
		"1 - 2 = last": {x: fq{1}, y: fq{2}, want: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
		"1 - 1 = 0":    {x: fq{1}, y: fq{1}, want: fq{}},
		"0 - last = 1": {x: fq{}, y: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, want: fq{1}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqSub(&got, &tc.x, &tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// note: to multiply x and y, they are first converted to Montgomery form.
func TestFqMul(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq
	}{
		"mont(1) * mont(1) = mont(1)": {
			x:    fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			y:    fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
		},
		"mont(1) * mont(last) = mont(last)": {
			x:    fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
			y:    fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
			want: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqMul(&got, &tc.x, &tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqExp(t *testing.T) {
	tests := map[string]struct {
		x    fq
		y    []uint64
		want fq
	}{
		"mont(1) ^ 7 = mont(1)": {x: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, y: []uint64{7}, want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqExp(&got, &tc.x, tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFqSqrt(t *testing.T) {
	tests := map[string]struct {
		input, want fq
	}{
		"Sqrt(mont(1)) = mont(1)": {input: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}, want: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fq
			fqSqrt(&got, &tc.input)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

/*
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

*/
