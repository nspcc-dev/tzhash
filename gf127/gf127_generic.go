//go:build !amd64 || generic

package gf127

// Add sets c to a+b.
func Add(a, b, c *GF127) {
	addGeneric(a, b, c)
}

// Mul sets c to a*b.
func Mul(a, b, c *GF127) {
	mulGeneric(a, b, c)
}

// Mul10 sets b to a*x.
func Mul10(a, b *GF127) {
	mul10Generic(a, b)
}

// Mul11 sets b to a*(x+1).
func Mul11(a, b *GF127) {
	mul11Generic(a, b)
}
