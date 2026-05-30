BINARY   = dispatch
GOPATH  ?= $(HOME)/go-packages
GO      ?= $(shell which go)

.PHONY: build install clean arm64

build:
	GOPATH=$(GOPATH) $(GO) build -o $(BINARY) .

install: build
	cp $(BINARY) $(HOME)/.local/bin/$(BINARY)

# cross-compile for Raspberry Pi (ARM64)
arm64:
	GOPATH=$(GOPATH) GOOS=linux GOARCH=arm64 $(GO) build -o $(BINARY)-arm64 .

clean:
	rm -f $(BINARY) $(BINARY)-arm64
