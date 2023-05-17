// Package tz contains Tillich-Zemor checksum implementations
// using different backends.
//
// Copyright 2022 (c) NSPCC
package tz

import (
	"errors"
)

// Concat performs combining of hashes based on homomorphic property.
func Concat(hs [][]byte) ([]byte, error) {
	var b, c sl2

	b = id
	for i := range hs {
		if err := c.UnmarshalBinary(hs[i]); err != nil {
			return nil, err
		}
		b.Mul(&b, &c)
	}
	return b.MarshalBinary()
}

// Validate checks if hashes in hs combined are equal to h.
func Validate(h []byte, hs [][]byte) (bool, error) {
	var (
		b             []byte
		got, expected [Size]byte
		err           error
	)

	if len(h) != Size {
		return false, errors.New("invalid hash")
	} else if len(hs) == 0 {
		return false, errors.New("empty slice")
	}

	copy(expected[:], h)

	b, err = Concat(hs)
	if err != nil {
		return false, errors.New("cant concatenate hashes")
	}

	copy(got[:], b)

	return expected == got, nil
}

// SubtractR returns hash a, such that Concat(a, b) == c
// This is possible, because Tillich-Zemor hash is actually a matrix
// which can be inversed.
func SubtractR(c, b []byte) (a []byte, err error) {
	var p1, p2, r sl2

	if err = r.UnmarshalBinary(c); err != nil {
		return nil, err
	}
	if err = p2.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	p1 = *p2.Inv()
	p1.Mul(&r, &p1)

	return p1.MarshalBinary()
}

// SubtractL returns hash b, such that Concat(a, b) == c
// This is possible, because Tillich-Zemor hash is actually a matrix
// which can be inversed.
func SubtractL(c, a []byte) (b []byte, err error) {
	var p1, p2, r sl2

	if err = r.UnmarshalBinary(c); err != nil {
		return nil, err
	}
	if err = p1.UnmarshalBinary(a); err != nil {
		return nil, err
	}

	p2 = *p1.Inv()
	p2.Mul(&p2, &r)

	return p2.MarshalBinary()
}
