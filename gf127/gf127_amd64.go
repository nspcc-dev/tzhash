//go:build amd64 && !generic

package gf127

import "golang.org/x/sys/cpu"

// x127x63 represents x^127 + x^63.
var x127x63 = GF127{msb64, msb64} //nolint:deadcode,varcheck,unused

// Add sets c to a+b.
func Add(a, b, c *GF127) {
	if cpu.X86.HasAVX {
		addAVX(a, b, c)
	} else {
		addGeneric(a, b, c)
	}
}

// Mul sets c to a*b.
func Mul(a, b, c *GF127) {
	if cpu.X86.HasAVX {
		mulAVX(a, b, c)
	} else {
		mulGeneric(a, b, c)
	}
}

// Mul10 sets b to a*x.
func Mul10(a, b *GF127) {
	if cpu.X86.HasAVX {
		mul10AVX(a, b)
	} else {
		mul10Generic(a, b)
	}
}

// Mul11 sets b to a*(x+1).
func Mul11(a, b *GF127) {
	if cpu.X86.HasAVX {
		mul11AVX(a, b)
	} else {
		mul11Generic(a, b)
	}
}

func addAVX(a, b, c *GF127)
func mulAVX(a, b, c *GF127)
func mul10AVX(a, b *GF127)
func mul11AVX(a, b *GF127)
