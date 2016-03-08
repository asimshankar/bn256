GO=go
DCLXVI_DIR=dclxvi-20130329
# Cross-compilers. On Ubuntu, these can be obtained using:
# apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
CROSSCOMPILE_AR=/usr/bin/arm-linux-gnueabihf-ar
CROSSCOMPILE_CC=/usr/bin/arm-linux-gnueabihf-gcc
CROSSCOMPILE_CPP=/usr/bin/arm-linux-gnueabihf-g++
CROSSCOMPILE_GOARCH=arm

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
	$(MAKE) AR=$(CROSSCOMPILE_AR) CC=$(CROSSCOMPILE_CC) CPP=$(CROSSCOMPILE_CPP) -C $(DCLXVI_DIR) libdclxvi.a
	CGO_ENABLED=1 GOARCH=$(CROSSCOMPILE_GOARCH) CC=$(CROSSCOMPILE_CC) CXX=$(CROSSCOMPILE_CPP) $(GO) test -c -o bn256.test.${CROSSCOMPILE_GOARCH}

.PHONY: android

ifndef ANDROID_NDK
android:
	@echo "ANDROID_NDK must be set. See http://developer.android.com/tools/sdk/ndk/index.html for installation" && false
else
android:
	USE_C=true $(MAKE) AR=$(ANDROID_NDK)/arm-linux-androideabi-ar CC=$(ANDROID_NDK)/arm-linux-androideabi-gcc CPP=$(ANDROID_NDK)/arm-linux-androideabi-g++ -C $(DCLXVI_DIR) libdclxvi.a
	$(GO) get -d -t ./...
	$(GO) get v.io/x/devtools/bendroid
endif

benchmark-android: android
	CC=$(ANDROID_NDK)/arm-linux-androideabi-gcc CXX=$(ANDROID_NDK)/arm-linux-androideabi-g++ CGO_ENABLED=1 GOARCH=arm GOOS=android $(GO) test -run NONE -bench . -exec $(GOPATH)/bin/bendroid github.com/asimshankar/bn256

test-android: android
	CC=$(ANDROID_NDK)/arm-linux-androideabi-gcc CXX=$(ANDROID_NDK)/arm-linux-androideabi-g++ CGO_ENABLED=1 GOARCH=arm GOOS=android $(GO) test -v -exec $(GOPATH)/bin/bendroid github.com/asimshankar/bn256
