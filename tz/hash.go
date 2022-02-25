// Package tz contains Tillich-Zemor checksum implementations
// using different backends.
//
// Copyright 2022 (c) NSPCC
package tz

import (
	"errors"
	"hash"

	"golang.org/x/sys/cpu"
)

type Implementation int

const (
	// Size is the size of a Tillich-Zemor hash sum in bytes.
	Size          = 64
	hashBlockSize = 128

	_ Implementation = iota
	AVX
	AVX2
	AVX2Inline
	PureGo
	AVXInline
)

var (
	hasAVX = cpu.X86.HasAVX
	// Having AVX2 does not guarantee
	// that AVX is also present.
	hasAVX2 = cpu.X86.HasAVX2 && hasAVX
)

func (impl Implementation) String() string {
	switch impl {
	case AVX:
		return "AVX"
	case AVXInline:
		return "AVXInline"
	case AVX2:
		return "AVX2"
	case AVX2Inline:
		return "AVX2Inline"
	case PureGo:
		return "PureGo"
	default:
		return "UNKNOWN"
	}
}

func NewWith(impl Implementation) hash.Hash {
	switch impl {
	case AVX:
		return newAVX()
	case AVXInline:
		return newAVXInline()
	case AVX2:
		return newAVX2()
	case AVX2Inline:
		return newAVX2Inline()
	case PureGo:
		return newPure()
	default:
		return New()
	}
}

// New returns a new hash.Hash computing the Tillich-Zémor checksum.
func New() hash.Hash {
	if hasAVX2 {
		return newAVX2Inline()
	} else if hasAVX {
		return newAVXInline()
	} else {
		return newPure()
	}
}

// Sum returns Tillich-Zémor checksum of data.
func Sum(data []byte) [Size]byte {
	if hasAVX2 {
		d := newAVX2Inline()
		_, _ = d.Write(data) // no errors
		return d.checkSum()
	} else if hasAVX {
		d := newAVXInline()
		_, _ = d.Write(data) // no errors
		return d.checkSum()
	} else {
		d := newPure()
		_, _ = d.Write(data) // no errors
		return d.checkSum()
	}
}

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

	p1 = *Inv(&p2)
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

	p2 = *Inv(&p1)
	p2.Mul(&p2, &r)

	return p2.MarshalBinary()
}
