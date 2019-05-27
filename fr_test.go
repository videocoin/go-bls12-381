package bls12

import "testing"

func TestFrMontgomeryEncode(t *testing.T) {

	tests := map[string]struct {
		input, want fr
	}{
		"0 > mont(0)": {input: fr{}, want: fr{}},
		"1 > mont(1)": {input: fr{1}, want: fr{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fr).MontgomeryEncode(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFrMontgomeryDecode(t *testing.T) {
	tests := map[string]struct {
		input, want fr
	}{
		"mont(0) > 0": {input: fr{}, want: fr{}},
		"mont(1) > 1": {input: fr{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f}, want: fr{1}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fr).MontgomeryDecode(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFrMod(t *testing.T) {
	tests := map[string]struct {
		input, want fr
	}{
		"0 > 0":       {input: fr{}, want: fr{}},
		"last > last": {input: fr{0xffffffff00000000, 0x53bda402fffe5bfe, 0x3339d80809a1d805, 0x73eda753299d7d48}, want: fr{0xffffffff00000000, 0x53bda402fffe5bfe, 0x3339d80809a1d805, 0x73eda753299d7d48}},
		"r > 0":       {input: fr{0xFFFFFFFF00000001, 0x53BDA402FFFE5BFE, 0x3339D80809A1D805, 0x73EDA753299D7D48}, want: fr{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fr).Set(&tc.input)
			frMod(got)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, *got)
			}
		})
	}
}

// note: addition in the montgomery form is the same as ordinary modular addition.
func TestFrAdd(t *testing.T) {
	tests := map[string]struct {
		x, y, want fr
	}{
		"0 + 0 = 0":    {x: fr{}, y: fr{}, want: fr{}},
		"0 + y = y":    {x: fr{}, y: fr{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69}, want: fr{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69}},
		"1 + 2 = 3":    {x: fr{1}, y: fr{2}, want: fr{3}},
		"last + 1 = 0": {x: fr{0xffffffff00000000, 0x53bda402fffe5bfe, 0x3339d80809a1d805, 0x73eda753299d7d48}, y: fr{1}, want: fr{}},
		"last + 2 = 1": {x: fr{0xffffffff00000000, 0x53bda402fffe5bfe, 0x3339d80809a1d805, 0x73eda753299d7d48}, y: fr{2}, want: fr{1}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fr
			frAdd(&got, &tc.x, &tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFrMul(t *testing.T) {
	tests := map[string]struct {
		x, y, want fr
	}{
		"mont(1) * mont(1) = mont(1)": {
			x:    fr{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f},
			y:    fr{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f},
			want: fr{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f},
		},

		/*
			TODO Review this one
				"mont(1) * mont(last) = mont(last)": {
					x:    fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493},
					y:    fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
					want: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
				},
		*/
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got fr
			frMul(&got, &tc.x, &tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
