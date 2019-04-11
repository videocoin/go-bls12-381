package bls12

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/*
func TestSign(t *testing.T) {
	privKey, err := GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		hash   []byte
		priv   *PrivateKey
		output []byte
	}{
		{
			hash:   []byte{7, 8, 9},
			priv:   privKey,
			output: []byte{1, 2, 3},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("hash: [% x], priv: %s\n", testCase.hash, testCase.priv.Secret.String()), func(t *testing.T) {
			sig := Sign(testCase.priv, testCase.hash)
			if !equal(sig, testCase.output) {
				t.Errorf("expected [% x], got [% x]\n", testCase.output, sig)
			}
		})
	}
}
*/
