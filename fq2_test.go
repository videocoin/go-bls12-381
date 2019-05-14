package bls12

import (
	"testing"
)

func TestFq2Set(t *testing.T) {
	tests := map[string]struct {
		input, want fq2
	}{
		"0 > input": {
			input: fq2{},
			want:  fq2{},
		},
		"c0 + c1X > input": {
			input: fq2{
				c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
				c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
			},
			want: fq2{
				c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
				c1: fq{0x40AB3263EFF0206, 0xEF148D1EA0F4C069, 0xECA8F3318332BB7A, 0x7E83A49A2E99D69, 0x32B7FFF2ED47FFFD, 0x43F5FFFFFFFCAAAE},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2SetZero(t *testing.T) {
	tests := map[string]struct {
		input fq2
	}{
		"0 > 0": {
			input: fq2{},
		},
		"c0 + c1X > 0": {
			input: fq2{
				c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
				c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SetZero()
			if (*got != fq2{}) {
				t.Fatalf("expected: %v, got: %v", fq2{}, got)
			}
		})
	}
}

func TestFq2SetOne(t *testing.T) {
	tests := map[string]struct {
		input fq2
	}{
		"0 > 1": {
			input: fq2{},
		},
		"c0 + c1X > 1": {
			input: fq2{
				c0: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
				c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.SetOne()
			if (*got != fq2{c0: *new(fq).SetUint64(1)}) {
				t.Fatalf("expected: %v, got: %v", fq2{c0: *new(fq).SetUint64(1)}, got)
			}
		})
	}
}

func TestFq2Add(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq2
	}{
		"(0 + 0X) + (0 + yX) = yX": {
			x:    fq2{},
			y:    fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
			want: fq2{c1: fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}},
		},
		"(1 + lastX) + (2 + 2X) = 3 + X": {
			x:    fq2{c0: fq{1}, c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A}},
			y:    fq2{c0: fq{2}, c1: fq{2}},
			want: fq2{c0: fq{3}, c1: fq{1}},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Add(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2Neg(t *testing.T) {
	tests := map[string]struct {
		input, want fq2
	}{
		"-(last + lastX) = 1 + X": {
			input: fq2{
				c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
			},
			want: fq2{
				c0: fq{1},
				c1: fq{1},
			},
		},
		"-(1 + X) = last + lastX": {
			input: fq2{
				c0: fq{1},
				c1: fq{1},
			},
			want: fq2{
				c0: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
				c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
			}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Neg(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2Sub(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq2
	}{
		"(2 + X) - (1 + 2X) = 1 + lastX": {
			x: fq2{
				c0: fq{2},
				c1: fq{1},
			},
			y: fq2{
				c0: fq{1},
				c1: fq{2},
			},
			want: fq2{
				c0: fq{1},
				c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
			},
		},
		"(1) - (1 + lastX) = X": {
			x: fq2{
				c0: fq{1},
			},
			y: fq2{
				c0: fq{1},
				c1: fq{0xB9FEFFFFFFFFAAAA, 0x1EABFFFEB153FFFF, 0x6730D2A0F6B0F624, 0x64774B84F38512BF, 0x4B1BA7B6434BACD7, 0x1A0111EA397FE69A},
			},
			want: fq2{
				c1: fq{1},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Sub(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2Mul(t *testing.T) {
	tests := map[string]struct {
		x, y, want fq2
	}{
		"tc 1": {
			x: fq2{
				c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
				c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
			},
			y: fq2{
				c0: fq{5508758831087832138, 6448303779119275098, 16710190169160573786, 13542242618704742751, 563980702369916322, 37152010398653157},
				c1: fq{12520284671833321565, 1777275927576994268, 9704602344324656032, 8739618045342622522, 16651875250601773805, 804950956836789234},
			},
			want: fq2{
				c0: fq{15827099533983733295, 5585758986002856792, 17089326063188852987, 9744188138385726684, 16136916081236286470, 539400444303524798},
				c1: fq{1625211833195412533, 6260280147765452499, 3544356884998474521, 9883248032979034166, 8703231040866562622, 339556463953896598},
			},
		},
		"tc 2": {
			x: fq2{
				c0: fq{2148348102093263616, 12232197882281708926, 11330363351339265390, 8919790901940522406, 17524282994943615806, 496043075549758450},
				c1: fq{406374380327561821, 16222300308001049590, 2744191801523148582, 9384378502465456106, 6088477103101489105, 1478173219842328727},
			},
			y: fq2{
				c0: fq{3656661975483800321, 2343032803162206357, 3148888360041182374, 13718521627034877454, 7582705438920995243, 1502510331453350514},
				c1: fq{12908108411122378854, 12884292956005326590, 14415309872191119313, 17916475641229970923, 8768506067738497560, 1577813823053251215},
			},
			want: fq2{
				c0: fq{16633317224619977984, 8114711876341744521, 5129895554745463763, 3028104440572865277, 600965037607129992, 1357055630732818745},
				c1: fq{16769752770056333816, 9539351424971908300, 1361078859430182954, 6919544897282332375, 8797177321797491027, 1057830777381385737},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Mul(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2MulXi(t *testing.T) {
	tests := map[string]struct {
		input, want fq2
	}{
		"(2 + X)ξ = 1 + 3X": {
			input: fq2{
				c0: fq{2},
				c1: fq{1},
			},
			want: fq2{
				c0: fq{1},
				c1: fq{3},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).MulXi(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2Sqr(t *testing.T) {
	tests := map[string]struct {
		input, want fq2
	}{
		"tc 1": {
			input: fq2{
				c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
				c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
			},
			want: fq2{
				c0: fq{5183593039390375737, 10963027822502823039, 6255345974967782363, 17684205669924779383, 13794376949041289905, 789947231065766105},
				c1: fq{2921497446257912465, 731946419108638375, 5871846982883770661, 14103266165668144248, 17935390935820665642, 1741923485045802819},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq2).Sqr(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq2Inv(t *testing.T) {
	// TODO
}
