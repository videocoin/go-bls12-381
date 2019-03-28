package bls12

/*
func TestTwistPointAdd(t *testing.T) {
	testCases := []struct {
		a, b, expectedPoint *twistPoint
	}{
		{
			a:             g2Generator,
			b:             twist0,
			expectedPoint: g2Generator,
		},
	}
	for _, testCase := range testCases {
		// TODO
		t.Run("", func(t *testing.T) {
			result := new(twistPoint).Add(testCase.a, testCase.b)
			if !result.Equal(testCase.expectedPoint) {
				// TODO type used in the printf
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
