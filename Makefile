build:
	go build -v
.PHONY: build

test:
	go test -v -cover
.PHONY: test

unit-test:
	go test -v -cover -short
.PHONY: unit-test

integration-test:
	go test -v -cover -run IntegrationTest
.PHONY: integration-test

clean:
	go clean
.PHONY: clean

all: clean build test
.PHONY: all
