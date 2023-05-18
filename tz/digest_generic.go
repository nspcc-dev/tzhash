//go:build !(amd64 && !generic)

package tz

func write(d *digest, data []byte) (int, error) {
	return writeGeneric(d, data)
}
