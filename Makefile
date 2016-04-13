GO=go
DCLXVI_DIR=dclxvi-20130329
# Cross-compilers. On Ubuntu, these can be obtained using:
# apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
AR_ARM=/usr/bin/arm-linux-gnueabihf-ar
CC_ARM=/usr/bin/arm-linux-gnueabihf-gcc
CPP_ARM=/usr/bin/arm-linux-gnueabihf-g++
AR_386=$(AR)
CC_386="$(CC) -m32"
CPP_386="$(CPP) -m32"

all: libdclxvi test benchmarks install

.PHONY: libdclxvi

libdclxvi:
	$(MAKE) -C $(DCLXVI_DIR) libdclxvi.a

.PHONY: clean

clean:
	$(GO) clean -i ./...
	$(MAKE) -C $(DCLXVI_DIR) clean
	rm -f bn256.test*

deps: libdclxvi
	$(GO) get -t ./...

test: deps
	$(GO) test ./... -v

install: deps
	$(GO) install ./...

benchmarks: deps
	$(GO) test ./... -run X -bench .

# Targets to build a cross-compiled binary that can be used to run the tests and
# benchmarks on a different architecture.
test-arm:
	USE_C=true $(MAKE) AR=$(AR_ARM) CC=$(CC_ARM) CPP=$(CPP_ARM) -C $(DCLXVI_DIR) libdclxvi.a
	CGO_ENABLED=1 GOARCH=arm CC=$(CC_ARM) CXX=$(CPP_ARM) $(GO) test -c -o bn256.test.arm

test-edison:
	USE_C=true $(MAKE) AR=$(AR_386) CC=$(CC_386) CPP=$(CPP_386) -C $(DCLXVI_DIR) libdclxvi.a
	CGO_ENABLED=1 GOARCH=386 CC=$(CC_386) CXX=$(CPP_386) $(GO) test -c -o bn256.test.edison

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
