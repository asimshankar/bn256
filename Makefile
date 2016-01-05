GO=go
DCLXVI_DIR=dclxvi-20130329
# Cross-compilers. On Ubuntu, these can be obtained using:
# apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
CROSSCOMPILE_CC=/usr/bin/arm-linux-gnueabihf-gcc
CROSSCOMPILE_CPP=/usr/bin/arm-linux-gnueabihf-g++
CROSSCOMPILE_GOARCH=arm
# The Android NDK is makde available by:
# go get golang.org/x/mobile/cmd/gomobile && ${GOPATH}/bin/gomobile init
CROSSCOMPILE_GOMOBILE=$(GOPATH)/pkg/gomobile/android-ndk-r10e/arm/bin

all: libdclxvi test benchmarks install

.PHONY: libdclxvi

libdclxvi:
	$(MAKE) -C $(DCLXVI_DIR) libdclxvi.a

.PHONY: clean

clean:
	$(GO) clean -i ./...
	$(MAKE) -C $(DCLXVI_DIR) clean
	rm -f bn256.test* android.apk

deps: libdclxvi
	$(GO) get -t ./...

test: deps
	$(GO) test ./... -v

install: deps
	$(GO) install ./...

benchmarks: deps
	$(GO) test ./... -run X -bench .

# Target to build a cross-compiled binary that can be used to run the tests and
# benchmarks on a different architecture.
test-crosscompiled:
	$(MAKE) CC=$(CROSSCOMPILE_CC) CPP=$(CROSSCOMPILE_CPP) -C $(DCLXVI_DIR) libdclxvi.a
	CGO_ENABLED=1 GOARCH=$(CROSSCOMPILE_GOARCH) CC=$(CROSSCOMPILE_CC) CXX=$(CROSSCOMPILE_CPP) $(GO) test -c -o bn256.test.${CROSSCOMPILE_GOARCH}

benchmark-android:
	USE_C=true $(MAKE) CC=$(CROSSCOMPILE_GOMOBILE)/arm-linux-androideabi-gcc CPP=$(CROSSCOMPILE_GOMOBILE)/arm-linux-androideabi-g++ -C $(DCLXVI_DIR) libdclxvi.a
	$(GO) get golang.org/x/mobile/cmd/gomobile
	$(GOPATH)/bin/gomobile init
	$(GOPATH)/bin/gomobile build github.com/asimshankar/bn256/android
