package bls12

import "testing"

func TestCurvePointSet(t *testing.T) {
	tests := map[string]struct {
		input, want curvePoint
	}{
		"zero":     {input: curvePoint{}, want: curvePoint{}},
		"non-zero": {input: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}}, want: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(curvePoint).Set(&tc.input)
			if *got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointEqual(t *testing.T) {
	tests := map[string]struct {
		a, b curvePoint
		want bool
	}{
		"zero": {
			a:    curvePoint{},
			b:    curvePoint{},
			want: true,
		},
		"tc 2": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			want: true,
		},
		"tc 3": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{2}, y: fq{2}, z: fq{3}},
			want: false,
		},
		"tc 4": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{3}, z: fq{3}},
			want: false,
		},
		"tc 5": {
			a:    curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			b:    curvePoint{x: fq{1}, y: fq{2}, z: fq{1}},
			want: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(&tc.b)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointIsInfinity(t *testing.T) {
	tests := map[string]struct {
		input curvePoint
		want  bool
	}{
		"zero": {
			input: curvePoint{},
			want:  true,
		},
		"z = 0": {
			input: curvePoint{x: fq{1}, y: fq{2}, z: fq{}},
			want:  true,
		},

		"z != 0": {
			input: curvePoint{x: fq{1}, y: fq{2}, z: fq{3}},
			want:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.IsInfinity()
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCurvePointAdd(t *testing.T) {
	// TODO
}

func TestCurvePointDouble(t *testing.T) {
	// TODO
}

func TestCurvePointScalarMult(t *testing.T) {
	// TODO
}

func TestCurvePointToAffine(t *testing.T) {
	// TODO
}

func TestCurvePointMarshal(t *testing.T) {
	// TODO
}

func TestCurvePointUnmarshal(t *testing.T) {
	// TODO
}

func TestCurvePointSetBytes(t *testing.T) {
	// TODO
}

func TestCurvePointSWEncode(t *testing.T) {
	// TODO
}
