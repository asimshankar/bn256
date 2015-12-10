# bn256
[![Build Status](https://travis-ci.org/asimshankar/bn256.svg)](https://travis-ci.org/asimshankar/bn256)
[![GoDoc](https://godoc.org/github.com/asimshankar/bn256?status.svg)](https://godoc.org/github.com/asimshankar/bn256)

Drop in replacement for
[golang.org/x/crypto/bn256](https://godoc.org/golang.org/x/crypto/bn256) backed
by the fast C implementation by Michael Naehrig, Ruben Niederhagen, and Peter
Schwabe described in https://www.cryptojedi.org/crypto/#dclxvi

## Getting Started
```
go get -d github.com/asimshankar/bn256
cd $GOPATH/src/github.com/asimshankar/bn256
make install
```

## LDFLAGS
To use this library in your own Go program, you will likely have to
update the `-L` flag set in `LDFLAGS=` in `bn256.go` to point to the
directory path where this library has been installed.
