//go:build !(amd64 && !generic)
// +build !amd64 generic

package gf127

// Mul10x2 sets (b1, b2) to (a1*x, a2*x)
func Mul10x2(a, b *GF127x2) {
	mul10x2Generic(a, b)
}

// Mul11x2 sets (b1, b2) to (a1*(x+1), a2*(x+1))
func Mul11x2(a, b *GF127x2) {
	mul11x2Generic(a, b)
}
