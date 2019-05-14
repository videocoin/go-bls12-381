package bls12

import "testing"

func TestFq12Set(t *testing.T) {
	tests := map[string]struct {
		input, want fq12
	}{
		"0 > input": {
			input: fq12{},
			want:  fq12{},
		},
		"c0 + c1X > input": {
			input: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			want: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12SetOne(t *testing.T) {
	tests := map[string]struct {
		input fq12
	}{
		"0 > 1": {
			input: fq12{},
		},
		"c0 + c1X > 1": {
			fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SetOne()
			if (*got != fq12{c0: fq6{c0: fq2{c0: *new(fq).SetUint64(1)}}}) {
				t.Fatalf("expected: %v, got: %v", fq12{c0: fq6{c0: fq2{c0: *new(fq).SetUint64(1)}}}, got)
			}
		})
	}
}

func TestFq12Equal(t *testing.T) {
	tests := map[string]struct {
		x, y fq12
		want bool
	}{
		"0 = 0": {
			x:    fq12{},
			y:    fq12{},
			want: true,
		},
		"tc 2": {
			x: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			y: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			want: true,
		},
		"tc 3": {
			x: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
			},
			y: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			want: false,
		},
		"tc 4": {
			x: fq12{
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			y: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			want: false,
		},
		"tc 5": {
			x: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
			},
			y: fq12{
				c1: fq6{
					c0: fq2{c0: fq{7}, c1: fq{8}},
					c1: fq2{c0: fq{9}, c1: fq{10}},
					c2: fq2{c0: fq{11}, c1: fq{12}},
				},
			},
			want: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.x.Equal(&tc.y)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Conjugate(t *testing.T) {
	tests := map[string]struct {
		input, want fq12
	}{
		//"conjugate(0) = 0": {input: fq12{}, want: fq12{}}, TODO
		"conjugate(c0 + c1U) = c0 - c1U": {
			input: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{1}, c1: fq{1}},
					c1: fq2{c0: fq{1}, c1: fq{1}},
					c2: fq2{c0: fq{1}, c1: fq{1}},
				},
			},
			want: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
					c1: fq2{c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
					c2: fq2{c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Conjugate(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Add(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq12
	}{
		"tc 1": {
			x: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{1}, c1: fq{1}},
					c1: fq2{c0: fq{2}, c1: fq{1}},
					c2: fq2{c0: fq{1}, c1: fq{1}},
				},
			},
			y: fq12{
				c0: fq6{
					c0: fq2{c0: fq{1}, c1: fq{2}},
					c1: fq2{c0: fq{3}, c1: fq{4}},
					c2: fq2{c0: fq{5}, c1: fq{6}},
				},
				c1: fq6{
					c0: fq2{c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, c1: fq{1}},
					c1: fq2{c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}, c1: fq{1}},
					c2: fq2{c0: fq{1}, c1: fq{1}},
				},
			},
			want: fq12{
				c0: fq6{
					c0: fq2{c0: fq{2}, c1: fq{4}},
					c1: fq2{c0: fq{6}, c1: fq{8}},
					c2: fq2{c0: fq{10}, c1: fq{12}},
				},
				c1: fq6{
					c0: fq2{c0: fq{}, c1: fq{2}},
					c1: fq2{c0: fq{1}, c1: fq{2}},
					c2: fq2{c0: fq{2}, c1: fq{2}},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Add(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Mul(t *testing.T) {
	// TODO
}

func TestFq12SparseMult(t *testing.T) {
	// TODO
}

func TestFq12Sqr(t *testing.T) {
	// TODO
}

func TestFq12Inv(t *testing.T) {
	// TODO
}

func TestFq12Exp(t *testing.T) {
	// TODO
}

func TestFq12Frobenius(t *testing.T) {
	// TODO
}
