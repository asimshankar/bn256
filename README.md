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
make -C ${GOPATH}/src/github.com/asimshankar/bn256 install
# If the above fails (on non-AMD64 platforms for example,
# where the optimized assembly implementation cannot be used)
# then use the slower-than-assembly but faster-than-pure-Go
# portable-C implementation via:
USE_C=true make -C ${GOPATH}/src/github.com/asimshankar/bn256 clean install
```

## ARM devices (like a RaspberryPi)
The instructions above work on an arm architecture device as well. Alternatively,
cross-compile from a more powerful laptop/desktop for the arm processor. From
ubuntu:
```
sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
USE_C=true make -C ${GOPATH}/src/github.com/asimshankar/bn256 clean test-crosscompiled
# This will generate a file like bn256.test.arm which can be copied
# to an run on an arm device, like a RaspberryPi via:
bn256.test.arm --test.v --test.bench=.
```

## Android devices
Tests and benchmarks can be run on Android as well.  It requires the 
[Android NDK](http://developer.android.com/ndk/downloads/index.html)
to be available on the host machine.
```
make -C ${GOPATH}/src/github.com/asimshankar/bn256 clean test-android benchmark-android
```

