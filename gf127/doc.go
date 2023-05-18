/*
Package gf127 implements the GF(2^127) arithmetic modulo (x^127 + x^63 + 1).

This is rather straight-forward re-implementation of C library
available here https://github.com/srijs/hwsl2-core .
Interfaces are highly influenced by [math/big].

gf127.go contains common definitions.
Other files contain architecture-specific implementations.

Copyright 2019 © NSPCC
*/
package gf127
