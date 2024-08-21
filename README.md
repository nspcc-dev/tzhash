[![codecov](https://codecov.io/gh/nspcc-dev/tzhash/branch/master/graph/badge.svg)](https://codecov.io/gh/nspcc-dev/tzhash)
[![Report](https://goreportcard.com/badge/github.com/nspcc-dev/tzhash)](https://goreportcard.com/report/github.com/nspcc-dev/tzhash)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nspcc-dev/tzhash?sort=semver)
![License](https://img.shields.io/github/license/nspcc-dev/tzhash.svg?style=popout)

# Demo

[![asciicast](https://asciinema.org/a/IArEDLTrQyabI3agSSpINoqNu.svg)](https://asciinema.org/a/IArEDLTrQyabI3agSSpINoqNu)

**In project root:**

```bash
# show help
make
# run auto demo
make auto
```

# Homomorphic hashing in golang

Package `tz` containts pure-Go implementation of hashing function described by Tillich and Źemor in [1] .

There are existing implementations already (e.g. [2]), however they are written in C.

Package `gf127` contains arithmetic in `GF(2^127)` with `x^127+x^63+1` as reduction polynomial.

# Description

It can be used instead of Merkle-tree for data-validation, because homomorphic hashes
are concatenable: hash sum of data can be calculated based on hashes of chunks.

The example of how it works can be seen in tests.

# Benchmarks

## go vs AVX vs AVX2 version

```
BenchmarkSum/AVX_digest-8             308       3889484 ns/op          25.71 MB/s         5 allocs/op
BenchmarkSum/AVXInline_digest-8       457       2455437 ns/op          40.73 MB/s         5 allocs/op
BenchmarkSum/AVX2_digest-8            399       3031102 ns/op          32.99 MB/s         3 allocs/op
BenchmarkSum/AVX2Inline_digest-8      602       2077719 ns/op          48.13 MB/s         3 allocs/op
BenchmarkSum/PureGo_digest-8           68       17795480 ns/op          5.62 MB/s         5 allocs/op
```

# Contributing

At this moment, we do not accept contributions. Follow us.

# Makefile

```
→ make
  Usage:

    make <target>

  Targets:

    attach   Attach to existing container
    auto     Auto Tillich-Zémor hasher demo
    down     Stop demo container
    help     Show this help prompt
    up       Run Tillich-Zémor hasher demo
```

# Links

[1] https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf

[2] https://github.com/srijs/hwsl2-core