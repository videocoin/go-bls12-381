package bls12

import "testing"

func TestFq6Set(t *testing.T) {
	tests := map[string]struct {
		input, want fq6
	}{
		"0 > input": {input: fq6{}, want: fq6{}},
		"c0 + c1X + c2X^2 > input": {
			input: fq6{
				c0: fq2{
					c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
					c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
				},
				c1: fq2{
					c0: fq{3},
					c1: fq{4},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{6},
				},
			},
			want: fq6{
				c0: fq2{
					c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
					c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
				},
				c1: fq2{
					c0: fq{3},
					c1: fq{4},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{6},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq6).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq6SetZero(t *testing.T) {
	tests := map[string]struct {
		input fq6
	}{
		"0 > 0": {
			input: fq6{},
		},
		"c0 + c1X + c2X^2 > 0": {
			input: fq6{
				c0: fq2{
					c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
					c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
				},
				c1: fq2{
					c0: fq{3},
					c1: fq{4},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{6},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SetZero()
			if (*got != fq6{}) {
				t.Fatalf("expected: %v, got: %v", fq6{}, got)
			}
		})
	}
}

func TestFq6SetOne(t *testing.T) {
	tests := map[string]struct {
		input fq6
	}{
		"0 > 1": {
			input: fq6{},
		},
		"c0 + c1X + c2X^2 > 1": {
			input: fq6{
				c0: fq2{
					c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
					c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
				},
				c1: fq2{
					c0: fq{3},
					c1: fq{4},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{6},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SetOne()
			if (*got != fq6{c0: fq2{c0: *new(fq).SetUint64(1)}}) {
				t.Fatalf("expected: %v, got: %v", fq6{c0: fq2{c0: *new(fq).SetUint64(1)}}, got)
			}
		})
	}
}

func TestFq6Add(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq6
	}{
		"0 + ((0 + yX) + (0 + yX)B + (0 + yX)B^2) = yX + yXB + yXB^2": {
			x: fq6{
				c0: fq2{},
				c1: fq2{},
				c2: fq2{},
			},
			y: fq6{
				c0: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
				c1: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
				c2: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
			},
			want: fq6{
				c0: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
				c1: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
				c2: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
			},
		},
		"((1 + lastX) + (1 + lastX)B + (1 + lastX)B^2) + ((2 + 2X) + (2 + 2X)B + (2 + 2X)B^2) = (3 + X) + (3 + X)B + (3 + X)B^2": {
			x: fq6{
				c0: fq2{c0: fq{1}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
				c1: fq2{c0: fq{1}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
				c2: fq2{c0: fq{1}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
			},
			y: fq6{
				c0: fq2{c0: fq{2}, c1: fq{2}},
				c1: fq2{c0: fq{2}, c1: fq{2}},
				c2: fq2{c0: fq{2}, c1: fq{2}},
			},
			want: fq6{
				c0: fq2{c0: fq{3}, c1: fq{1}},
				c1: fq2{c0: fq{3}, c1: fq{1}},
				c2: fq2{c0: fq{3}, c1: fq{1}},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq6).Add(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq6Neg(t *testing.T) {
	tests := map[string]struct {
		input, want fq6
	}{
		"-((last + lastX) + (last + lastX)B + (last + lastX)B^2) = (1 + X) + (1 + X)B + (1 + X)B^2": {
			input: fq6{
				c0: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
				c1: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
			},
			want: fq6{
				c0: fq2{
					c0: fq{1},
					c1: fq{1},
				},
				c1: fq2{
					c0: fq{1},
					c1: fq{1},
				},
				c2: fq2{
					c0: fq{1},
					c1: fq{1},
				},
			},
		},
		"-((1 + X) + (1 + X)B + (1 + X)B^2) = (last + lastX) + (last + lastX)B + (last + lastX)B^2": {
			input: fq6{
				c0: fq2{
					c0: fq{1},
					c1: fq{1},
				},
				c1: fq2{
					c0: fq{1},
					c1: fq{1},
				},
				c2: fq2{
					c0: fq{1},
					c1: fq{1},
				},
			},
			want: fq6{
				c0: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
				c1: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
				c2: fq2{
					c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
					c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq6).Neg(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

/*



func TestFq6Add(t *testing.T) {
	// TODO
}

func TestFq6Sub(t *testing.T) {
	// TODO
}

func TestFq6Mul(t *testing.T) {
	// TODO
}

func TestFq6SparseMult(t *testing.T) {
	// TODO
}

func TestFq6MulQuadraticNonResidue(t *testing.T) {
	// TODO
}

func TestFq6Sqr(t *testing.T) {
	// TODO
}

func TestFq6Inv(t *testing.T) {
	// TODO
}
*/
