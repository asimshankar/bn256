DCLXVI_DIR=dclxvi-20130329

all: libdclxvi test install

.PHONY: libdclxvi

libdclxvi:
	$(MAKE) -C $(DCLXVI_DIR) libdclxvi.a

.PHONY: clean

clean:
	go clean -i ./...
	$(MAKE) -C $(DCLXVI_DIR) clean

test: libdclxvi
	go test ./... -v -bench .

install: libdclxvi
	go install ./...

