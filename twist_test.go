package bls12

/*
func TestTwistPointAdd(t *testing.T) {
	testCases := []struct {
		a, b, expectedPoint *twistPoint
	}{
		{
			a: g2Generator,
			b: &twistPoint{
				x: fq2{},
				y: fq2{},
				z: fq2{FqMont1, Fq0},
			},
			expectedPoint: g2Generator,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("A: %v, B: %v \n", testCase.a, testCase.b), func(t *testing.T) {
			result := new(twistPoint).Add(testCase.a, testCase.b)
			if !result.Equal(testCase.expectedPoint) {
				t.Errorf("expected %v, got %v\n", testCase.expectedPoint, result)
			}
		})
	}
}

func TestScalarMult(t *testing.T) {
	testCases := []struct {
		point         *twistPoint
		scalar        *big.Int
		expectedPoint *twistPoint
	}{
		{
			point:         g2Generator,
			scalar:        big1,
			expectedPoint: g2Generator,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("point: %v, scalar: %s\n", testCase.point, testCase.scalar.String()), func(t *testing.T) {
			result := new(twistPoint).ScalarMult(testCase.point, testCase.scalar)
			if !result.Equal(testCase.expectedPoint) {
				t.Errorf("expected %v, got %v\n", testCase.expectedPoint, result)
			}
		})
	}
}
*/
