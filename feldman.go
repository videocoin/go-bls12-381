package bls12

import (
	"errors"
	"io"
	"math/big"
)

var (
	errReqGreaterThanTotal    = errors.New("feldman: minimum required number of shares is greater than the total number of shares")
	errMissingRequiredSecret  = errors.New("feldman: missing required secret")
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
	coefficients []Fq
}

func newPolynomial(coefficients []*PrivateKey) (*polynomial, error) {
	if len(coefficients) == 0 {
		return nil, errInvalidNumCoefficients
	}

	p := &polynomial{
		coefficients: make([]Fq, 0, len(coefficients)),
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
	mul := FqMont1
	sum := p.coefficients[0]
	for _, coeff := range p.coefficients[1:] {
		term := new(Fq)
		FqMul(&mul, &mul, &fqX)
		FqMul(term, &coeff, &mul)
		FqAdd(&sum, &sum, term)
	}

	return privKeyFromScalar(FqToBig(FqFromFqMontgomery(sum))), nil
}

// CreateShares divides the secret into parts, giving each participant its own unique part.
func CreateShares(reader io.Reader, threshold uint64, numShares uint64) ([]*PublicKey, []*Share, *PrivateKey, error) {
	if threshold > numShares {
		return nil, nil, nil, errReqGreaterThanTotal
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

/*
// SecretToShares
func SecretFromShares(reader io.Reader, threshold uint64, numShares uint64, priv *PrivateKey) ([]*Share, *PrivateKey, error) {
	if priv == nil {
		return nil, nil, errMissingRequiredSecret
	}


	freeCoeff, err := FqMontgomeryFromBig(priv.Secret)
	if err != nil {
		return nil, nil, err
	}
	coefficient := freeCoeff

	mul := make([]Fq, numShares) // caches the variable multiplications
	ids := make([]Fq, numShares) // caches the share index (montgomery form)

	for i := uint64(0); i < threshold; i++ {
		// generate coefficient
		if i != 0 {
			var err error
			pk, err := GenerateKey(reader)
			if err != nil {
				return nil, nil, err
			}
			coefficient, err = FqMontgomeryFromBig(pk.Secret)
			if err != nil {
				return nil, nil, err
			}
		}

		// calculate polynomial term
		for j := uint64(0); j < numShares; j++ {
			if i == 0 {
				index := j + 1


				// create share
				shares = append(shares, &Point{X: x, Y: coefficients[i]})

				mul[j] = FqMont1

				continue
			}

			// sum term
			term := new(Fq)
			FqMul(&mul[j], &mul[j], &shares[j].X)
			FqMul(term, &coefficients[i], &mul[j])
			FqAdd(&shares[j].Y, &shares[j].Y, term)
		}
	}

	return shares, coefficients, nil
}

*/
// PrivKeyFromShares reconstructs the secret using Lagrange polynomials.
// Passing less shares than the minimum required results in the wrong secret.
func PrivKeyFromShares(shares []*Share) (*PrivateKey, error) {
	// See https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing#Computationally_efficient_approach
	ids := make([]Fq, 0, len(shares))
	secrets := make([]Fq, 0, len(shares))
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

	sum := new(Fq)
	for i := 0; i < len(shares); i++ {
		mul := FqMont1
		for j := 0; j < len(shares); j++ {
			if j != i {
				term := new(Fq)
				FqSub(term, &ids[j], &ids[i])
				FqInv(term, term)
				FqMul(term, term, &ids[j])
				FqMul(&mul, &mul, term)
			}
		}
		FqMul(&mul, &mul, &secrets[i])
		FqAdd(sum, sum, &mul)
	}

	return privKeyFromScalar(FqToBig(FqFromFqMontgomery(*sum))), nil
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
	sum := newG2Point().Set(p.coefficients[0])
	bigX := new(big.Int).SetUint64(x)
	mul := new(big.Int).SetUint64(1)
	for _, coeff := range p.coefficients[1:] {
		mul.Mul(mul, bigX)
		tmp := newG2Point().Add(sum, newG2Point().ScalarMult(coeff, mul))
		sum = tmp
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
	if err != nil {
		return err
	}
	if !pubKey.Equal(&expectedPubKey) {
		return errInvalidShare
	}

	return nil
}
