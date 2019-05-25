package bls12

import (
	"errors"
	"io"
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
	coefficients []*fr
}

func newPolynomial(coefficients []*PrivateKey) (*polynomial, error) {
	if len(coefficients) == 0 {
		return nil, errInvalidNumCoefficients
	}

	p := &polynomial{
		coefficients: make([]*fr, 0, len(coefficients)),
	}

	for _, priv := range coefficients {
		coeff, _ := new(fr).SetInt(priv.Secret)
		p.coefficients = append(p.coefficients, coeff)
	}

	return p, nil
}

func (p *polynomial) evaluate(x uint64) *PrivateKey {
	frX := new(fr).SetUint64(x)
	mul := new(fr).SetUint64(1)
	sum := new(fr).Set(p.coefficients[0])
	term := new(fr)
	for _, coeff := range p.coefficients[1:] {
		frMul(mul, mul, frX)
		frMul(term, coeff, mul)
		frAdd(sum, sum, term)
	}

	return privKeyFromScalar(sum.Int())
}

// CreateShares divides the secret into parts, giving each participant its own unique part.
func CreateShares(reader io.Reader, priv *PrivateKey, threshold uint64, numShares uint64) ([]*PublicKey, []*Share, *PrivateKey, error) {
	if threshold > numShares {
		return nil, nil, nil, errInvalidThreshold
	}

	secrets := make([]*PrivateKey, 0, threshold)
	verification := make([]*PublicKey, 0, threshold)

	// use an existent secret instead of generating a new one
	if priv != nil {
		secrets = append(secrets, priv)
		verification = append(verification, &priv.PublicKey)
	}

	for i := uint64(len(secrets)); i < threshold; i++ {
		priv, err := GenerateKey(reader)
		if err != nil {
			return nil, nil, nil, err
		}

		secrets = append(secrets, priv)
		verification = append(verification, &priv.PublicKey)
	}

	randPolynomial, err := newPolynomial(secrets)
	if err != nil {
		return nil, nil, nil, err
	}

	shares := make([]*Share, 0, numShares)
	for i := uint64(1); i <= numShares; i++ {
		shares = append(shares, &Point{
			X: i,
			Y: randPolynomial.evaluate(i),
		})
	}

	return verification, shares, secrets[0], nil
}

// PrivKeyFromShares reconstructs the secret using Lagrange polynomials.
// Passing less shares than the minimum required results in the wrong secret.
func PrivKeyFromShares(shares []*Share) (*PrivateKey, error) {
	// See https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing#Computationally_efficient_approach
	ids := make([]*fr, 0, len(shares))
	secrets := make([]*fr, 0, len(shares))
	for _, share := range shares {
		ids = append(ids, new(fr).SetUint64(share.X))
		secret, err := new(fr).SetInt(share.Y.Secret)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	sum := new(fr)
	for i := 0; i < len(shares); i++ {
		mul := new(fr).SetUint64(1)
		for j := 0; j < len(shares); j++ {
			if j != i {
				term := new(fr)
				frSub(term, ids[j], ids[i])
				frInv(term, term)
				frMul(term, term, ids[j])
				frMul(mul, mul, term)
			}
		}
		frMul(mul, mul, secrets[i])
		frAdd(sum, sum, mul)
	}

	return privKeyFromScalar(sum.Int()), nil
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

func (p *publicPolynomial) evaluate(x uint64) *PublicKey {
	frX := new(fr).SetUint64(x)
	mul := new(fr).SetUint64(1)
	sum := new(g2Point).Set(p.coefficients[0])
	term := new(g2Point)
	for _, coeff := range p.coefficients[1:] {
		frMul(mul, mul, frX)
		sum.Add(sum, term.ScalarMult(coeff, mul.Int()))
	}

	return sum.ToAffine()
}

// VerifyShare verifies that a received secret key share is actually the result
// of the evaluation of the secret polynomial.
func VerifyShare(share *Share, verificationVec []*PublicKey) error {
	if len(verificationVec) == 0 {
		return errEmptyVerificationVec
	}

	expectedPubKey := share.Y.Public()
	publicPolynomial, err := newPublicPolynomial(verificationVec)
	if err != nil {
		return err
	}
	pubKey := publicPolynomial.evaluate(share.X)

	if !pubKey.Equal(&expectedPubKey) {
		return errInvalidShare
	}

	return nil
}
