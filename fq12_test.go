package bls12

import (
	"math/big"
	"testing"
)

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
			// TODO replace set uint with hardcoded value
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
		"conjugate(0) = 0": {input: fq12{}, want: fq12{}},
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
	tests := map[string]struct {
		x, y, want fq12
	}{
		"0 * 0 = 0": {
			x:    fq12{},
			y:    fq12{},
			want: fq12{},
		},
		"0 * mont(1) = 0": {
			x:    fq12{},
			y:    fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			want: fq12{},
		},
		"mont(1) * mont(1) = mont(1)": {
			x:    fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			y:    fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			want: fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Mul(&tc.x, &tc.y)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12SparseMul014(t *testing.T) {
	tests := map[string]struct {
		a          fq12
		b0, b1, b2 fq2
		want       fq12
	}{
		"tc1": {
			a: fq12{c0: fq6{c0: fq2{c0: fq{8505329371266088957, 17002214543764226050, 6865905132761471162, 8632934651105793861, 6631298214892334189, 1582556514881692819}}}},
			b0: fq2{
				c0: fq{16521250154091893390, 11550767693645120122, 3517755699210907354, 8035647930683680220, 4019796556552419771, 451420072861892737},
				c1: fq{18293378936511423336, 13143741527255503601, 7522477062776238534, 11407560980233213841, 15716159397276202569, 1403106181814332372},
			},
			b1: fq2{
				c0: fq{4449984651635275401, 11680217845942513937, 11140531346451103816, 10500440204658996834, 16771094099020012859, 527545683060335776},
				c1: fq{18097767337011887344, 12600677166075290804, 13638019292209587518, 16188658819672867276, 12891218261365710021, 1175533863824620827},
			},
			b2: fq2{
				c0: fq{11033082401322288522, 16152799148149621672, 5480658672849822104, 12178979763090086206, 16683590538352871296, 1233599904214129433},
				c1: fq{9019646775837872212, 2854805697774761111, 7183090960508055194, 11070723832269891145, 15933306181532061767, 354030584485544670},
			},
			want: fq12{
				c0: fq6{
					c0: fq2{
						c0: fq{16521250154091893390, 11550767693645120122, 3517755699210907354, 8035647930683680220, 4019796556552419771, 451420072861892737},
						c1: fq{18293378936511423336, 13143741527255503601, 7522477062776238534, 11407560980233213841, 15716159397276202569, 1403106181814332372},
					},
					c1: fq2{
						c0: fq{4449984651635275401, 11680217845942513937, 11140531346451103816, 10500440204658996834, 16771094099020012859, 527545683060335776},
						c1: fq{18097767337011887344, 12600677166075290804, 13638019292209587518, 16188658819672867276, 12891218261365710021, 1175533863824620827},
					},
				},
				c1: fq6{
					c1: fq2{
						c0: fq{11033082401322288522, 16152799148149621672, 5480658672849822104, 12178979763090086206, 16683590538352871296, 1233599904214129433},
						c1: fq{9019646775837872212, 2854805697774761111, 7183090960508055194, 11070723832269891145, 15933306181532061767, 354030584485544670},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).SparseMul014(&tc.a, &tc.b0, &tc.b1, &tc.b2)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Sqr(t *testing.T) {
	tests := map[string]struct {
		input, want fq12
	}{
		"sqr(0) = 0": {
			input: fq12{},
			want:  fq12{},
		},
		"sqr(mont(1)) = mont(1)": {
			input: fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			want:  fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Sqr(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Inv(t *testing.T) {
	tests := map[string]struct {
		input, want fq12
	}{
		"inv(1) = 1": {
			input: fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			want:  fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Inv(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Exp(t *testing.T) {
	tests := map[string]struct {
		base fq12
		exp  *big.Int
		want fq12
	}{
		"1^X = 1": {
			base: fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
			exp:  bigU,
			want: fq12{c0: fq6{c0: fq2{c0: fq{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}}}},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Exp(&tc.base, tc.exp)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFq12Frobenius(t *testing.T) {
	tests := map[string]struct {
		input fq12
		power uint64
		want  fq12
	}{
		"frob(input, 0) = input": {
			input: fq12{
				c0: fq6{
					c0: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c1: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c2: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
				},
				c1: fq6{
					c0: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c1: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c2: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
				},
			},
			power: 0,
			want: fq12{
				c0: fq6{
					c0: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c1: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c2: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
				},
				c1: fq6{
					c0: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c1: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
					c2: fq2{
						c0: fq{17722385409647053328, 12967546844987299354, 11648722842835150208, 10994581490347323113, 8027586497049998955, 396758299565931735},
						c1: fq{11937283898719073798, 12295044263989567683, 4301357764460312582, 1953074377943790439, 14030662337566180679, 1266120665323335155},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(fq12).Frobenius(&tc.input, tc.power)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
