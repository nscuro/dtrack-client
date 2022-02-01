build:
	go build -v
.PHONY: build

test:
	go test -v -cover
.PHONY: test

unit-test:
	go test -v -cover -short
.PHONY: unit-test

clean:
	go clean
.PHONY: clean

all: clean build test
.PHONY: all
