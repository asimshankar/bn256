DCLXVI_DIR=dclxvi-20130329

all: libdclxvi test benchmarks install

.PHONY: libdclxvi

libdclxvi:
	$(MAKE) -C $(DCLXVI_DIR) libdclxvi.a

.PHONY: clean

clean:
	go clean -i ./...
	$(MAKE) -C $(DCLXVI_DIR) clean

deps: libdclxvi
	go get -t ./...

test: deps
	go test ./... -v

install: deps
	go install ./...

benchmarks: deps
	go test ./... -run X -bench .

