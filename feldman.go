package bls12

import (
	"errors"
	"io"
	"math/big"
)

var (
	errInvalidThreshold       = errors.New("feldman: minimum required number of shares is greater than the total number of shares")
	errInvalidNumCoefficients = errors.New("feldman: polynomial requires at least one coefficient")
	errEmptyVerificationVec   = errors.New("feldman: empty verification vector")
	errInvalidShare           = errors.New("feldman: share is not valid")
)

// Share represents a unique part of a secret.
type Share = Point

// Point is a polynomial point.
type Point struct {
	X uint64
	Y *PrivateKey
}

type polynomial struct {
	coefficients []fq
}

func newPolynomial(coefficients []*PrivateKey) (*polynomial, error) {
	if len(coefficients) == 0 {
		return nil, errInvalidNumCoefficients
	}

	p := &polynomial{
		coefficients: make([]fq, 0, len(coefficients)),
	}

	for _, priv := range coefficients {
		coeff, err := FqMontgomeryFromBig(priv.Secret)
		if err != nil {
			return nil, err
		}
		p.coefficients = append(p.coefficients, coeff)
	}

	return p, nil
}

func (p *polynomial) evaluate(x uint64) (*PrivateKey, error) {
	fqX, err := FqMontgomeryFromBig(new(big.Int).SetUint64(x))
	if err != nil {
		return nil, err
	}
	mul := fqMont1
	sum := p.coefficients[0]
	for _, coeff := range p.coefficients[1:] {
		term := new(fq)
		fqMul(&mul, &mul, &fqX)
		fqMul(term, &coeff, &mul)
		fqAdd(&sum, &sum, term)
	}

	return privKeyFromScalar(fqToBig(fqFromFqMontgomery(sum))), nil
}

// CreateShares divides the secret into parts, giving each participant its own unique part.
func CreateShares(reader io.Reader, threshold uint64, numShares uint64) ([]*PublicKey, []*Share, *PrivateKey, error) {
	if threshold > numShares {
		return nil, nil, nil, errInvalidThreshold
	}

	// generate coefficients
	secrets := make([]*PrivateKey, 0, threshold)
	verification := make([]*PublicKey, 0, threshold)
	for i := uint64(0); i < threshold; i++ {
		/*
			priv, err := GenerateKey(reader)
			if err != nil {
				return nil, nil, nil, err
			}
		*/
		priv := privKeyFromScalar(new(big.Int).SetUint64(1))
		pub := priv.Public()
		secrets = append(secrets, priv)
		verification = append(verification, &pub)
	}

	randPolynomial, err := newPolynomial(secrets)
	if err != nil {
		return nil, nil, nil, err
	}

	shares := make([]*Share, 0, numShares)
	for i := uint64(1); i <= numShares; i++ {
		secret, err := randPolynomial.evaluate(i)
		if err != nil {
			return nil, nil, nil, err
		}
		shares = append(shares, &Point{
			X: i,
			Y: secret,
		})
	}

	return verification, shares, secrets[0], nil
}

// PrivKeyFromShares reconstructs the secret using Lagrange polynomials.
// Passing less shares than the minimum required results in the wrong secret.
func PrivKeyFromShares(shares []*Share) (*PrivateKey, error) {
	// See https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing#Computationally_efficient_approach
	ids := make([]fq, 0, len(shares))
	secrets := make([]fq, 0, len(shares))
	for _, share := range shares {
		idMont, err := FqMontgomeryFromBig(new(big.Int).SetUint64(share.X))
		if err != nil {
			return nil, err
		}
		ids = append(ids, idMont)
		secretMont, err := FqMontgomeryFromBig(share.Y.Secret)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secretMont)
	}

	sum := new(fq)
	for i := 0; i < len(shares); i++ {
		mul := fqMont1
		for j := 0; j < len(shares); j++ {
			if j != i {
				term := new(fq)
				fqSub(term, &ids[j], &ids[i])
				fqInv(term, term)
				fqMul(term, term, &ids[j])
				fqMul(&mul, &mul, term)
			}
		}
		fqMul(&mul, &mul, &secrets[i])
		fqAdd(sum, sum, &mul)
	}

	return privKeyFromScalar(fqToBig(fqFromFqMontgomery(*sum))), nil
}

type publicPolynomial struct {
	coefficients []*PublicKey
}

func newPublicPolynomial(coefficients []*PublicKey) (*publicPolynomial, error) {
	if len(coefficients) == 0 {
		return nil, errInvalidNumCoefficients
	}

	return &publicPolynomial{
		coefficients: coefficients,
	}, nil
}

func (p *publicPolynomial) evaluate(x uint64) (*PublicKey, error) {
	bigX := new(big.Int).SetUint64(x)
	sum := newG2Point().Set(p.coefficients[0])
	mul := new(big.Int).SetUint64(1)
	for _, coeff := range p.coefficients[1:] {
		mul.Mul(mul, bigX)
		sum.Add(sum, newG2Point().ScalarMult(coeff, mul))
	}

	return sum, nil
}

// VerifyShare verifies that a received secret key share is actually the result
// of the evaluation of the secret polynomial.
func VerifyShare(share *Share, verificationVec []*PublicKey) error {
	if len(verificationVec) == 0 {
		return errEmptyVerificationVec
	}

	expectedPubKey := share.Y.Public()
	publicPolynomial, err := newPublicPolynomial(verificationVec)
	pubKey, err := publicPolynomial.evaluate(share.X)
	//fmt.Printf("Expected: %s, got: %s", expectedPubKey.String(), pubKey.String())
	if err != nil {
		return err
	}
	if !pubKey.Equal(&expectedPubKey) {
		return errInvalidShare
	}

	return nil
}
